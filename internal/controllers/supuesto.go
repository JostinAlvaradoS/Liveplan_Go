package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

func ListSupuestos(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var items []models.Supuesto
	query := db
	if pid := r.URL.Query().Get("plan_id"); pid != "" {
		id, err := strconv.Atoi(pid)
		if err != nil {
			http.Error(w, "invalid plan_id", http.StatusBadRequest)
			return
		}
		query = query.Where("plan_negocio_id = ?", id)
	}
	if err := query.Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func GetSupuesto(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.Supuesto
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

func CreateSupuesto(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var item models.Supuesto
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

func UpdateSupuestoPatch(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.Supuesto
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

func DeleteSupuesto(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&models.Supuesto{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func ListSupuestosByPlan(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	var items []models.Supuesto
	if err := db.Where("plan_negocio_id = ?", planID).Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
