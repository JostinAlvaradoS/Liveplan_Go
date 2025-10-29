package procedimientos

import (
    "fmt"
    "log"

    "github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
    "gorm.io/gorm"
)

// CalcularCostosVentas calcula la tabla CostosVentas para un plan.
// Para cada registro de VentasDinero (plan, producto, anio) toma el Mensual
// (que puede variar por año) y lo multiplica por la suma de los costos
// asociados al producto (CostosProdServ.CostoCalc | Costo). Ese valor se guarda
// como costo mensual y se repite en los 12 meses (mes 1..12) del año.
func CalcularCostosVentas(db *gorm.DB, planID uint) error {
    // Cargar todos los registros VentasDinero del plan
    var ventasDin []models.VentasDinero
    if err := db.Where("plan_negocio_id = ?", planID).Find(&ventasDin).Error; err != nil {
        return fmt.Errorf("obtener VentasDinero: %w", err)
    }

    // Sumar costos por producto (CostosProdServ) para el plan
    var costos []models.CostosProdServ
    if err := db.Where("plan_negocio_id = ?", planID).Find(&costos).Error; err != nil {
        return fmt.Errorf("obtener CostosProdServ: %w", err)
    }
    // mapa productoID -> suma de costos
    costosPorProducto := make(map[uint]float64)
    for _, c := range costos {
        var val float64
        if c.CostoCalc != nil {
            val = *c.CostoCalc
        } else if c.Costo != nil {
            val = *c.Costo
        } else {
            val = 0
        }
        costosPorProducto[c.ProductoServicioID] += val
    }

    // Iterar cada VentasDinero (cada producto por año) y crear/actualizar 12 meses
    for _, vd := range ventasDin {
        sumaCostosProducto := costosPorProducto[vd.ProductoID]
        // si no existen costos asociados, asumimos 0
        if sumaCostosProducto == 0 {
            // no es error, solo resultado 0
        }

        // Costo mensual total asociado a las ventas = Mensual * sumaCostosProducto
        costoMensual := vd.Mensual * sumaCostosProducto

        for mes := 1; mes <= 12; mes++ {
            var cv models.CostosVentas
            q := db.Where("plan_negocio_id = ? AND producto_id = ? AND anio = ? AND mes = ?", planID, vd.ProductoID, vd.Anio, mes).First(&cv)
            if q.Error == nil {
                cv.Mensual = costoMensual
                if err := db.Save(&cv).Error; err != nil {
                    log.Printf("CalcularCostosVentas: error actualizando costo ventas P:%d Prod:%d A:%d M:%d: %v", planID, vd.ProductoID, vd.Anio, mes, err)
                    return fmt.Errorf("actualizar costos_ventas: %w", err)
                }
            } else if q.Error == gorm.ErrRecordNotFound {
                newCv := models.CostosVentas{
                    PlanNegocioID: planID,
                    ProductoID:    vd.ProductoID,
                    Anio:          vd.Anio,
                    Mes:           mes,
                    Mensual:      costoMensual,
					Anual: costoMensual * 12,
                }
                if err := db.Create(&newCv).Error; err != nil {
                    log.Printf("CalcularCostosVentas: error creando costo ventas P:%d Prod:%d A:%d M:%d: %v", planID, vd.ProductoID, vd.Anio, mes, err)
                    return fmt.Errorf("crear costos_ventas: %w", err)
                }
            } else {
                return fmt.Errorf("consultar costos_ventas: %w", q.Error)
            }
        }
    }

    return nil
}
