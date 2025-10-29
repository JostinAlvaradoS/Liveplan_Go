package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

func ListEvaluacionProyecto(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var items []models.EvaluacionProyecto
	if err := db.Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func GetEvaluacionProyectoByPlan(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	var evaluacion models.EvaluacionProyecto
	if err := db.Where("plan_negocio_id = ?", planID).First(&evaluacion).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "EvaluacionProyecto not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(evaluacion)
}

func CreateEvaluacionProyecto(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var evaluacion models.EvaluacionProyecto
	if err := json.NewDecoder(r.Body).Decode(&evaluacion); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := db.Create(&evaluacion).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(evaluacion)
}

func GetEvaluacionProyecto(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var evaluacion models.EvaluacionProyecto
	if err := db.First(&evaluacion, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "EvaluacionProyecto not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(evaluacion)
}

func UpdateEvaluacionProyecto(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var evaluacion models.EvaluacionProyecto
	if err := db.First(&evaluacion, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "EvaluacionProyecto not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var updateData models.EvaluacionProyecto
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Actualizar campos
	evaluacion.VAN = updateData.VAN
	evaluacion.TIR = updateData.TIR
	evaluacion.TREMA = updateData.TREMA

	if err := db.Save(&evaluacion).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(evaluacion)
}

func DeleteEvaluacionProyecto(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&models.EvaluacionProyecto{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
