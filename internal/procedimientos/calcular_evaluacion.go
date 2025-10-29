package procedimientos

import (
	"fmt"
	"math"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

// CalcularEvaluacion calcula y actualiza los registros de ConceptosEvaluacion
// para los años 0..5 de un plan. Reglas:
// - FlujoEfectivoNominal = (ComposicionFinanciamiento.CapitalPorcentaje/100) * Total_Inversion
// - ValorRescate = 0 para todos los años
// - TotalFlujoEfectivo = FlujoEfectivoNominal + ValorRescate
// - ValorActualFlujosFuturos:
//   - para años != 0: TotalFlujoEfectivo / (1 + TREMA/100)^{anio}
//   - para año 0: suma de los valores actualizados de los años 1..5
func CalcularEvaluacion(db *gorm.DB, planID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// cargar composicion de financiamiento
		var comp models.ComposicionFinanciamiento
		if err := tx.Where("plan_negocio_id = ?", planID).First(&comp).Error; err != nil {
			return fmt.Errorf("composicion_financiamiento no encontrada para plan %d: %w", planID, err)
		}

		// cargar evaluacion proyecto (contiene TREMA)
		var eval models.EvaluacionProyecto
		if err := tx.Where("plan_negocio_id = ?", planID).First(&eval).Error; err != nil {
			return fmt.Errorf("evaluacion_proyecto no encontrada para plan %d: %w", planID, err)
		}

		// Calcular ValorRescate: 0 para años 0-4, calculado para año 5 basado en BalanceGeneral
		valorRescatePorAnio := make(map[int]float64)

		// Años 0-4: ValorRescate = 0
		for anio := 0; anio <= 4; anio++ {
			valorRescatePorAnio[anio] = 0.0
		}

		// Año 5: calcular basado en BalanceGeneral
		var balancesAnio5 []models.BalanceGeneral
		if err := tx.Where("plan_negocio_id = ? AND anio = ?", planID, 5).Find(&balancesAnio5).Error; err != nil {
			return fmt.Errorf("error al buscar BalanceGeneral para año 5: %w", err)
		}

		// Sumar valores del año 5
		var sumaCuentasxCobrar, sumaInventarios, sumaNoCorrientes, sumaCuentasxPagar, sumaOtrosCortoplazo float64
		for _, balance := range balancesAnio5 {
			sumaCuentasxCobrar += balance.Corrientes_CuentasxCobrar
			sumaInventarios += balance.Corrientes_Inventarios
			sumaNoCorrientes += balance.NoCorrientes_Suma
			sumaCuentasxPagar += balance.PasivoCuentasxPagarCortoPlazo
			sumaOtrosCortoplazo += balance.PasivoOtrosCortoPlazo
		}

		// ValorRescate Año 5 = CuentasxCobrar + Inventarios + NoCorrientes - CuentasxPagar - OtrosCortoPlazo
		valorRescatePorAnio[5] = sumaCuentasxCobrar + sumaInventarios + sumaNoCorrientes - sumaCuentasxPagar - sumaOtrosCortoplazo

		// Calcular FlujoEfectivoNominal por año:
		// - Año 0: composicion (capital_porcentaje * total_inversion)
		// - Años 1-5: suma de FlujoCaja de FlujoEfectivo para ese año
		flujoNominalPorAnio := make(map[int]float64)

		// Año 0: usar composición financiera (con signo negativo)
		flujoNominalPorAnio[0] = -((comp.CapitalPorcentaje / 100.0) * comp.Total_Inversion)

		// Años 1-5: sumar FlujoCaja de FlujoEfectivo por año
		for anio := 1; anio <= 5; anio++ {
			var flujos []models.FlujoEfectivo
			if err := tx.Where("plan_negocio_id = ? AND anio = ?", planID, anio).Find(&flujos).Error; err != nil {
				return fmt.Errorf("error al buscar FlujoEfectivo para año %d: %w", anio, err)
			}
			suma := 0.0
			for _, flujo := range flujos {
				suma += flujo.FlujoCaja
			}
			flujoNominalPorAnio[anio] = suma
		}

		// Calcular TotalFlujoEfectivo y ValorActualFlujosFuturos para años 1-5
		totalFlujoPorAnio := make(map[int]float64)
		valorActualPorAnio := make(map[int]float64)
		sumaValoresActuales := 0.0

		for anio := 1; anio <= 5; anio++ {
			totalFlujo := flujoNominalPorAnio[anio] + valorRescatePorAnio[anio]
			totalFlujoPorAnio[anio] = totalFlujo

			// Calcular valor actual descontado
			factor := math.Pow(1.0+(eval.TREMA/100.0), float64(anio))
			valorActual := totalFlujo / factor
			valorActualPorAnio[anio] = valorActual
			sumaValoresActuales += valorActual
		}

		// Año 0: TotalFlujoEfectivo y ValorActualFlujosFuturos
		totalFlujoPorAnio[0] = flujoNominalPorAnio[0] + valorRescatePorAnio[0]
		valorActualPorAnio[0] = sumaValoresActuales // Suma de valores actuales de años 1-5

		// Iterar años 0-5 y upsert en ConceptosEvaluacion
		for anio := 0; anio <= 5; anio++ {
			flujoNominalStr := fmt.Sprintf("%.2f", flujoNominalPorAnio[anio])
			valorRescateStr := fmt.Sprintf("%.2f", valorRescatePorAnio[anio])
			totalFlujoStr := fmt.Sprintf("%.2f", totalFlujoPorAnio[anio])
			valorActualStr := fmt.Sprintf("%.2f", valorActualPorAnio[anio])

			var ce models.ConceptosEvaluacion
			err := tx.Where("plan_negocio_id = ? AND anio = ?", planID, anio).First(&ce).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					ce = models.ConceptosEvaluacion{
						PlanNegocioID:            planID,
						Anio:                     anio,
						FlujoEfectivoNominal:     flujoNominalStr,
						ValorRescate:             valorRescateStr,
						TotalFlujoEfectivo:       totalFlujoStr,
						ValorActualFlujosFuturos: valorActualStr,
					}
					if err := tx.Create(&ce).Error; err != nil {
						return fmt.Errorf("crear ConceptosEvaluacion anio %d: %w", anio, err)
					}
					continue
				}
				return fmt.Errorf("leer ConceptosEvaluacion anio %d: %w", anio, err)
			}

			// actualizar campos existentes
			ce.FlujoEfectivoNominal = flujoNominalStr
			ce.ValorRescate = valorRescateStr
			ce.TotalFlujoEfectivo = totalFlujoStr
			ce.ValorActualFlujosFuturos = valorActualStr

			if err := tx.Save(&ce).Error; err != nil {
				return fmt.Errorf("actualizar ConceptosEvaluacion anio %d: %w", anio, err)
			}
		}

		// Calcular VAN y TIR para EvaluacionProyecto
		// VAN = TotalFlujoEfectivo año 0 + ValorActualFlujosFuturos año 0
		van := totalFlujoPorAnio[0] + valorActualPorAnio[0]

		// Preparar flujos para cálculo de TIR (años 0 a 5)
		flujosTIR := make([]float64, 6)
		for anio := 0; anio <= 5; anio++ {
			flujosTIR[anio] = totalFlujoPorAnio[anio]
		}

		// Calcular TIR usando algoritmo de Newton-Raphson
		tir := calcularTIR(flujosTIR)

		// Actualizar EvaluacionProyecto
		var evalProyecto models.EvaluacionProyecto
		err := tx.Where("plan_negocio_id = ?", planID).First(&evalProyecto).Error
		if err != nil {
			return fmt.Errorf("error al buscar EvaluacionProyecto: %w", err)
		}

		evalProyecto.VAN = van
		evalProyecto.TIR = tir

		if err := tx.Save(&evalProyecto).Error; err != nil {
			return fmt.Errorf("error al actualizar EvaluacionProyecto: %w", err)
		}

		return nil
	})
}

