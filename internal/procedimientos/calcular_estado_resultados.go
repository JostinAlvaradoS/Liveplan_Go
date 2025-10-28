package procedimientos

import (
	"fmt"
	"log"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

func CalcularEstadoResultados(db *gorm.DB, planID uint) error {
	// Sumar intereses de PrestamoCuotas por anio y mes (gastos financieros)
	prestamosPorAnioMes := make(map[int]map[int]float64)
	var cuotas []models.PrestamoCuotas
	if err := db.Where("plan_negocio_id = ?", planID).Find(&cuotas).Error; err != nil {
		return fmt.Errorf("obtener prestamo_cuotas: %w", err)
	}
	for _, c := range cuotas {
		if _, ok := prestamosPorAnioMes[c.Anio]; !ok {
			prestamosPorAnioMes[c.Anio] = make(map[int]float64)
		}
		prestamosPorAnioMes[c.Anio][c.Mes] += c.Interes
	}

	// Obtener indicadores macro (PTU, tasa impuesto)
	var ind models.IndicadoresMacro
	if err := db.Where("plan_negocio_id = ?", planID).First(&ind).Error; err != nil {
		ind = models.IndicadoresMacro{}
	}
	// Sumar costos por año y mes desde CostosVentas
	costosPorAnioMes := make(map[int]map[int]float64) // anio -> mes -> suma
	var costosVentas []models.CostosVentas
	if err := db.Where("plan_negocio_id = ?", planID).Find(&costosVentas).Error; err != nil {
		return fmt.Errorf("obtener costos_ventas: %w", err)
	}
	for _, cv := range costosVentas {
		if _, ok := costosPorAnioMes[cv.Anio]; !ok {
			costosPorAnioMes[cv.Anio] = make(map[int]float64)
		}
		costosPorAnioMes[cv.Anio][cv.Mes] += cv.Costo
	}
	// Determinar los años a procesar. Preferir los años ya creados en EstadoResultados
	yearsSet := make(map[int]struct{})
	var ers []models.EstadoResultados
	if err := db.Where("plan_negocio_id = ?", planID).Find(&ers).Error; err != nil {
		return fmt.Errorf("obtener estado_resultados existentes: %w", err)
	}
	for _, e := range ers {
		yearsSet[e.Anio] = struct{}{}
	}

	// Obtener ventas por año
	var ventas []models.Ventas
	if err := db.Where("plan_negocio_id = ?", planID).Find(&ventas).Error; err != nil {
		return fmt.Errorf("obtener ventas: %w", err)
	}
	ventasPorAnio := make(map[int]float64)
	for _, v := range ventas {
		ventasPorAnio[v.Anio] += v.Venta
	}

	// Si no hay años creados, usar los años hallados en ventas
	if len(yearsSet) == 0 {
		for anio := range ventasPorAnio {
			yearsSet[anio] = struct{}{}
		}
	}

	// Si aún no hay años (no hay ventas ni registros), no hacemos nada
	if len(yearsSet) == 0 {
		log.Printf("CalcularEstadoResultados: no hay años para procesar en plan %d", planID)
		return nil
	}

	// Sumar GastosVentaAdm por mes y año
	gastosVentaAdmPorAnioMes := make(map[int]map[int]float64)
	var gastosOperacion []models.GastosOperacion
	if err := db.Where("plan_negocio_id = ?", planID).Find(&gastosOperacion).Error; err != nil {
		return fmt.Errorf("obtener gastos_operacion: %w", err)
	}
	// GastosOperacion no tiene mes/anio, se asigna igual a todos los meses
	for anio := range yearsSet {
		if _, ok := gastosVentaAdmPorAnioMes[anio]; !ok {
			gastosVentaAdmPorAnioMes[anio] = make(map[int]float64)
		}
		totalGastos := 0.0
		for _, gope := range gastosOperacion {
			totalGastos += gope.Costo
		}
		for mes := 1; mes <= 12; mes++ {
			gastosVentaAdmPorAnioMes[anio][mes] = totalGastos
		}
	}

	// Sumar Depreciacion y Amortizacion por mes y año
	depreciacionPorAnioMes := make(map[int]map[int]float64)
	amortizacionPorAnioMes := make(map[int]map[int]float64)
	var depreciaciones []models.Depreciacion
	if err := db.Where("plan_negocio_id = ?", planID).Preload("DetalleInversion").Find(&depreciaciones).Error; err != nil {
		return fmt.Errorf("obtener depreciaciones: %w", err)
	}
	for anio := range yearsSet {
		if _, ok := depreciacionPorAnioMes[anio]; !ok {
			depreciacionPorAnioMes[anio] = make(map[int]float64)
		}
		if _, ok := amortizacionPorAnioMes[anio]; !ok {
			amortizacionPorAnioMes[anio] = make(map[int]float64)
		}
		for mes := 1; mes <= 12; mes++ {
			depSum := 0.0
			amoSum := 0.0
			for _, dep := range depreciaciones {
				tipo := 0
				if dep.DetalleInversion != nil {
					tipo = int(dep.DetalleInversion.TipoID)
				}
				val := 0.0
				if dep.DepreciacionMensual != nil {
					val = *dep.DepreciacionMensual
				}
				if tipo == 1 {
					depSum += val
				} else if tipo == 2 {
					amoSum += val
				}
			}
			depreciacionPorAnioMes[anio][mes] = depSum
			amortizacionPorAnioMes[anio][mes] = amoSum
		}
	}

	// Para cada año y cada mes actualizar o crear el registro correspondiente
	for anio := range yearsSet {
		ventasAnio := ventasPorAnio[anio]
		for mes := 1; mes <= 12; mes++ {
			costosMes := 0.0
			if m, ok := costosPorAnioMes[anio]; ok {
				costosMes = m[mes]
			}
			gastosVentaAdm := 0.0
			if m, ok := gastosVentaAdmPorAnioMes[anio]; ok {
				gastosVentaAdm = m[mes]
			}
			depreciacion := 0.0
			if m, ok := depreciacionPorAnioMes[anio]; ok {
				depreciacion = m[mes]
			}
			amortizacion := 0.0
			if m, ok := amortizacionPorAnioMes[anio]; ok {
				amortizacion = m[mes]
			}
			utilidadBruta := ventasAnio - costosMes
			utilidadPrevioIntImp := utilidadBruta - gastosVentaAdm - depreciacion - amortizacion

			gastosFinancieros := 0.0
			if pa, ok := prestamosPorAnioMes[anio]; ok {
				gastosFinancieros = pa[mes]
			}
			utilidadAntesPTU := utilidadPrevioIntImp - gastosFinancieros
			ptu := utilidadAntesPTU * (ind.PTU / 100.0)
			utilidadAntesImpuestos := utilidadAntesPTU - ptu
			isr := utilidadAntesImpuestos * (ind.TasaImpuesto / 100.0)
			utilidadNeta := utilidadAntesImpuestos - isr

			var er models.EstadoResultados
			q := db.Where("plan_negocio_id = ? AND anio = ? AND mes = ?", planID, anio, mes).First(&er)
			if q.Error == nil {
				er.Ventas = ventasAnio
				er.CostosVentas = costosMes
				er.UtilidadBruta = utilidadBruta
				er.GastosVentaAdm = gastosVentaAdm
				er.Depreciacion = depreciacion
				er.Amortizacion = amortizacion
				er.UtilidadprevioIntImp = utilidadPrevioIntImp
				er.GastosFinancieros = gastosFinancieros
				er.UtilidadAntesPTU = utilidadAntesPTU
				er.PTU = ptu
				er.UtilidadAntesImpuestos = utilidadAntesImpuestos
				er.ISR = isr
				er.UtilidadNeta = utilidadNeta
				if err := db.Save(&er).Error; err != nil {
					log.Printf("CalcularEstadoResultados: error actualizando estado resultados P:%d A:%d M:%d: %v", planID, anio, mes, err)
					return fmt.Errorf("actualizar estado_resultados: %w", err)
				}
			} else if q.Error == gorm.ErrRecordNotFound {
				newEr := models.EstadoResultados{
					PlanNegocioID:          planID,
					Anio:                   anio,
					Mes:                    mes,
					Ventas:                 ventasAnio,
					CostosVentas:           costosMes,
					UtilidadBruta:          utilidadBruta,
					GastosVentaAdm:         gastosVentaAdm,
					Depreciacion:           depreciacion,
					Amortizacion:           amortizacion,
					UtilidadprevioIntImp:   utilidadPrevioIntImp,
					GastosFinancieros:      gastosFinancieros,
					UtilidadAntesPTU:       utilidadAntesPTU,
					PTU:                    ptu,
					UtilidadAntesImpuestos: utilidadAntesImpuestos,
					ISR:                    isr,
					UtilidadNeta:           utilidadNeta,
				}
				if err := db.Create(&newEr).Error; err != nil {
					log.Printf("CalcularEstadoResultados: error creando estado resultados P:%d A:%d M:%d: %v", planID, anio, mes, err)
					return fmt.Errorf("crear estado_resultados: %w", err)
				}
			} else {
				return fmt.Errorf("consultar estado_resultados: %w", q.Error)
			}
		}
	}

	return nil
}
