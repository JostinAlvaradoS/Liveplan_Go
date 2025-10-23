package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/procedimientos"
	"gorm.io/gorm"
)

func ListCostosProdServ(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var items []models.CostosProdServ
	query := db.Preload("ProductoServicio").Preload("CategoriaCosto").Preload("PlanNegocio")
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

func ListCostosProdServByPlan(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	var items []models.CostosProdServ
	if err := db.Preload("ProductoServicio").Preload("CategoriaCosto").Where("plan_negocio_id = ?", planID).Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func GetCostosProdServ(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.CostosProdServ
	if err := db.Preload("ProductoServicio").Preload("CategoriaCosto").Preload("PlanNegocio").First(&item, id).Error; err != nil {
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

func CreateCostosProdServ(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var item models.CostosProdServ
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

func UpdateCostosProdServPatch(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.CostosProdServ
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

func DeleteCostosProdServ(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&models.CostosProdServ{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ReportCostosPorProducto genera un reporte agregado por producto para un plan:
// para cada producto devuelve los costos por categoría y la sumatoria total.
func ReportCostosPorProducto(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	var items []models.CostosProdServ
	if err := db.Preload("ProductoServicio").Preload("CategoriaCosto").Where("plan_negocio_id = ?", planID).Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type CatCost struct {
		CategoriaID   uint    `json:"categoria_id"`
		CategoriaName string  `json:"categoria_nombre"`
		Costo         float64 `json:"costo"`
	}

	type ProductReport struct {
		ProductoID   uint      `json:"producto_id"`
		ProductoName string    `json:"producto_nombre"`
		Costos       []CatCost `json:"costos_por_categoria"`
		Total        float64   `json:"total_producto"`
	}

	// Aggregate by producto
	prodMap := make(map[uint]*ProductReport)
	for _, c := range items {
		pid := c.ProductoServicioID
		pr, ok := prodMap[pid]
		if !ok {
			name := ""
			if c.ProductoServicio != nil {
				name = c.ProductoServicio.Nombre
			}
			pr = &ProductReport{ProductoID: pid, ProductoName: name, Costos: []CatCost{}, Total: 0}
			prodMap[pid] = pr
		}
		catName := ""
		if c.CategoriaCosto != nil {
			catName = c.CategoriaCosto.Nombre
		}
		costoVal := 0.0
		if c.Costo != nil {
			costoVal = *c.Costo
		}
		pr.Costos = append(pr.Costos, CatCost{CategoriaID: c.CategoriaCostoID, CategoriaName: catName, Costo: costoVal})
		pr.Total += costoVal
	}

	// Build result slice
	var result []ProductReport
	for _, pr := range prodMap {
		result = append(result, *pr)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
