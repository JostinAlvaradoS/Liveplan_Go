package controllers

import (
	"encoding/json"
	"net/http"
	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

func ListCostosVentaByPlan(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	var items []models.CostosVentas
	if err := db.Where("plan_negocio_id = ?", planID).Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
