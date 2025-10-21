package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/procedimientos"
	"gorm.io/gorm"
)

func ListVariacionAnual(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var items []models.VariacionAnual
	query := db.Preload("PlanNegocio")
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

func ListVariacionAnualByPlan(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	var items []models.VariacionAnual
	if err := db.Preload("PlanNegocio").Where("plan_negocio_id = ?", planID).Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func GetVariacionAnual(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.VariacionAnual
	if err := db.Preload("PlanNegocio").First(&item, id).Error; err != nil {
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

func CreateVariacionAnual(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var item models.VariacionAnual
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

func UpdateVariacionAnualPatch(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.VariacionAnual
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
	// If any of the anio1..anio5 fields were updated, propagate the value to
	// PresupuestoVenta.crecimiento for that plan and year.
	for k, v := range body {
		var year int
		switch k {
		case "anio1":
			year = 1
		case "anio2":
			year = 2
		case "anio3":
			year = 3
		case "anio4":
			year = 4
		case "anio5":
			year = 5
		default:
			continue
		}
		// v will be float64 when decoded from JSON numbers, or nil if null
		if v == nil {
			if err := db.Model(&models.PresupuestoVenta{}).
				Where("plan_negocio_id = ? AND anio = ?", item.PlanNegocioID, year).
				Update("crecimiento", nil).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			continue
		}
		// try numeric
		if f, ok := v.(float64); ok {
			if err := db.Model(&models.PresupuestoVenta{}).
				Where("plan_negocio_id = ? AND anio = ?", item.PlanNegocioID, year).
				Update("crecimiento", f).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	json.NewEncoder(w).Encode(item)
}

func DeleteVariacionAnual(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&models.VariacionAnual{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
