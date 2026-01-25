package procedimientos

import (
	"fmt"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

func SetStatusSensibilidad(db *gorm.DB, planID uint, status bool) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var analisis []models.AnalisisSensibilidad_Status
		if err := tx.Where("plan_negocio_id = ?", planID).Find(&analisis).Error; err != nil {
			return fmt.Errorf("cargando el status de sensibilidad: %w", err)
		}
		for _, a := range analisis {
			a.Status = status
			if err := tx.Save(&a).Error; err != nil {
				return fmt.Errorf("actualizando el status de sensibilidad: %w", err)
			}
		}

		return nil
	})
}
