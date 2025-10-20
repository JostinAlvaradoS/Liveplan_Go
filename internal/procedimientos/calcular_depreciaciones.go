package procedimientos

import (
	"fmt"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

// floatPtr helper
func floatPtr(v float64) *float64 { return &v }

// intPtr helper
func intPtr(v int) *int { return &v }

// CalcularDepreciaciones recalcula las filas de la tabla depreciaciones para un plan dado.
// Reglas:
//   - Recorre todos los DetalleInversionInicial del plan.
//   - Si Detalle.VidaUtil <= 0 entonces se ignora (no se calcula).
//   - VidaUtil se interpreta en meses. La depreciación mensual = Importe / VidaUtil.
//   - Para cada año i=1..5 calculamos la depreciación anual como monthly * meses_en_el_año
//     (meses_en_el_año = min(12, max(0, VidaUtil - (i-1)*12))). Si meses_en_el_año == 0 el año queda NULL.
//   - DepreciacionMensual se guarda como monthly, DepreciacionAnio1..5 con los valores calculados
//   - ValorRescate = Importe - suma(depreciaciones de los 5 años)
//   - Si no existe un registro en `depreciaciones` para el detalle, se crea.
func CalcularDepreciaciones(db *gorm.DB, planID uint) error {
	// run in transaction for consistency
	return db.Transaction(func(tx *gorm.DB) error {
		var detalles []models.DetalleInversionInicial
		if err := tx.Where("plan_negocio_id = ?", planID).Find(&detalles).Error; err != nil {
			return err
		}

		for _, d := range detalles {
			if d.VidaUtil <= 0 {
				// If vida util is zero or negative we must ensure the depreciation entry
				// exists and that all depreciation fields are NULL (no depreciation applies).
				var existing models.Depreciacion
				err := tx.Where("detalle_inversion_id = ?", d.ID).First(&existing).Error
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						placeholder := models.Depreciacion{
							PlanNegocioID:      d.PlanNegocioID,
							DetalleInversionID: d.ID,
						}
						if err := tx.Create(&placeholder).Error; err != nil {
							return fmt.Errorf("creating placeholder depreciacion for detalle %d: %w", d.ID, err)
						}
					} else {
						return err
					}
				} else {
					// existing record -> nullify calculated fields
					existing.DepreciacionMensual = nil
					existing.DepreciacionAnio1 = nil
					existing.DepreciacionAnio2 = nil
					existing.DepreciacionAnio3 = nil
					existing.DepreciacionAnio4 = nil
					existing.DepreciacionAnio5 = nil
					existing.ValorRescate = nil
					if err := tx.Save(&existing).Error; err != nil {
						return fmt.Errorf("clearing depreciacion for detalle %d: %w", d.ID, err)
					}
				}
				continue
			}

			importe := d.Importe
			vidaMeses := d.VidaUtil
			monthly := importe / float64(vidaMeses)

			// compute per-year depreciation for up to 5 years
			years := make([]*float64, 5)
			monthsRemaining := vidaMeses
			var sumYears float64
			for i := 0; i < 5; i++ {
				if monthsRemaining <= 0 {
					years[i] = nil
					continue
				}
				monthsInYear := 12
				if monthsRemaining < 12 {
					monthsInYear = monthsRemaining
				}
				val := monthly * float64(monthsInYear)
				years[i] = floatPtr(val)
				sumYears += val
				monthsRemaining -= monthsInYear
			}

			valorRescate := importe - sumYears
			dep := models.Depreciacion{
				PlanNegocioID:       d.PlanNegocioID,
				DetalleInversionID:  d.ID,
				DepreciacionMensual: floatPtr(monthly),
				DepreciacionAnio1:   years[0],
				DepreciacionAnio2:   years[1],
				DepreciacionAnio3:   years[2],
				DepreciacionAnio4:   years[3],
				DepreciacionAnio5:   years[4],
				ValorRescate:        floatPtr(valorRescate),
			}

			var existing models.Depreciacion
			err := tx.Where("detalle_inversion_id = ?", d.ID).First(&existing).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					if err := tx.Create(&dep).Error; err != nil {
						return fmt.Errorf("creating depreciacion for detalle %d: %w", d.ID, err)
					}
				} else {
					return err
				}
			} else {
				existing.DepreciacionMensual = dep.DepreciacionMensual
				existing.DepreciacionAnio1 = dep.DepreciacionAnio1
				existing.DepreciacionAnio2 = dep.DepreciacionAnio2
				existing.DepreciacionAnio3 = dep.DepreciacionAnio3
				existing.DepreciacionAnio4 = dep.DepreciacionAnio4
				existing.DepreciacionAnio5 = dep.DepreciacionAnio5
				existing.ValorRescate = dep.ValorRescate
				if err := tx.Save(&existing).Error; err != nil {
					return fmt.Errorf("updating depreciacion for detalle %d: %w", d.ID, err)
				}
			}
		}
		return nil
	})
}
