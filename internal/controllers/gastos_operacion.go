package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

func ListGastosOperacion(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var items []models.GastosOperacion
	if err := db.Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func ListGastosOperacionByPlan(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	// Obtener gastos operacion base (ahora con anual por año)
	var gastos []models.GastosOperacion
	if err := db.Where("plan_negocio_id = ?", planID).Find(&gastos).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Obtener cuotas de préstamo
	var cuotas []models.PrestamoCuotas
	if err := db.Where("plan_negocio_id = ?", planID).Find(&cuotas).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Obtener depreciaciones
	var depreciaciones []models.Depreciacion
	if err := db.Where("plan_negocio_id = ?", planID).Preload("DetalleInversion").Find(&depreciaciones).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Mapas para sumar por año
	var interesesPorAnio map[int]float64
	interesesPorAnio = make(map[int]float64)
	// Declarar mapas para depreciaciones y amortizaciones mensuales
	var depreciacionPorMes map[int]map[int]float64 = make(map[int]map[int]float64)
	var amortizacionPorMes map[int]map[int]float64 = make(map[int]map[int]float64)
	for _, c := range cuotas {
		interesesPorAnio[c.Anio] += c.Interes
	}

	var depreciacionPorAnio map[int]float64
	var amortizacionPorAnio map[int]float64
	depreciacionPorAnio = make(map[int]float64)
	amortizacionPorAnio = make(map[int]float64)
	depreciacionPorMes = make(map[int]map[int]float64)
	amortizacionPorMes = make(map[int]map[int]float64)
	for _, dep := range depreciaciones {
		tipo := 0
		if dep.DetalleInversion != nil {
			tipo = int(dep.DetalleInversion.TipoID)
		}
		val := 0.0
		if dep.DepreciacionMensual != nil {
			val = *dep.DepreciacionMensual
		}
		for anio := 1; anio <= 5; anio++ {
			for mes := 1; mes <= 12; mes++ {
				if tipo == 1 {
					if _, ok := depreciacionPorMes[anio]; !ok {
						depreciacionPorMes[anio] = make(map[int]float64)
					}
					depreciacionPorMes[anio][mes] += val
				} else if tipo == 2 {
					if _, ok := amortizacionPorMes[anio]; !ok {
						amortizacionPorMes[anio] = make(map[int]float64)
					}
					amortizacionPorMes[anio][mes] += val
				}
			}
		}
	}
	// Sumar las mensuales para obtener el anual
	for anio := 1; anio <= 5; anio++ {
		for mes := 1; mes <= 12; mes++ {
			depreciacionPorAnio[anio] += depreciacionPorMes[anio][mes]
			amortizacionPorAnio[anio] += amortizacionPorMes[anio][mes]
		}
	}

	// Mapas para sumar por año y mes
	interesesPorAnio = make(map[int]float64)
	interesesPorMes := make(map[int]map[int]float64)
	for _, c := range cuotas {
		interesesPorAnio[c.Anio] += c.Interes
		if _, ok := interesesPorMes[c.Anio]; !ok {
			interesesPorMes[c.Anio] = make(map[int]float64)
		}
		interesesPorMes[c.Anio][c.Mes] += c.Interes
	}

	depreciacionPorAnio = make(map[int]float64)
	amortizacionPorAnio = make(map[int]float64)
	for _, dep := range depreciaciones {
		tipo := 0
		if dep.DetalleInversion != nil {
			tipo = int(dep.DetalleInversion.TipoID)
		}
		val := 0.0
		if dep.DepreciacionMensual != nil {
			val = *dep.DepreciacionMensual
		}
		for anio := 1; anio <= 5; anio++ {
			for mes := 1; mes <= 12; mes++ {
				if tipo == 1 {
					depreciacionPorAnio[anio] += val
					if _, ok := depreciacionPorMes[anio]; !ok {
						depreciacionPorMes[anio] = make(map[int]float64)
					}
					depreciacionPorMes[anio][mes] += val
				} else if tipo == 2 {
					amortizacionPorAnio[anio] += val
					if _, ok := amortizacionPorMes[anio]; !ok {
						amortizacionPorMes[anio] = make(map[int]float64)
					}
					amortizacionPorMes[anio][mes] += val
				}
			}
		}
	}
	// ...eliminado: ya no se multiplica por 12, el anual es la suma de las mensuales...

	// Sumar gastos operacion anual y mensual por año
	gastosOperacionPorAnio := make(map[int]float64)
	gastosOperacionPorMes := make(map[int]map[int]float64)
	for _, gope := range gastos {
		for anio := 1; anio <= 5; anio++ {
			gastosOperacionPorAnio[anio] += gope.Anual
			if _, ok := gastosOperacionPorMes[anio]; !ok {
				gastosOperacionPorMes[anio] = make(map[int]float64)
			}
			for mes := 1; mes <= 12; mes++ {
				gastosOperacionPorMes[anio][mes] += gope.Mensual
			}
		}
	}

	// Armar reporte por año
	type ReporteAnual struct {
		Anio            int     `json:"anio"`
		GastosOperacion float64 `json:"gastos_operacion_anual"`
		Intereses       float64 `json:"intereses_prestamo_anual"`
		Depreciacion    float64 `json:"depreciacion_anual"`
		Amortizacion    float64 `json:"amortizacion_anual"`
		Total           float64 `json:"total_anual"`
	}
	type ReporteMensual struct {
		Anio            int     `json:"anio"`
		Mes             int     `json:"mes"`
		GastosOperacion float64 `json:"gastos_operacion_mensual"`
		Intereses       float64 `json:"intereses_prestamo_mensual"`
		Depreciacion    float64 `json:"depreciacion_mensual"`
		Amortizacion    float64 `json:"amortizacion_mensual"`
		Total           float64 `json:"total_mensual"`
	}
	var reporteAnual []ReporteAnual
	var reporteMensual []ReporteMensual
	for anio := 1; anio <= 5; anio++ {
		totalAnual := gastosOperacionPorAnio[anio] + interesesPorAnio[anio] + depreciacionPorAnio[anio] + amortizacionPorAnio[anio]
		reporteAnual = append(reporteAnual, ReporteAnual{
			Anio:            anio,
			GastosOperacion: gastosOperacionPorAnio[anio],
			Intereses:       interesesPorAnio[anio],
			Depreciacion:    depreciacionPorAnio[anio],
			Amortizacion:    amortizacionPorAnio[anio],
			Total:           totalAnual,
		})
		for mes := 1; mes <= 12; mes++ {
			totalMensual := gastosOperacionPorMes[anio][mes] + interesesPorMes[anio][mes] + depreciacionPorMes[anio][mes] + amortizacionPorMes[anio][mes]
			reporteMensual = append(reporteMensual, ReporteMensual{
				Anio:            anio,
				Mes:             mes,
				GastosOperacion: gastosOperacionPorMes[anio][mes],
				Intereses:       interesesPorMes[anio][mes],
				Depreciacion:    depreciacionPorMes[anio][mes],
				Amortizacion:    amortizacionPorMes[anio][mes],
				Total:           totalMensual,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"anuales":   reporteAnual,
		"mensuales": reporteMensual,
		"gastos":    gastos,
	})
}

func UpdateGastosOperacionPatch(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.GastosOperacion
	if err := db.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	delete(body, "id")
	delete(body, "ID")
	if err := db.Model(&item).Updates(body).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(item)
}

func DeleteGastosOperacion(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&models.GastosOperacion{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
