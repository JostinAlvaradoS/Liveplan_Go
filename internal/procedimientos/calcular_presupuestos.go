package procedimientos

import (
	"fmt"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

// CalcularPresupuestos recalcula Mensual y Anual en PresupuestoVenta para un plan.
// Reglas:
// - Se toma DiasxMes desde IndicadoresMacro del plan; si no existe, se usa 30.
// - Para cada fila de PresupuestoVenta del plan:
//   - Se busca VentaDiaria del mismo plan y producto. Si no existe o VentaDia es NULL, Mensual y Anual se dejan NULL.
//   - growth = Crecimiento (porcentaje) / 100.0 (ej: 5 -> 0.05). Si Crecimiento es NULL se asume 0.
//   - mensual = VentaDia * (1 + growth) * diasxmes
//   - anual = mensual * 12
func CalcularPresupuestos(db *gorm.DB, planID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// get diasxmes
		var im models.IndicadoresMacro
		diasxmes := 30
		if err := tx.Where("plan_negocio_id = ?", planID).First(&im).Error; err == nil {
			if im.DiasxMes > 0 {
				diasxmes = im.DiasxMes
			}
		} else if err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("loading indicadores_macro: %w", err)
		}

		var presupuestos []models.PresupuestoVenta
		if err := tx.Where("plan_negocio_id = ?", planID).Find(&presupuestos).Error; err != nil {
			return fmt.Errorf("loading presupuestos: %w", err)
		}

		for _, p := range presupuestos {
			// find venta diaria for this product
			var vd models.VentaDiaria
			err := tx.Where("plan_negocio_id = ? AND producto_servicio_id = ?", planID, p.ProductoID).First(&vd).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					// set Mensual/Anual to NULL
					if err := tx.Model(&models.PresupuestoVenta{}).
						Where("id = ?", p.ID).
						Updates(map[string]interface{}{"mensual": nil, "anual": nil}).Error; err != nil {
						return fmt.Errorf("clearing mensual/anual for presupuesto %d: %w", p.ID, err)
					}
					continue
				}
				return fmt.Errorf("finding venta_diaria for producto %d: %w", p.ProductoID, err)
			}

			if vd.VentaDia == nil {
				if err := tx.Model(&models.PresupuestoVenta{}).
					Where("id = ?", p.ID).
					Updates(map[string]interface{}{"mensual": nil, "anual": nil}).Error; err != nil {
					return fmt.Errorf("clearing mensual/anual for presupuesto %d: %w", p.ID, err)
				}
				continue
			}

			// compute
			var growth float64
			if p.Crecimiento != nil {
				growth = *p.Crecimiento / 100.0
			}
			ventaDia := float64(*vd.VentaDia)
			mensual := ventaDia * (1.0 + growth) * float64(diasxmes)
			anual := mensual * 12.0

			if err := tx.Model(&models.PresupuestoVenta{}).
				Where("id = ?", p.ID).
				Updates(map[string]interface{}{"mensual": mensual, "anual": anual}).Error; err != nil {
				return fmt.Errorf("updating presupuesto %d: %w", p.ID, err)
			}
		}

		return nil
	})
}
