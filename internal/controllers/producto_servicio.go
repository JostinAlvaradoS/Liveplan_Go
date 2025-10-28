package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/procedimientos"
	"gorm.io/gorm"
)

func ListProductoServicios(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var items []models.ProductoServicio
	query := db.Preload("PlanNegocio")
	// optional filter by plan_id query param (plan_negocio_id)
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

func ListProductoServiciosByPlan(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	var items []models.ProductoServicio
	if err := db.Where("plan_negocio_id = ?", planID).Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func GetProductoServicio(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.ProductoServicio
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

func CreateProductoServicio(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var item models.ProductoServicio
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// create producto and related default records in a transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&item).Error; err != nil {
			return err
		}
		// create a VentaDiaria with null VentaDia
		vd := models.VentaDiaria{
			PlanNegocioID:      item.PlanNegocioID,
			ProductoServicioID: item.ID,
			VentaDia:           nil,
		}
		if err := tx.Create(&vd).Error; err != nil {
			return err
		}
		// create PreciosProdServ with nil precio and precio_calc
		pp := models.PreciosProdServ{
			PlanNegocioID:      item.PlanNegocioID,
			ProductoServicioID: item.ID,
			Precio:             nil,
			PrecioCalc:         nil,
		}
		if err := tx.Create(&pp).Error; err != nil {
			return err
		}
		// create one CostosProdServ per existing CategoriaCosto (categoria cannot be null)
		var cats []models.CategoriaCosto
		if err := tx.Find(&cats).Error; err != nil {
			return err
		}
		for _, c := range cats {
			cp := models.CostosProdServ{
				PlanNegocioID:      item.PlanNegocioID,
				ProductoServicioID: item.ID,
				CategoriaCostoID:   c.ID,
				Costo:              nil,
				CostoCalc:          nil,
			}
			if err := tx.Create(&cp).Error; err != nil {
				return err
			}
		}

		// initialize PresupuestoVenta for 5 years with null values
		for anio := 1; anio <= 5; anio++ {
			pv := models.PresupuestoVenta{
				PlanNegocioID: item.PlanNegocioID,
				ProductoID:    item.ID,
				Anio:          anio,
				Crecimiento:   nil,
				Mensual:       nil,
				Anual:         nil,
			}
			if err := tx.Create(&pv).Error; err != nil {
				return err
			}
		}

		// initialize VentasDinero for 5 years (anio 1..5) with zeros
		for anio := 1; anio <= 5; anio++ {
			vd := models.VentasDinero{
				PlanNegocioID: item.PlanNegocioID,
				ProductoID:    item.ID,
				Anio:          anio,
				Mensual:       0,
				Anual:         0,
			}
			if err := tx.Create(&vd).Error; err != nil {
				return err
			}
		}


		// initialize VentasDinero for 5 years (anio 1..5) with zeros
		for anio := 1; anio <= 5; anio++ {
			vd := models.Ventas{
				PlanNegocioID: item.PlanNegocioID,
				ProductoID:    item.ID,
				Anio:          anio,
				Venta:         0,
			}
			if err := tx.Create(&vd).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func UpdateProductoServicioPatch(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.ProductoServicio
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
		// Call Recalcular for the associated plan
		planID := item.PlanNegocioID
		_ = procedimientos.Recalcular(db, planID)
	}
	json.NewEncoder(w).Encode(item)
}

func DeleteProductoServicio(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&models.ProductoServicio{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
