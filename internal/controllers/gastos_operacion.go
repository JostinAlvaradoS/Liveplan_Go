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
	// Obtener gastos operacion base
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

	// Armar respuesta por año y mes
	type MesResult struct {
		Anio         int                      `json:"anio"`
		Mes          int                      `json:"mes"`
		Gastos       []models.GastosOperacion `json:"gastos_operacion"`
		Subtotal     float64                  `json:"subtotal"`
		Intereses    float64                  `json:"intereses_prestamo"`
		Depreciacion float64                  `json:"depreciacion"`
		Amortizacion float64                  `json:"amortizacion"`
		Total        float64                  `json:"total"`
	}

	// GastosOperacion no tiene mes/anio, se asigna igual a todos los meses
	gastosTotal := 0.0
	for _, gope := range gastos {
		gastosTotal += gope.Costo
	}

	// Mapas para sumar por año/mes
	interesesPorAnioMes := make(map[int]map[int]float64)
	for _, c := range cuotas {
		if _, ok := interesesPorAnioMes[c.Anio]; !ok {
			interesesPorAnioMes[c.Anio] = make(map[int]float64)
		}
		interesesPorAnioMes[c.Anio][c.Mes] += c.Interes
	}

	depreciacionPorAnioMes := make(map[int]map[int]float64)
	amortizacionPorAnioMes := make(map[int]map[int]float64)
	for _, dep := range depreciaciones {
		tipo := 0
		if dep.DetalleInversion != nil {
			tipo = int(dep.DetalleInversion.TipoID)
		}
		val := 0.0
		if dep.DepreciacionMensual != nil {
			val = *dep.DepreciacionMensual
		}
		anios := []int{1, 2, 3, 4, 5}
		for _, anio := range anios {
			if _, ok := depreciacionPorAnioMes[anio]; !ok {
				depreciacionPorAnioMes[anio] = make(map[int]float64)
			}
			if _, ok := amortizacionPorAnioMes[anio]; !ok {
				amortizacionPorAnioMes[anio] = make(map[int]float64)
			}
			for mes := 1; mes <= 12; mes++ {
				if tipo == 1 {
					depreciacionPorAnioMes[anio][mes] += val
				} else if tipo == 2 {
					amortizacionPorAnioMes[anio][mes] += val
				}
			}
		}
	}

	var result []MesResult
	for anio := 1; anio <= 5; anio++ {
		for mes := 1; mes <= 12; mes++ {
			subtotal := gastosTotal
			intereses := 0.0
			if m, ok := interesesPorAnioMes[anio]; ok {
				intereses = m[mes]
			}
			depreciacion := 0.0
			if m, ok := depreciacionPorAnioMes[anio]; ok {
				depreciacion = m[mes]
			}
			amortizacion := 0.0
			if m, ok := amortizacionPorAnioMes[anio]; ok {
				amortizacion = m[mes]
			}
			total := subtotal + intereses + depreciacion + amortizacion
			result = append(result, MesResult{
				Anio:         anio,
				Mes:          mes,
				Gastos:       gastos, // todos los gastos operacion
				Subtotal:     subtotal,
				Intereses:    intereses,
				Depreciacion: depreciacion,
				Amortizacion: amortizacion,
				Total:        total,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func GetGastosOperacion(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.GastosOperacion
	if err := db.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func CreateGastosOperacion(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var item models.GastosOperacion
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := db.Create(&item).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
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
