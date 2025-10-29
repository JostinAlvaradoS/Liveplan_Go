package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

func ListBalanceGeneral(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var items []models.BalanceGeneral
	if err := db.Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func ListBalanceGeneralByPlan(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	var items []models.BalanceGeneral
	if err := db.Where("plan_negocio_id = ?", planID).Order("anio, mes").Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calcular sumas anuales por año
	type SumaAnual struct {
		Anio                          int     `json:"anio"`
		Corrientes_Efectivo           float64 `json:"corrientes_efectivo"`
		Corrientes_CuentasxCobrar     float64 `json:"corrientes_cuentasx_cobrar"`
		Corrientes_Inventarios        float64 `json:"corrientes_inventarios"`
		Corrientes_Otros              float64 `json:"corrientes_otros"`
		Corrientes_Suma               float64 `json:"corrientes_suma"`
		NoCorrientes_Suma             float64 `json:"no_corrientes_suma"`
		TotalActivo                   float64 `json:"total_activo"`
		PasivoProveedoresCortoPlazo   float64 `json:"pasivo_proveedores_corto_plazo"`
		PasivoPrestamosCortoPlazo     float64 `json:"pasivo_prestamos_corto_plazo"`
		PasivoCuentasxPagarCortoPlazo float64 `json:"pasivo_cuentasx_pagar_corto_plazo"`
		PasivoOtrosCortoPlazo         float64 `json:"pasivo_otros_corto_plazo"`
		PasivoCortoPlazo_Suma         float64 `json:"pasivo_corto_plazo_suma"`
		PasivoPrestamosLargoPlazo     float64 `json:"pasivo_prestamos_largo_plazo"`
		PasivoOtrosLargoPlazo         float64 `json:"pasivo_otros_largo_plazo"`
		PasivoLargoPlazo_Suma         float64 `json:"pasivo_largo_plazo_suma"`
		TotalPasivo                   float64 `json:"total_pasivo"`
		CapitalSocial                 float64 `json:"capital_social"`
		CapitalAdicional              float64 `json:"capital_adicional"`
		UtilidadesRetenidas           float64 `json:"utilidades_retenidas"`
		UtilidadDelEjercicio          float64 `json:"utilidad_del_ejercicio"`
		TotalCapitalContable          float64 `json:"total_capital_contable"`
	}

	sumas := make(map[int]*SumaAnual)
	for _, bg := range items {
		// Solo incluir meses 1-12 en las sumas anuales (excluir mes 0)
		if bg.Mes == 0 {
			continue
		}

		s, ok := sumas[bg.Anio]
		if !ok {
			s = &SumaAnual{Anio: bg.Anio}
			sumas[bg.Anio] = s
		}
		s.Corrientes_Efectivo += bg.Corrientes_Efectivo
		s.Corrientes_CuentasxCobrar += bg.Corrientes_CuentasxCobrar
		s.Corrientes_Inventarios += bg.Corrientes_Inventarios
		s.Corrientes_Otros += bg.Corrientes_Otros
		s.Corrientes_Suma += bg.Corrientes_Suma
		s.NoCorrientes_Suma += bg.NoCorrientes_Suma
		s.TotalActivo += bg.TotalActivo
		s.PasivoProveedoresCortoPlazo += bg.PasivoProveedoresCortoPlazo
		s.PasivoPrestamosCortoPlazo += bg.PasivoPrestamosCortoPlazo
		s.PasivoCuentasxPagarCortoPlazo += bg.PasivoCuentasxPagarCortoPlazo
		s.PasivoOtrosCortoPlazo += bg.PasivoOtrosCortoPlazo
		s.PasivoCortoPlazo_Suma += bg.PasivoCortoPlazo_Suma
		s.PasivoPrestamosLargoPlazo += bg.PasivoPrestamosLargoPlazo
		s.PasivoOtrosLargoPlazo += bg.PasivoOtrosLargoPlazo
		s.PasivoLargoPlazo_Suma += bg.PasivoLargoPlazo_Suma
		s.TotalPasivo += bg.TotalPasivo
		s.CapitalSocial += bg.CapitalSocial
		s.CapitalAdicional += bg.CapitalAdicional
		s.UtilidadesRetenidas += bg.UtilidadesRetenidas
		s.UtilidadDelEjercicio += bg.UtilidadDelEjercicio
		s.TotalCapitalContable += bg.TotalCapitalContable
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

func GetBalanceGeneral(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.BalanceGeneral
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

func CreateBalanceGeneral(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var item models.BalanceGeneral
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

func UpdateBalanceGeneralPatch(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.BalanceGeneral
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

func DeleteBalanceGeneral(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&models.BalanceGeneral{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
