package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

func ListEstadoResultados(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var items []models.EstadoResultados
	if err := db.Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func ListEstadoResultadosByPlan(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	var items []models.EstadoResultados
	if err := db.Where("plan_negocio_id = ?", planID).Order("anio, mes").Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calcular sumas anuales por año
	type SumaAnual struct {
		Anio                   int     `json:"anio"`
		Ventas                 float64 `json:"ventas"`
		CostosVentas           float64 `json:"costos_ventas"`
		UtilidadBruta          float64 `json:"utilidad_bruta"`
		GastosVentaAdm         float64 `json:"gastos_venta_adm"`
		Depreciacion           float64 `json:"depreciacion"`
		Amortizacion           float64 `json:"amortizacion"`
		UtilidadprevioIntImp   float64 `json:"utilidad_previo_int_imp"`
		GastosFinancieros      float64 `json:"gastos_financieros"`
		UtilidadAntesPTU       float64 `json:"utilidad_antes_ptu"`
		PTU                    float64 `json:"ptu"`
		UtilidadAntesImpuestos float64 `json:"utilidad_antes_impuestos"`
		ISR                    float64 `json:"isr"`
		UtilidadNeta           float64 `json:"utilidad_neta"`
	}

	sumas := make(map[int]*SumaAnual)
	for _, er := range items {
		s, ok := sumas[er.Anio]
		if !ok {
			s = &SumaAnual{Anio: er.Anio}
			sumas[er.Anio] = s
		}
		s.Ventas += er.Ventas
		s.CostosVentas += er.CostosVentas
		s.UtilidadBruta += er.UtilidadBruta
		s.GastosVentaAdm += er.GastosVentaAdm
		s.Depreciacion += er.Depreciacion
		s.Amortizacion += er.Amortizacion
		s.UtilidadprevioIntImp += er.UtilidadprevioIntImp
		s.GastosFinancieros += er.GastosFinancieros
		s.UtilidadAntesPTU += er.UtilidadAntesPTU
		s.PTU += er.PTU
		s.UtilidadAntesImpuestos += er.UtilidadAntesImpuestos
		s.ISR += er.ISR
		s.UtilidadNeta += er.UtilidadNeta
	}
	var sumasAnuales []SumaAnual
	for _, s := range sumas {
		sumasAnuales = append(sumasAnuales, *s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"items":         items,
		"sumas_anuales": sumasAnuales,
	})
}

func GetEstadoResultados(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.EstadoResultados
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

func CreateEstadoResultados(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var item models.EstadoResultados
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

func UpdateEstadoResultadosPatch(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.EstadoResultados
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

func DeleteEstadoResultados(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&models.EstadoResultados{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
