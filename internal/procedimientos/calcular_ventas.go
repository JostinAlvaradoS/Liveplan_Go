package procedimientos

import (
	"fmt"
	"log"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

// CalcularVentas calcula las filas de la tabla Ventas para un plan dado.
// Para cada registro de VentasDinero (PlanNegocioID, ProductoID, Anio)
// calcula Venta = VentasDinero.Mensual * PreciosProdServ.PrecioCalc
// y hace upsert en la tabla Ventas (por PlanNegocioID, ProductoID, Anio).
func CalcularVentas(db *gorm.DB, planID uint) error {
	var ventasDin []models.VentasDinero
	if err := db.Where("plan_negocio_id = ?", planID).Find(&ventasDin).Error; err != nil {
		return fmt.Errorf("obtener VentasDinero: %w", err)
	}
	// Cargar precios una vez en memoria por producto para evitar consultas repetidas
	var precios []models.PreciosProdServ
	if err := db.Where("plan_negocio_id = ?", planID).Find(&precios).Error; err != nil {
		return fmt.Errorf("obtener precios prodserv: %w", err)
	}
	precioMap := make(map[uint]float64)
	for _, p := range precios {
		if p.PrecioCalc != nil {
			precioMap[p.ProductoServicioID] = *p.PrecioCalc
		} else if p.Precio != nil {
			precioMap[p.ProductoServicioID] = *p.Precio
		} else {
			precioMap[p.ProductoServicioID] = 0
		}
	}

	// Iterar por cada registro de VentasDinero (producto + anio). Mensual puede cambiar por anio.
	for _, vd := range ventasDin {
		precioCalc, ok := precioMap[vd.ProductoID]
		if !ok {
			// si no existe precio para ese producto, asumimos 0
			precioCalc = 0
		}

		// ventasDinero.Mensual corresponde al valor mensual para ese anio
		venta := vd.Mensual * precioCalc
		var v models.Ventas
		q := db.Where("plan_negocio_id = ? AND producto_id = ? AND anio = ?", planID, vd.ProductoID, vd.Anio).First(&v)
		if q.Error == nil {
			// actualizar
			v.Venta = venta
			if err := db.Save(&v).Error; err != nil {
				log.Printf("CalcularVentas: error actualizando ventas para producto %d anio %d: %v", vd.ProductoID, vd.Anio, err)
				return fmt.Errorf("actualizar ventas: %w", err)
			}
		} else if q.Error == gorm.ErrRecordNotFound {
			// crear nuevo registro
			newV := models.Ventas{
				PlanNegocioID: planID,
				ProductoID:    vd.ProductoID,
				Anio:          vd.Anio,
				Venta:         venta,
			}
			if err := db.Create(&newV).Error; err != nil {
				log.Printf("CalcularVentas: error creando ventas para producto %d anio %d: %v", vd.ProductoID, vd.Anio, err)
				return fmt.Errorf("crear ventas: %w", err)
			}
		} else {
			log.Printf("CalcularVentas: error consultando ventas existentes: %v", q.Error)
			return fmt.Errorf("consultar ventas existentes: %w", q.Error)
		}
	}

	return nil
}
