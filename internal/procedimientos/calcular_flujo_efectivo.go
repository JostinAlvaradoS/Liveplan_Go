package procedimientos

import (
	"fmt"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

// CalcularFlujoEfectivo calcula y actualiza los valores de FlujoEfectivo para cada año y mes de un plan
func CalcularFlujoEfectivo(db *gorm.DB, planID uint) error {

	print("Inicio de procedimiento")
	var ers []models.EstadoResultados
	if err := db.Where("plan_negocio_id = ?", planID).Find(&ers).Error; err != nil {
		return err
	}
	var politicas []models.PoliticasVenta
	if err := db.Where("plan_negocio_id = ?", planID).Find(&politicas).Error; err != nil {
		return err
	}

	// Map: anio, mes -> PoliticasVenta
	pvMap := make(map[int]map[int]models.PoliticasVenta)
	for _, pv := range politicas {
		if _, ok := pvMap[pv.Anio]; !ok {
			pvMap[pv.Anio] = make(map[int]models.PoliticasVenta)
		}
		pvMap[pv.Anio][pv.Mes] = pv
	}

	for _, er := range ers {
		anio := er.Anio
		mes := er.Mes
		ventas := er.Ventas

		// Solo procesar mes 0 si es el primer año
		if mes == 0 && anio != 1 {
			continue
		}

		pv, ok := pvMap[anio][mes]
		if !ok {
			fmt.Printf("[FlujoEfectivo] PoliticasVenta no encontradas para anio=%d mes=%d\n", anio, mes)
		}
		contado := pv.PorcentajeContado

		anioAnt := anio
		mesAnt := mes - 1
		if mesAnt == 0 {
			mesAnt = 12
			anioAnt = anio - 1
		}
		ventasAnt := 0.0
		for _, erAnt := range ers {
			if erAnt.Anio == anioAnt && erAnt.Mes == mesAnt {
				ventasAnt = erAnt.Ventas
				break
			}
		}
		pvAnt, okAnt := pvMap[anioAnt][mesAnt]
		if !okAnt {
			fmt.Printf("[FlujoEfectivo] PoliticasVenta anterior no encontrada para anio=%d mes=%d\n", anioAnt, mesAnt)
		}
		creditoAnt := pvAnt.PorcentajeCredito

		ingContado := ventas * contado / 100.0
		ingCredito := ventasAnt * creditoAnt / 100.0

		fmt.Printf("[FlujoEfectivo] anio=%d mes=%d ventas=%.2f contado=%.2f ingContado=%.2f ventasAnt=%.2f creditoAnt=%.2f ingCredito=%.2f\n",
			anio, mes, ventas, contado, ingContado, ventasAnt, creditoAnt, ingCredito)

		var flujo models.FlujoEfectivo
		if err := db.Where("plan_negocio_id = ? AND anio = ? AND mes = ?", planID, anio, mes).First(&flujo).Error; err != nil {
			flujo = models.FlujoEfectivo{
				PlanNegocioID: planID,
				Anio:          anio,
				Mes:           mes,
			}
		}
		flujo.Ingresos_VentaContado = ingContado
		flujo.Ingresos_CobrosVentasCredito = ingCredito
		if err := db.Save(&flujo).Error; err != nil {
			return err
		}
	}
	return nil
}
