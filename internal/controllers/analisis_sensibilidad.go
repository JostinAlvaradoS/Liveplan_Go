package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

func ListAnalisisSensibilidad(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var items []models.AnalisisSensibilidad
	if err := db.Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func ListAnalisisSensibilidadByPlan(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	var items []models.AnalisisSensibilidad
	if err := db.Where("plan_negocio_id = ?", planID).Order("volumen, costo").Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type ResponseData struct {
		PlanNegocioID uint                          `json:"plan_negocio_id"`
		Analisis      []models.AnalisisSensibilidad `json:"analisis"`
	}

	response := ResponseData{
		PlanNegocioID: planID,
		Analisis:      items,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CreateAnalisisSensibilidad(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var a models.AnalisisSensibilidad
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := db.Create(&a).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(a)
}

func GetAnalisisSensibilidad(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var a models.AnalisisSensibilidad
	if err := db.First(&a, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "AnalisisSensibilidad not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a)
}

func UpdateAnalisisSensibilidad(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var a models.AnalisisSensibilidad
	if err := db.First(&a, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "AnalisisSensibilidad not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var updateData models.AnalisisSensibilidad
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	a.Volumen = updateData.Volumen
	a.Costo = updateData.Costo
	a.Valor = updateData.Valor

	if err := db.Save(&a).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a)
}

func DeleteAnalisisSensibilidad(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&models.AnalisisSensibilidad{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
