package procedimientos

import (
	"fmt"
	"math"

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
		} else if err != gorm.ErrRecordNotFound {
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

			// --- Update VentasDinero for the same plan/product and año ---
			// As requested: VentasDinero.mensual = PresupuestoVenta.mensual * indicadoresMacro.diasxmes
			// and VentasDinero.anual = VentasDinero.mensual * 12
			var ventasMensualFloat float64
			if mensual == 0 {
				ventasMensualFloat = 0
			} else {
				ventasMensualFloat = mensual * float64(diasxmes)
			}
			ventasMensualInt := int(math.Round(ventasMensualFloat))
			ventasAnualInt := ventasMensualInt * 12

			// try update existing VentasDinero row for this plan/product/anio
			upd := map[string]interface{}{"mensual": ventasMensualInt, "anual": ventasAnualInt}
			res := tx.Model(&models.VentasDinero{}).
				Where("plan_negocio_id = ? AND producto_id = ? AND anio = ?", planID, p.ProductoID, p.Anio).
				Updates(upd)
			if res.Error != nil {
				return fmt.Errorf("updating ventas_dinero for producto %d anio %d: %w", p.ProductoID, p.Anio, res.Error)
			}
			if res.RowsAffected == 0 {
				// create a new row if none existed
				vdNew := models.VentasDinero{
					PlanNegocioID: planID,
					ProductoID:    p.ProductoID,
					Anio:          p.Anio,
					Mensual:       ventasMensualInt,
					Anual:         ventasAnualInt,
				}
				if err := tx.Create(&vdNew).Error; err != nil {
					return fmt.Errorf("creating ventas_dinero for producto %d anio %d: %w", p.ProductoID, p.Anio, err)
				}
			}
		}

		return nil
	})
}
