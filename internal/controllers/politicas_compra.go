package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/procedimientos"
	"gorm.io/gorm"
)

func ListPoliticasCompra(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var items []models.PoliticasCompra
	if err := db.Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func ListPoliticasCompraByPlan(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	var items []models.PoliticasCompra
	if err := db.Where("plan_negocio_id = ?", planID).Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func GetPoliticasCompra(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.PoliticasCompra
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

func CreatePoliticasCompra(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var item models.PoliticasCompra
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var body map[string]interface{}
	recalc := false
	if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
		recalc, _ = body["recalc"].(bool)
	}
	if err := db.Create(&item).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if recalc {
		planID := item.PlanNegocioID
		_ = procedimientos.Recalcular(db, planID)
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func UpdatePoliticasCompraPatch(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.PoliticasCompra
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func DeletePoliticasCompra(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&models.PoliticasCompra{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
