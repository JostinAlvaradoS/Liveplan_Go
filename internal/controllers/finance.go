package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

// PlanNegocio CRUD
func ListPlanNegocios(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var items []models.PlanNegocio
	if err := db.Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func GetPlanNegocio(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.PlanNegocio
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

func CreatePlanNegocio(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var item models.PlanNegocio
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

// Patch (partial) update for PlanNegocio
func UpdatePlanNegocioPatch(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.PlanNegocio
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

func DeletePlanNegocio(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&models.PlanNegocio{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// TipoInversionInicial CRUD
func ListTipos(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var items []models.TipoInversionInicial
	if err := db.Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func CreateTipo(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var item models.TipoInversionInicial
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

// InversionInicial CRUD
func ListInversiones(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var items []models.InversionInicial
	// Preload both Tipo and PlanNegocio so nested data is returned
	query := db.Preload("Tipo").Preload("PlanNegocio")
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

func CreateInversion(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var item models.InversionInicial
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

// Patch (partial) update for InversionInicial
func UpdateInversionPatch(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.InversionInicial
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

// DetalleInversionInicial CRUD
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

// Get single InversionInicial by id with preloaded associations
func GetInversion(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.InversionInicial
	if err := db.Preload("Tipo").First(&item, id).Error; err != nil {
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

// Get single DetalleInversionInicial by id with nested preloads
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

// ListInversionesByPlan lists inversiones filtered by plan_negocio_id (planID)
func ListInversionesByPlan(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	var items []models.InversionInicial
	if err := db.Preload("Tipo").Preload("PlanNegocio").Where("plan_negocio_id = ?", planID).Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// ListDetallesByPlan lists detalles filtered by plan_negocio_id (planID)
func ListDetallesByPlan(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	var items []models.DetalleInversionInicial
	if err := db.Preload("Tipo").Preload("Inversion").Where("plan_negocio_id = ?", planID).Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// DeleteInversion deletes an inversion by id
func DeleteInversion(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&models.InversionInicial{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteDetalle deletes a detalle by id
func DeleteDetalle(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&models.DetalleInversionInicial{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ProductoServicio CRUD
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

// ListProductoServiciosByPlan lists producto_servicio filtered by plan_negocio_id (planID)
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
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// Patch (partial) update for ProductoServicio
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
	delete(body, "id")
	delete(body, "ID")

	if err := db.Model(&item).Updates(body).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// Note: full-update PUT handlers removed; use partial PATCH handlers instead

// Patch (partial) update for DetalleInversionInicial
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
	json.NewEncoder(w).Encode(item)
}

// Supuesto CRUD
func ListSupuestos(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var items []models.Supuesto
	query := db
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

func GetSupuesto(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.Supuesto
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

func CreateSupuesto(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var item models.Supuesto
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

func UpdateSupuestoPatch(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.Supuesto
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

func DeleteSupuesto(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&models.Supuesto{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListSupuestosByPlan lists supuestos filtered by plan_negocio_id (planID)
func ListSupuestosByPlan(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	var items []models.Supuesto
	if err := db.Where("plan_negocio_id = ?", planID).Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// VentaDiaria CRUD
func ListVentaDiarias(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var items []models.VentaDiaria
	query := db.Preload("ProductoServicio").Preload("PlanNegocio")
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

func ListVentaDiariasByPlan(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	var items []models.VentaDiaria
	if err := db.Preload("ProductoServicio").Where("plan_negocio_id = ?", planID).Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func GetVentaDiaria(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.VentaDiaria
	if err := db.Preload("ProductoServicio").Preload("PlanNegocio").First(&item, id).Error; err != nil {
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

func CreateVentaDiaria(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var item models.VentaDiaria
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

func UpdateVentaDiariaPatch(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.VentaDiaria
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

func DeleteVentaDiaria(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&models.VentaDiaria{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// VariablesDeSensibilidad CRUD
func ListVariablesDeSensibilidad(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var items []models.VariablesDeSensibilidad
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

func ListVariablesDeSensibilidadByPlan(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	var items []models.VariablesDeSensibilidad
	if err := db.Where("plan_negocio_id = ?", planID).Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func GetVariablesDeSensibilidad(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.VariablesDeSensibilidad
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

func CreateVariablesDeSensibilidad(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var item models.VariablesDeSensibilidad
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

func UpdateVariablesDeSensibilidadPatch(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.VariablesDeSensibilidad
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

func DeleteVariablesDeSensibilidad(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&models.VariablesDeSensibilidad{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// VariacionAnual CRUD
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
	delete(body, "id")
	delete(body, "ID")
	if err := db.Model(&item).Updates(body).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

// PreciosProdServ CRUD
func ListPreciosProdServ(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var items []models.PreciosProdServ
	query := db.Preload("PlanNegocio").Preload("ProductoServicio")
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

func ListPreciosProdServByPlan(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	var items []models.PreciosProdServ
	if err := db.Preload("ProductoServicio").Where("plan_negocio_id = ?", planID).Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func GetPreciosProdServ(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.PreciosProdServ
	if err := db.Preload("ProductoServicio").Preload("PlanNegocio").First(&item, id).Error; err != nil {
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

func CreatePreciosProdServ(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var item models.PreciosProdServ
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

func UpdatePreciosProdServPatch(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.PreciosProdServ
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

func DeletePreciosProdServ(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&models.PreciosProdServ{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// CategoriaCosto CRUD (catalog)
func ListCategoriaCosto(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var items []models.CategoriaCosto
	if err := db.Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func GetCategoriaCosto(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.CategoriaCosto
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

func CreateCategoriaCosto(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var item models.CategoriaCosto
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

func UpdateCategoriaCostoPatch(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	var item models.CategoriaCosto
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

func DeleteCategoriaCosto(db *gorm.DB, w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&models.CategoriaCosto{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// CostosProdServ CRUD
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
	delete(body, "id")
	delete(body, "ID")
	if err := db.Model(&item).Updates(body).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

// Helper to parse id from path like /plan/123
func ParseUintFromPath(path string) (uint, error) {
	// find last slash
	idx := -1
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			idx = i
			break
		}
	}
	if idx == -1 || idx == len(path)-1 {
		return 0, strconv.ErrSyntax
	}
	v, err := strconv.Atoi(path[idx+1:])
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}
