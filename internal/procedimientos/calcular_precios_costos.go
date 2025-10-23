package procedimientos

import (
	"fmt"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

// CalcularPreciosYCostosPorPlan recalcula todos los precios del plan y, para
// cada producto, recalcula todos los costos asociados usando multiplicadores
// según la categoría de costo:
//
//	categoria 1 -> 0.10
//	categoria 2 -> 0.12
//	categoria 3 -> 0.15
//
// La fórmula de precio: precio_calc = precio * (1 + variables.precio)
// La fórmula de costo: costo_calc = (costo * multiplicador) * (1 + variables.costo)
func CalcularPreciosYCostosPorPlan(db *gorm.DB, planID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// load variables de sensibilidad for plan (if none found assume 0)
		var vs models.VariablesDeSensibilidad
		precioFactor := 0.0
		costoFactor := 0.0
		if err := tx.Where("plan_negocio_id = ?", planID).First(&vs).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return fmt.Errorf("loading variables_de_sensibilidad: %w", err)
			}
		} else {
			precioFactor = vs.Precio
			costoFactor = vs.Costo
		}

		// load all precios for plan
		var precios []models.PreciosProdServ
		if err := tx.Where("plan_negocio_id = ?", planID).Find(&precios).Error; err != nil {
			return fmt.Errorf("loading precios for plan %d: %w", planID, err)
		}

		for _, p := range precios {

			if p.Precio == nil {
				// clear precio_calc
				if err := tx.Model(&models.PreciosProdServ{}).
					Where("id = ?", p.ID).
					Updates(map[string]interface{}{"precio_calc": nil}).Error; err != nil {
					return fmt.Errorf("clearing precio_calc for precio id %d: %w", p.ID, err)
				}

			} else {
				v := (*p.Precio) * (1.0 + precioFactor/100)
				if err := tx.Model(&models.PreciosProdServ{}).
					Where("id = ?", p.ID).
					Updates(map[string]interface{}{"precio_calc": v}).Error; err != nil {
					return fmt.Errorf("updating precio_calc for precio id %d: %w", p.ID, err)
				}
			}

			// For this product, update all costos that belong to same plan & product
			var costos []models.CostosProdServ
			if err := tx.Where("plan_negocio_id = ? AND producto_servicio_id = ?", planID, p.ProductoServicioID).
				Find(&costos).Error; err != nil {
				return fmt.Errorf("loading costos for product %d (plan %d): %w", p.ProductoServicioID, planID, err)
			}

			for _, c := range costos {
				// If there is no computed precio for this product, set costo_calc NULL
				if p.Precio== nil {
					if err := tx.Model(&models.CostosProdServ{}).
						Where("id = ?", c.ID).
						Updates(map[string]interface{}{"costo_calc": nil}).Error; err != nil {
						return fmt.Errorf("clearing costo_calc for costo id %d due missing precio: %w", c.ID, err)
					}
					continue
				}

				// determine multiplicador según categoria
				multiplicador := 0.10
				switch c.CategoriaCostoID {
				case 2:
					multiplicador = 0.12
				case 3:
					multiplicador = 0.15
				}

				costoCalc := ((*p.Precio) * multiplicador) * (1.0 + costoFactor/100)
				if err := tx.Model(&models.CostosProdServ{}).
					Where("id = ?", c.ID).
					Updates(map[string]interface{}{"costo": costoCalc}).Error; err != nil {
					return fmt.Errorf("updating costo_calc for costo id %d: %w", c.ID, err)
				}
			}
		}

		return nil
	})
}
