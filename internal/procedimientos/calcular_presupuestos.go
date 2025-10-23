package procedimientos

import (
	"fmt"
	"sort"

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

		// group presupuestos by producto
		byProducto := make(map[uint][]models.PresupuestoVenta)
		for _, p := range presupuestos {
			byProducto[p.ProductoID] = append(byProducto[p.ProductoID], p)
		}

		for productoID, slice := range byProducto {
			// sort by Anio asc
			sort.Slice(slice, func(i, j int) bool { return slice[i].Anio < slice[j].Anio })

			// fetch venta diaria once per producto
			var vd models.VentaDiaria
			err := tx.Where("plan_negocio_id = ? AND producto_servicio_id = ?", planID, productoID).First(&vd).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					// for all years set Mensual/Anual to NULL
					for _, p := range slice {
						if err := tx.Model(&models.PresupuestoVenta{}).
							Where("id = ?", p.ID).
							Updates(map[string]interface{}{"mensual": nil, "anual": nil}).Error; err != nil {
							return fmt.Errorf("clearing mensual/anual for presupuesto %d: %w", p.ID, err)
						}
					}
					continue
				}
				return fmt.Errorf("finding venta_diaria for producto %d: %w", productoID, err)
			}

			if vd.VentaDia == nil {
				for _, p := range slice {
					if err := tx.Model(&models.PresupuestoVenta{}).
						Where("id = ?", p.ID).
						Updates(map[string]interface{}{"mensual": nil, "anual": nil}).Error; err != nil {
						return fmt.Errorf("clearing mensual/anual for presupuesto %d: %w", p.ID, err)
					}
				}
				continue
			}

			var prevMensual float64
			ventaDia := float64(*vd.VentaDia)

			for idx, p := range slice {
				var growth float64
				if p.Crecimiento != nil {
					growth = *p.Crecimiento / 100.0
				}

				var mensual float64
				if idx == 0 {
					// first year: use ventaDia * (1+growth) * diasxmes
					mensual = ventaDia * (1.0 + growth) 
				} else {
					// subsequent years: previous year's mensual * (1 + growth)
					mensual = prevMensual * (1.0 + growth)
				}
				anual := mensual * 12.0 * float64(diasxmes)

				// persist presupuesto
				if err := tx.Model(&models.PresupuestoVenta{}).
					Where("id = ?", p.ID).
					Updates(map[string]interface{}{"mensual": mensual, "anual": anual}).Error; err != nil {
					return fmt.Errorf("updating presupuesto %d: %w", p.ID, err)
				}

				// update/create ventas_dinero
				ventasMensual := mensual
				ventasAnual := anual
				upd := map[string]interface{}{"mensual": ventasMensual, "anual": ventasAnual}
				res := tx.Model(&models.VentasDinero{}).
					Where("plan_negocio_id = ? AND producto_id = ? AND anio = ?", planID, p.ProductoID, p.Anio).
					Updates(upd)
				if res.Error != nil {
					return fmt.Errorf("updating ventas_dinero for producto %d anio %d: %w", p.ProductoID, p.Anio, res.Error)
				}
				if res.RowsAffected == 0 {
					vdNew := models.VentasDinero{
						PlanNegocioID: planID,
						ProductoID:    p.ProductoID,
						Anio:          p.Anio,
						Mensual:       ventasMensual,
						Anual:         ventasAnual,
					}
					if err := tx.Create(&vdNew).Error; err != nil {
						return fmt.Errorf("creating ventas_dinero for producto %d anio %d: %w", p.ProductoID, p.Anio, err)
					}
				}

				prevMensual = mensual
			}
		}

		return nil
	})
}
