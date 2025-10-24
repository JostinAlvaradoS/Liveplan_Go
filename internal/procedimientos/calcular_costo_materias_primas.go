package procedimientos

import (
	"fmt"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

// CalcularCostoMateriasPrimas calcula el costo de materias primas por producto/plan/año
// Regla:
//   - Para cada fila de VentasDinero del plan: obtener la suma de CostosProdServ donde CategoriaCostoID = 2
//     (usar siempre Costo). El costo mensual se calcula como:
//     costoMensual = VentasDinero.Mensual * sumaCostos
//     y el costo anual como costoMensual * 12.
//   - Crear o actualizar UNA fila en CostoMateriasPrimas por (plan_negocio_id, producto_id, anio)
//     guardando `costo_mensual` y `costo_anual`.
func CalcularCostoMateriasPrimas(db *gorm.DB, planID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var ventas []models.VentasDinero
		if err := tx.Where("plan_negocio_id = ?", planID).Find(&ventas).Error; err != nil {
			return fmt.Errorf("loading ventas_dinero: %w", err)
		}

		for _, v := range ventas {
			// obtener costos tipo 2 para este producto
			var costos []models.CostosProdServ
			if err := tx.Where("plan_negocio_id = ? AND producto_servicio_id = ? AND categoria_costo_id = ?", planID, v.ProductoID, 2).Find(&costos).Error; err != nil {
				return fmt.Errorf("loading costos_prodserv for producto %d: %w", v.ProductoID, err)
			}

			var sumaCostos float64
			for _, c := range costos {
				// usar siempre el campo Costo directamente (ignorar CostoCalc)
				if c.Costo != nil {
					sumaCostos += *c.Costo
				}
			}

			costoMensual := v.Mensual * sumaCostos
			costoAnual := costoMensual * 12.0

			// actualizar o crear una sola fila para el año
			upd := map[string]interface{}{"costo_mensual": costoMensual, "costo_anual": costoAnual}
			res := tx.Model(&models.CostoMateriasPrimas{}).
				Where("plan_negocio_id = ? AND producto_id = ? AND anio = ?", planID, v.ProductoID, v.Anio).
				Updates(upd)
			if res.Error != nil {
				return fmt.Errorf("updating costo_materias_primas for producto %d anio %d: %w", v.ProductoID, v.Anio, res.Error)
			}
			if res.RowsAffected == 0 {
				cm := models.CostoMateriasPrimas{
					PlanNegocioID: planID,
					ProductoID:    v.ProductoID,
					Anio:          v.Anio,
					CostoMensual:  costoMensual,
					CostoAnual:    costoAnual,
				}
				if err := tx.Create(&cm).Error; err != nil {
					return fmt.Errorf("creating costo_materias_primas for producto %d anio %d: %w", v.ProductoID, v.Anio, err)
				}
			}
		}

		return nil
	})
}
