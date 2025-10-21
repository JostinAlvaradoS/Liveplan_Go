package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/procedimientos"
	"gorm.io/gorm"
)

func ListPreciosProdServ(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var items []models.PreciosProdServ
	query := db.Preload("PlanNegocio").Preload("ProductoServicio")
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

func ListPreciosProdServByPlan(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	var items []models.PreciosProdServ
	if err := db.Preload("ProductoServicio").Where("plan_negocio_id = ?", planID).Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func GetPreciosProdServ(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.PreciosProdServ
	if err := db.Preload("ProductoServicio").Preload("PlanNegocio").First(&item, id).Error; err != nil {
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

func CreatePreciosProdServ(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var item models.PreciosProdServ
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

func UpdatePreciosProdServPatch(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.PreciosProdServ
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
	recalc, _ := body["recalc"].(bool)
	delete(body, "id")
	delete(body, "ID")
	delete(body, "recalc")

	if err := db.Model(&item).Updates(body).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if recalc {
		planID := item.PlanNegocioID
		_ = procedimientos.Recalcular(db, planID)
	}
	json.NewEncoder(w).Encode(item)
}

func DeletePreciosProdServ(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&models.PreciosProdServ{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
