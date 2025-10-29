package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

func ListConceptosEvaluacion(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var items []models.ConceptosEvaluacion
	if err := db.Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func ListConceptosEvaluacionByPlan(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	var items []models.ConceptosEvaluacion
	if err := db.Where("plan_negocio_id = ?", planID).Order("anio").Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Estructura de respuesta con los datos organizados
	type ResponseData struct {
		PlanNegocioID uint                         `json:"plan_negocio_id"`
		Conceptos     []models.ConceptosEvaluacion `json:"conceptos"`
	}

	response := ResponseData{
		PlanNegocioID: planID,
		Conceptos:     items,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CreateConceptoEvaluacion(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var concepto models.ConceptosEvaluacion
	if err := json.NewDecoder(r.Body).Decode(&concepto); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := db.Create(&concepto).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(concepto)
}

func GetConceptoEvaluacion(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var concepto models.ConceptosEvaluacion
	if err := db.First(&concepto, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "ConceptoEvaluacion not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(concepto)
}

func UpdateConceptoEvaluacion(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var concepto models.ConceptosEvaluacion
	if err := db.First(&concepto, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "ConceptoEvaluacion not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var updateData models.ConceptosEvaluacion
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Actualizar campos
	concepto.FlujoEfectivoNominal = updateData.FlujoEfectivoNominal
	concepto.ValorRescate = updateData.ValorRescate
	concepto.TotalFlujoEfectivo = updateData.TotalFlujoEfectivo
	concepto.ValorActualFlujosFuturos = updateData.ValorActualFlujosFuturos

	if err := db.Save(&concepto).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(concepto)
}

func DeleteConceptoEvaluacion(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&models.ConceptosEvaluacion{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
