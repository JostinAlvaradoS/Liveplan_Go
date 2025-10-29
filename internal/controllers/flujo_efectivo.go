package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

func ListFlujoEfectivo(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var items []models.FlujoEfectivo
	if err := db.Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func ListFlujoEfectivoByPlan(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	var items []models.FlujoEfectivo
	if err := db.Where("plan_negocio_id = ?", planID).Order("anio, mes").Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type SumaAnual struct {
		Anio                         int     `json:"anio"`
		Ingresos_VentaContado        float64 `json:"ingresos_venta_contado"`
		Ingresos_CobrosVentasCredito float64 `json:"ingresos_cobros_ventas_credito"`
		Ingresos_OtrosIngresos       float64 `json:"ingresos_otros_ingresos"`
		Ingresos_Prestamos           float64 `json:"ingresos_prestamos"`
		Ingresos_AportesCapital      float64 `json:"ingresos_aportes_capital"`
		Ingresos                     float64 `json:"ingresos"`

		Egresos_ComprasCostosContado float64 `json:"egresos_compras_costos_contado"`
		Egresos_ComprasCostosCredito float64 `json:"egresos_compras_costos_credito"`
		Egresos_GastosOperacion      float64 `json:"egresos_gastos_operacion"`
		Egresos_Intereses            float64 `json:"egresos_intereses"`
		Egresos_PagosPrestamos       float64 `json:"egresos_pagos_prestamos"`
		Egresos_PagosSRI             float64 `json:"egresos_pagos_sri"`
		Egresos_PagoPTU              float64 `json:"egresos_pago_ptu"`
		Egresos                      float64 `json:"egresos"`

		Flujo_Caja      float64 `json:"flujo_caja"`
		EfectivoInicial float64 `json:"efectivo_inicial"`
		EfectivoFinal   float64 `json:"efectivo_final"`
	}

	sumas := make(map[int]*SumaAnual)
	for _, f := range items {
		s, ok := sumas[f.Anio]
		if !ok {
			s = &SumaAnual{Anio: f.Anio}
			sumas[f.Anio] = s
		}
		s.Ingresos_VentaContado += f.Ingresos_VentaContado
		s.Ingresos_CobrosVentasCredito += f.Ingresos_CobrosVentasCredito
		s.Ingresos_OtrosIngresos += f.Ingresos_OtrosIngresos
		s.Ingresos_Prestamos += f.Ingresos_Prestamos
		s.Ingresos_AportesCapital += f.Ingresos_AportesCapital
		s.Ingresos += f.Ingresos

		s.Egresos_ComprasCostosContado += f.Egresos_ComprasCostosContado
		s.Egresos_ComprasCostosCredito += f.Egresos_ComprasCostosCredito
		s.Egresos_GastosOperacion += f.Egresos_GastosOperacion
		s.Egresos_Intereses += f.Egresos_Intereses
		s.Egresos_PagosPrestamos += f.Egresos_PagosPrestamos
		s.Egresos_PagosSRI += f.Egresos_PagosSRI
		s.Egresos_PagoPTU += f.Egresos_PagoPTU
		s.Egresos += f.Egresos

		s.Flujo_Caja += f.FlujoCaja
		s.EfectivoInicial += f.EfectivoInicial
		s.EfectivoFinal += f.EfectivoFinal
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

func GetFlujoEfectivo(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.FlujoEfectivo
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

func CreateFlujoEfectivo(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var item models.FlujoEfectivo
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

func UpdateFlujoEfectivoPatch(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.FlujoEfectivo
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

func DeleteFlujoEfectivo(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&models.FlujoEfectivo{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
