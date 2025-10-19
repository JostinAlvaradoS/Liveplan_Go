package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

// PlanNegocio CRUD
func ListPlanNegocios(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var items []models.PlanNegocio
	if err := db.Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func GetPlanNegocio(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.PlanNegocio
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

func CreatePlanNegocio(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var item models.PlanNegocio
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

// Patch (partial) update for PlanNegocio
func UpdatePlanNegocioPatch(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.PlanNegocio
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

func DeletePlanNegocio(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&models.PlanNegocio{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
