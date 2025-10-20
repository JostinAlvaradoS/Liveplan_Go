package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/procedimientos"
	"gorm.io/gorm"
)

func ListDetalles(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var items []models.DetalleInversionInicial
	// Preload Tipo and the related Inversion and its PlanNegocio to return nested data
	query := db.Preload("Tipo").Preload("Inversion").Preload("Inversion.PlanNegocio")
	// optional filter by plan_id query param
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

func GetDetalle(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.DetalleInversionInicial
	if err := db.Preload("Tipo").Preload("Inversion").Preload("Inversion.PlanNegocio").First(&item, id).Error; err != nil {
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

func CreateDetalle(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var item models.DetalleInversionInicial
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := db.Create(&item).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// create initial Depreciacion record for this detalle with vida_util = 0 and other fields NULL
	dep := models.Depreciacion{
		PlanNegocioID:       item.PlanNegocioID,
		DetalleInversionID:  item.ID,
		DepreciacionMensual: nil,
		DepreciacionAnio1:   nil,
		DepreciacionAnio2:   nil,
		DepreciacionAnio3:   nil,
		DepreciacionAnio4:   nil,
		DepreciacionAnio5:   nil,
		ValorRescate:        nil,
	}
	// ignore error if it fails; creation should be best-effort but we log to response if necessary
	_ = db.Create(&dep).Error
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func UpdateDetallePatch(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.DetalleInversionInicial
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
	procedimientos.CalcularDepreciaciones(db, item.PlanNegocioID)
	json.NewEncoder(w).Encode(item)
}

func DeleteDetalle(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&models.DetalleInversionInicial{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func ListDetallesByPlan(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	var items []models.DetalleInversionInicial
	if err := db.Preload("Tipo").Preload("Inversion").Where("plan_negocio_id = ?", planID).Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
