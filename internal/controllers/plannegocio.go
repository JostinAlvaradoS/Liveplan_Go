package controllers

import (
	"encoding/json"
	"net/http"

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
	// Use transaction to ensure related default records are created atomically
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&item).Error; err != nil {
			return err
		}

		// Create default VariacionAnual (all zeros)
		va := models.VariacionAnual{
			PlanNegocioID: item.ID,
			Año1:          0,
			Año2:          0,
			Año3:          0,
			Año4:          0,
			Año5:          0,
		}
		if err := tx.Create(&va).Error; err != nil {
			return err
		}

		// Create default VariablesDeSensibilidad (zeros)
		vs := models.VariablesDeSensibilidad{
			Cantidad_volumen: 0,
			Precio:           0,
			Costo:            0,
			PlanNegocioID:    item.ID,
		}
		if err := tx.Create(&vs).Error; err != nil {
			return err
		}

		// Create default Supuesto
		sup := models.Supuesto{
			PlanNegocioID:         item.ID,
			PorcenVentas:          0,
			VariacionPorcenVentas: 0,
			PTU:                   0,
			ISR:                   0,
		}
		if err := tx.Create(&sup).Error; err != nil {
			return err
		}

		// Create default IndicadoresMacro
		im := models.IndicadoresMacro{
			PlanNegocioID: item.ID,
			TipoCambio:    0,
			Inflacion:     0,
			TasaDeuda:     0,
			TasaInteres:   0,
			TasaImpuesto:  0,
			PTU:           0,
			DiasxMes:      30,
		}
		if err := tx.Create(&im).Error; err != nil {
			return err
		}

		// Create default ComposicionFinanciamiento
		cf := models.ComposicionFinanciamiento{
			PlanNegocioID:     item.ID,
			CapitalPorcentaje: 50,
			DeudaPorcentaje:   50,
			Total_Inversion:   0,
		}
		if err := tx.Create(&cf).Error; err != nil {
			return err
		}

		// Create default DatosPrestamo and PrestamoCuotas for 5 years (60 meses)
		dp := models.DatosPrestamo{
			PlanNegocioID:          item.ID,
			Monto:                  0,
			TasaAnual:              12,
			PeriodosCapitalizacion: 12,
			TasaMensual:            1,
			PeriodosAmortizacion:   60,
			Cuota:                  0,
		}
		if err := tx.Create(&dp).Error; err != nil {
			return err
		}

		// Create 60 PrestamoCuotas rows, one per month (periodo_mes 1..60)
		for m := 1; m <= 60; m++ {
			anio := (m-1)/12 + 1 // 1..5
			mes := (m-1)%12 + 1  // 1..12
			pc := models.PrestamoCuotas{
				PlanNegocioID:  item.ID,
				PeriodoMes:     m,
				Anio:           anio,
				Mes:            mes,
				SaldoInicial:   0,
				Interes:        0,
				Amortizacion:   0,
				CuotaTotal:     0,
				SaldoPendiente: 0,
			}
			if err := tx.Create(&pc).Error; err != nil {
				return err
			}
		}

		// Create default EstadoResultados rows: 5 years x 12 months (anio 1..5, mes 1..12)
		for anio := 1; anio <= 5; anio++ {
			for mes := 1; mes <= 12; mes++ {
				er := models.EstadoResultados{
					PlanNegocioID:          item.ID,
					Anio:                   anio,
					Mes:                    mes,
					Ventas:                 0,
					CostosVentas:           0,
					UtilidadBruta:          0,
					GastosVentaAdm:         0,
					Depreciacion:           0,
					Amortizacion:           0,
					UtilidadprevioIntImp:   0,
					GastosFinancieros:      0,
					UtilidadAntesPTU:       0,
					PTU:                    0,
					UtilidadAntesImpuestos: 0,
					ISR:                    0,
					UtilidadNeta:           0,
				}
				if err := tx.Create(&er).Error; err != nil {
					return err
				}
			}
		}

		// Poblar GastosOperacion desde GastosOperacionBase
		var gastosBase []models.GastosOperacionBase
		if err := tx.Find(&gastosBase).Error; err != nil {
			return err
		}
		for _, gb := range gastosBase {
			gope := models.GastosOperacion{
				PlanNegocioID: item.ID,
				Descripcion:   gb.Descripcion,
				Mensual:       gb.Valor,
				Anual:         gb.Valor * 12,
			}
			if err := tx.Create(&gope).Error; err != nil {
				return err
			}
		}

		// Crear PoliticasVenta y PoliticasCompra por defecto (80-20) para todos los meses de los 5 años
		for anio := 1; anio <= 5; anio++ {
			for mes := 1; mes <= 12; mes++ {
				pv := models.PoliticasVenta{
					PlanNegocioID:     item.ID,
					Anio:              anio,
					Mes:               mes,
					PorcentajeCredito: 20,
					PorcentajeContado: 80,
				}
				if err := tx.Create(&pv).Error; err != nil {
					return err
				}
				pc := models.PoliticasCompra{
					PlanNegocioID:     item.ID,
					Anio:              anio,
					Mes:               mes,
					PorcentajeCredito: 20,
					PorcentajeContado: 80,
				}
				if err := tx.Create(&pc).Error; err != nil {
					return err
				}
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
	delete(body, "recalc")

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

func ListPlanesByUser(db *gorm.DB, w http.ResponseWriter, r *http.Request, userID uint) {
	var items []models.PlanNegocio
	if err := db.Where("autor = ?", userID).Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// ListPlanesByUserUID lista planes cuyo campo Autor coincide con el UID (string).
func ListPlanesByUserUID(db *gorm.DB, w http.ResponseWriter, r *http.Request, userUID string) {
	var items []models.PlanNegocio
	if err := db.Where("autor = ?", userUID).Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
