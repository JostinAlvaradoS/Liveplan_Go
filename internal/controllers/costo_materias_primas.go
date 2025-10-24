package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/procedimientos"
	"gorm.io/gorm"
)

func ListCostoMateriasPrimas(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var items []models.CostoMateriasPrimas
	query := db.Preload("Producto").Preload("PlanNegocio")
	if pid := r.URL.Query().Get("plan_id"); pid != "" {
		// reuse existing ParseUintFromPath? keep simple: filter by plan_id string -> GORM accepts string
		query = query.Where("plan_negocio_id = ?", pid)
	}
	if err := query.Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func ListCostoMateriasPrimasByPlan(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	var items []models.CostoMateriasPrimas
	if err := db.Preload("Producto").Where("plan_negocio_id = ?", planID).Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func GetCostoMateriasPrimas(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.CostoMateriasPrimas
	if err := db.Preload("Producto").Preload("PlanNegocio").First(&item, id).Error; err != nil {
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

func CreateCostoMateriasPrimas(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var item models.CostoMateriasPrimas
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

func UpdateCostoMateriasPrimasPatch(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.CostoMateriasPrimas
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

func DeleteCostoMateriasPrimas(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&models.CostoMateriasPrimas{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
