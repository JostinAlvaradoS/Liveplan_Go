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

		// ValorRescate siempre es 0
		valorRescate := 0.0

		// Calcular FlujoEfectivoNominal por año:
		// - Año 0: composicion (capital_porcentaje * total_inversion)
		// - Años 1-5: suma de FlujoCaja de FlujoEfectivo para ese año
		flujoNominalPorAnio := make(map[int]float64)

		// Año 0: usar composición financiera
		flujoNominalPorAnio[0] = (comp.CapitalPorcentaje / 100.0) * comp.Total_Inversion

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
			totalFlujo := flujoNominalPorAnio[anio] + valorRescate
			totalFlujoPorAnio[anio] = totalFlujo

			// Calcular valor actual descontado
			factor := math.Pow(1.0+(eval.TREMA/100.0), float64(anio))
			valorActual := totalFlujo / factor
			valorActualPorAnio[anio] = valorActual
			sumaValoresActuales += valorActual
		}

		// Año 0: TotalFlujoEfectivo y ValorActualFlujosFuturos
		totalFlujoPorAnio[0] = flujoNominalPorAnio[0] + valorRescate
		valorActualPorAnio[0] = sumaValoresActuales // Suma de valores actuales de años 1-5

		// Iterar años 0-5 y upsert en ConceptosEvaluacion
		for anio := 0; anio <= 5; anio++ {
			flujoNominalStr := fmt.Sprintf("%.2f", flujoNominalPorAnio[anio])
			valorRescateStr := fmt.Sprintf("%.2f", valorRescate)
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

		return nil
	})
}