// calcularTIR implementa el algoritmo de Newton-Raphson para calcular la TIR
// Equivalente a la función IRR de Excel
func calcularTIR(flujos []float64) float64 {
	// Valores iniciales
	tasa := 0.1 // Tasa inicial del 10%
	precision := 0.0001
	maxIteraciones := 100

	for i := 0; i < maxIteraciones; i++ {
		// Calcular VPN con la tasa actual
		vpn := 0.0
		derivada := 0.0

		for j, flujo := range flujos {
			factor := math.Pow(1+tasa, float64(j))
			vpn += flujo / factor

			// Calcular derivada para Newton-Raphson
			if j > 0 {
				derivada -= float64(j) * flujo / math.Pow(1+tasa, float64(j+1))
			}
		}

		// Si VPN es suficientemente cercano a cero, hemos encontrado la TIR
		if math.Abs(vpn) < precision {
			return tasa * 100 // Convertir a porcentaje
		}

		// Evitar división por cero
		if math.Abs(derivada) < precision {
			break
		}

		// Actualizar tasa usando Newton-Raphson
		nuevaTasa := tasa - vpn/derivada

		// Limitar la tasa para evitar valores extremos
		if nuevaTasa < -0.99 {
			nuevaTasa = -0.99
		} else if nuevaTasa > 10 {
			nuevaTasa = 10
		}

		tasa = nuevaTasa
	}

	return tasa * 100 // Convertir a porcentaje
}
