package procedimientos

import (
	"fmt"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

// CalcularComposicion calcula el total de inversión para la tabla
// ComposicionFinanciamiento de un plan sumando los importes de todos los
// DetalleInversionInicial asociados al plan y actualiza el campo
// total_inversion.
func CalcularComposicion(db *gorm.DB, planID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var detalles []models.DetalleInversionInicial
		if err := tx.Where("plan_negocio_id = ?", planID).Find(&detalles).Error; err != nil {
			return fmt.Errorf("loading detalle_inversion for plan %d: %w", planID, err)
		}

		var total float64
		for _, d := range detalles {
			total += d.Importe
		}

		// upsert: try to find existing composicion_financiamiento row for plan
		var comp models.ComposicionFinanciamiento
		err := tx.Where("plan_negocio_id = ?", planID).First(&comp).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// create new
				comp = models.ComposicionFinanciamiento{
					PlanNegocioID:   planID,
					Total_Inversion: total,
				}
				if err := tx.Create(&comp).Error; err != nil {
					return fmt.Errorf("creating composicion_financiamiento for plan %d: %w", planID, err)
				}
				return nil
			}
			return fmt.Errorf("loading composicion_financiamiento for plan %d: %w", planID, err)
		}

		// update existing
		comp.Total_Inversion = total
		if err := tx.Save(&comp).Error; err != nil {
			return fmt.Errorf("updating composicion_financiamiento for plan %d: %w", planID, err)
		}
		return nil
	})
}
