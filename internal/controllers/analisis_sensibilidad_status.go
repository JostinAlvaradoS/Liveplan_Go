package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

// GetAnalisisSensibilidadStatus obtiene el status del análisis de sensibilidad para un plan
func GetAnalisisSensibilidadStatus(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	var status models.AnalisisSensibilidad_Status
	if err := db.Where("plan_negocio_id = ?", planID).First(&status).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Si no existe, crear uno con status false
			status = models.AnalisisSensibilidad_Status{
				PlanNegocioID: planID,
				Status:        false,
			}
			if err := db.Create(&status).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// UpdateAnalisisSensibilidadStatus actualiza el status del análisis de sensibilidad
func UpdateAnalisisSensibilidadStatus(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	var updateData struct {
		Status bool `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var status models.AnalisisSensibilidad_Status
	if err := db.Where("plan_negocio_id = ?", planID).First(&status).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Si no existe, crear uno nuevo
			status = models.AnalisisSensibilidad_Status{
				PlanNegocioID: planID,
				Status:        updateData.Status,
			}
			if err := db.Create(&status).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Actualizar el status existente
		if err := db.Model(&status).Update("status", updateData.Status).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
