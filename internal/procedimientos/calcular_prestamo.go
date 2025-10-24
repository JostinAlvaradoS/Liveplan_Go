package procedimientos

import (
	"fmt"
	"math"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

// CalcularPrestamo genera la tabla de amortización en prestamo_cuotas para
// el plan indicado usando los datos en datos_prestamo. Se asume fórmula de
// anualidad (cuota fija mensual):
//
//	r = tasa_interes / 100 / 12 (ajustado por periodos de capitalización)
//	cuota = P * r / (1 - (1+r)^-n) (si no está definida en DatosPrestamo)
//
// Donde P = Monto, n = PeriodosAmortizacion. La función actualiza las filas
// existentes en prestamo_cuotas ordenadas por periodo_mes; si hay más filas
// de las necesarias, las deja; pero actualiza hasta n periodos.
func CalcularPrestamo(db *gorm.DB, planID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var dp models.DatosPrestamo
		if err := tx.Where("plan_negocio_id = ?", planID).First(&dp).Error; err != nil {
			return fmt.Errorf("loading datos_prestamo for plan %d: %w", planID, err)
		}

		// Determine principal P. Prefer composition-based debt amount when possible:
		// DeudaPorcentaje (%) of ComposicionFinanciamiento.Total_Inversion.
		P := dp.Monto
		var cf models.ComposicionFinanciamiento
		if err := tx.Where("plan_negocio_id = ?", planID).First(&cf).Error; err == nil {
			if cf.Total_Inversion > 0 && cf.DeudaPorcentaje > 0 {
				derived := cf.Total_Inversion * (cf.DeudaPorcentaje / 100.0)
				// if derived differs from dp.Monto (and dp.Monto is zero or smaller), prefer derived
				if dp.Monto == 0 || math.Abs(derived-dp.Monto) > 0.01 {
					P = derived
					// persist the derived monto back to datos_prestamo for transparency
					if err := tx.Model(&dp).Update("monto", P).Error; err != nil {
						return fmt.Errorf("updating derived monto in datos_prestamo for plan %d: %w", planID, err)
					}
				}
			}
		}
		n := dp.PeriodosAmortizacion
		if n <= 0 {
			var count int64
			if err := tx.Model(&models.PrestamoCuotas{}).Where("plan_negocio_id = ?", planID).Count(&count).Error; err != nil {
				return fmt.Errorf("counting prestamo_cuotas for plan %d: %w", planID, err)
			}
			if count > 0 {
				n = int(count)
			} else {
				n = 60 // Default a 5 años si no se especifica
			}
		}

		// Recalcular siempre la tasa mensual a partir de la tasa anual y periodos de capitalización
		// según: tasa_mensual = tasa_anual / periodos_capitalizacion
		if dp.TasaAnual != 0 && dp.PeriodosCapitalizacion > 0 {
			// mantendremos la misma unidad que usa DatosPrestamo (porcentaje), por lo que
			// asignamos el valor en unidades de porcentaje y lo persistimos.
			dp.TasaMensual = dp.TasaAnual / float64(dp.PeriodosCapitalizacion)
			if err := tx.Model(&dp).Update("tasa_mensual", dp.TasaMensual).Error; err != nil {
				return fmt.Errorf("updating tasa_mensual in datos_prestamo for plan %d: %w", planID, err)
			}
		}

		// Determinar tasa mensual (convertir de porcentaje a decimal)
		var r float64
		if dp.TasaMensual != 0 {
			r = dp.TasaMensual / 100.0 // Convertir de porcentaje (ej. 1% -> 0.01)
		} else if dp.TasaAnual != 0 {
			// Convertir tasa anual (porcentaje) a tasa mensual efectiva considerando periodos de capitalización
			annualRate := dp.TasaAnual / 100.0 // Convertir de porcentaje a decimal (ej. 12% -> 0.12)
			r = math.Pow(1+annualRate, 1.0/float64(dp.PeriodosCapitalizacion)) - 1.0
		} else {
			return fmt.Errorf("no tasa de interés (anual o mensual) proporcionada para plan %d", planID)
		}

		// Usar cuota proporcionada si existe, sino calcularla
		var cuota float64
		if dp.Cuota != 0 {
			cuota = dp.Cuota
		} else if r == 0 {
			cuota = P / float64(n) // Pago igual si no hay interés
		} else {
			cuota = P * r / (1 - math.Pow(1+r, -float64(n)))
		}

		// Cargar cuotas existentes ordenadas por periodo_mes
		var cuotas []models.PrestamoCuotas
		if err := tx.Where("plan_negocio_id = ?", planID).Order("periodo_mes asc").Find(&cuotas).Error; err != nil {
			return fmt.Errorf("loading prestamo_cuotas for plan %d: %w", planID, err)
		}

		// Ajuste de número de cuotas: si el conjunto actual no tiene exactamente n filas
		// (por ejemplo, quedó muy grande por periodos de capitalización altos), eliminamos
		// todas las filas y recreamos exactamente n entradas (periodo_mes = 1..n).
		if len(cuotas) != n {
			if err := tx.Where("plan_negocio_id = ?", planID).Delete(&models.PrestamoCuotas{}).Error; err != nil {
				return fmt.Errorf("resetting prestamo_cuotas for plan %d: %w", planID, err)
			}
			for m := 1; m <= n; m++ {
				anio := (m-1)/12 + 1
				mes := (m-1)%12 + 1
				pc := models.PrestamoCuotas{
					PlanNegocioID:  planID,
					PeriodoMes:     m,
					Anio:           anio,
					Mes:            mes,
					SaldoInicial:   0,
					Interes:        0,
					Amortizacion:   0,
					CuotaTotal:     0,
					SaldoPendiente: 0,
				}
				if err := tx.Create(&pc).Error; err != nil {
					return fmt.Errorf("creating prestamo_cuotas for plan %d: %w", planID, err)
				}
			}
			// Recargar cuotas recién creadas
			if err := tx.Where("plan_negocio_id = ?", planID).Order("periodo_mes asc").Find(&cuotas).Error; err != nil {
				return fmt.Errorf("reloading prestamo_cuotas for plan %d: %w", planID, err)
			}
		}

		// Calcular y actualizar la amortización
		saldo := P
		for i := 0; i < n; i++ {
			cuotaRow := &cuotas[i]
			// saldo inicial para este periodo = saldo antes del pago
			saldoInicial := saldo
			interes := saldoInicial * r
			amort := cuota - interes
			if amort < 0 {
				amort = 0 // Ajuste para el último pago
			}
			if saldoInicial < cuota {
				cuota = saldoInicial + interes // Ajuste final
				amort = saldoInicial
			}
			saldo = saldoInicial - amort
			if saldo < 0 {
				saldo = 0
			}

			updates := map[string]interface{}{
				"saldo_inicial":   saldoInicial,
				"interes":         interes,
				"amortizacion":    amort,
				"cuota_total":     cuota,
				"saldo_pendiente": saldo,
			}
			if err := tx.Model(&models.PrestamoCuotas{}).Where("id = ?", cuotaRow.ID).Updates(updates).Error; err != nil {
				return fmt.Errorf("updating prestamo_cuotas id %d: %w", cuotaRow.ID, err)
			}
		}

		// Actualizar el campo Cuota en DatosPrestamo si fue calculado
		if dp.Cuota == 0 {
			if err := tx.Model(&dp).Update("cuota", cuota).Error; err != nil {
				return fmt.Errorf("updating cuota in datos_prestamo for plan %d: %w", planID, err)
			}
		}

		return nil
	})
}
