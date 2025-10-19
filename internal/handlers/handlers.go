package handlers

import (
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/controllers"
	"gorm.io/gorm"
)

type App struct {
	DB *gorm.DB
}

func (a *App) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// PlanNegocio routes
	mux.HandleFunc("/plan", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListPlanNegocios(a.DB, w, r)
		case http.MethodPost:
			controllers.CreatePlanNegocio(a.DB, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/plan/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetPlanNegocio(a.DB, w, r, id)
		case http.MethodPatch:
			controllers.UpdatePlanNegocioPatch(a.DB, w, r, id)
		case http.MethodDelete:
			controllers.DeletePlanNegocio(a.DB, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// TipoInversionInicial routes
	mux.HandleFunc("/tipos_inversion", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListTipos(a.DB, w, r)
		case http.MethodPost:
			controllers.CreateTipo(a.DB, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// ProductoServicio routes
	mux.HandleFunc("/producto_servicio", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListProductoServicios(a.DB, w, r)
		case http.MethodPost:
			controllers.CreateProductoServicio(a.DB, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	// Treat GET /producto_servicio/{plan_id} as list of producto_servicio for that plan
	mux.HandleFunc("/producto_servicio/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.ListProductoServiciosByPlan(a.DB, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	// Item-specific routes for producto (use /producto_servicio/item/{id})
	mux.HandleFunc("/producto_servicio/item/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetProductoServicio(a.DB, w, r, id)
		case http.MethodPatch:
			controllers.UpdateProductoServicioPatch(a.DB, w, r, id)
		case http.MethodDelete:
			controllers.DeleteProductoServicio(a.DB, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Supuesto routes
	mux.HandleFunc("/supuestos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListSupuestos(a.DB, w, r)
		case http.MethodPost:
			controllers.CreateSupuesto(a.DB, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	// GET /supuestos/{plan_id} -> list by plan
	mux.HandleFunc("/supuestos/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.ListSupuestosByPlan(a.DB, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	// Item-specific supuestos
	mux.HandleFunc("/supuestos/item/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetSupuesto(a.DB, w, r, id)
		case http.MethodPatch:
			controllers.UpdateSupuestoPatch(a.DB, w, r, id)
		case http.MethodDelete:
			controllers.DeleteSupuesto(a.DB, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// VentaDiaria routes
	mux.HandleFunc("/ventas_diarias", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListVentaDiarias(a.DB, w, r)
		case http.MethodPost:
			controllers.CreateVentaDiaria(a.DB, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	// GET /ventas_diarias/{plan_id} -> list by plan
	mux.HandleFunc("/ventas_diarias/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.ListVentaDiariasByPlan(a.DB, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	// Item-specific ventas diarias
	mux.HandleFunc("/ventas_diarias/item/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetVentaDiaria(a.DB, w, r, id)
		case http.MethodPatch:
			controllers.UpdateVentaDiariaPatch(a.DB, w, r, id)
		case http.MethodDelete:
			controllers.DeleteVentaDiaria(a.DB, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// VariablesDeSensibilidad routes
	mux.HandleFunc("/variables_sensibilidad", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListVariablesDeSensibilidad(a.DB, w, r)
		case http.MethodPost:
			controllers.CreateVariablesDeSensibilidad(a.DB, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	// GET /variables_sensibilidad/{plan_id} -> list by plan
	mux.HandleFunc("/variables_sensibilidad/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.ListVariablesDeSensibilidadByPlan(a.DB, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	// Item-specific variables_sensibilidad
	mux.HandleFunc("/variables_sensibilidad/item/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetVariablesDeSensibilidad(a.DB, w, r, id)
		case http.MethodPatch:
			controllers.UpdateVariablesDeSensibilidadPatch(a.DB, w, r, id)
		case http.MethodDelete:
			controllers.DeleteVariablesDeSensibilidad(a.DB, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// VariacionAnual routes
	mux.HandleFunc("/variacion_anual", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListVariacionAnual(a.DB, w, r)
		case http.MethodPost:
			controllers.CreateVariacionAnual(a.DB, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	// GET /variacion_anual/{plan_id}
	mux.HandleFunc("/variacion_anual/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.ListVariacionAnualByPlan(a.DB, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	// Item-specific variacion_anual
	mux.HandleFunc("/variacion_anual/item/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetVariacionAnual(a.DB, w, r, id)
		case http.MethodPatch:
			controllers.UpdateVariacionAnualPatch(a.DB, w, r, id)
		case http.MethodDelete:
			controllers.DeleteVariacionAnual(a.DB, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// PreciosProdServ routes
	mux.HandleFunc("/precios_prodserv", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListPreciosProdServ(a.DB, w, r)
		case http.MethodPost:
			controllers.CreatePreciosProdServ(a.DB, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	// GET /precios_prodserv/{plan_id}
	mux.HandleFunc("/precios_prodserv/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.ListPreciosProdServByPlan(a.DB, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	// Item-specific precios_prodserv
	mux.HandleFunc("/precios_prodserv/item/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetPreciosProdServ(a.DB, w, r, id)
		case http.MethodPatch:
			controllers.UpdatePreciosProdServPatch(a.DB, w, r, id)
		case http.MethodDelete:
			controllers.DeletePreciosProdServ(a.DB, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// CategoriaCosto routes (catalog)
	mux.HandleFunc("/categoria_costo", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListCategoriaCosto(a.DB, w, r)
		case http.MethodPost:
			controllers.CreateCategoriaCosto(a.DB, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/categoria_costo/item/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetCategoriaCosto(a.DB, w, r, id)
		case http.MethodPatch:
			controllers.UpdateCategoriaCostoPatch(a.DB, w, r, id)
		case http.MethodDelete:
			controllers.DeleteCategoriaCosto(a.DB, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// CostosProdServ routes
	mux.HandleFunc("/costos_prodserv", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListCostosProdServ(a.DB, w, r)
		case http.MethodPost:
			controllers.CreateCostosProdServ(a.DB, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	// GET /costos_prodserv/{plan_id}
	mux.HandleFunc("/costos_prodserv/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.ListCostosProdServByPlan(a.DB, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	// Item-specific costos_prodserv
	mux.HandleFunc("/costos_prodserv/item/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetCostosProdServ(a.DB, w, r, id)
		case http.MethodPatch:
			controllers.UpdateCostosProdServPatch(a.DB, w, r, id)
		case http.MethodDelete:
			controllers.DeleteCostosProdServ(a.DB, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// InversionInicial routes
	mux.HandleFunc("/inversiones", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListInversiones(a.DB, w, r)
		case http.MethodPost:
			controllers.CreateInversion(a.DB, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	// Treat GET /inversiones/{plan_id} as list of inversiones for that plan
	mux.HandleFunc("/inversiones/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.ListInversionesByPlan(a.DB, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Item-specific routes for inversiones
	mux.HandleFunc("/inversiones/item/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetInversion(a.DB, w, r, id)
		case http.MethodPatch:
			controllers.UpdateInversionPatch(a.DB, w, r, id)
		case http.MethodDelete:
			controllers.DeleteInversion(a.DB, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// DetalleInversionInicial routes
	mux.HandleFunc("/detalles_inversion", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListDetalles(a.DB, w, r)
		case http.MethodPost:
			controllers.CreateDetalle(a.DB, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	// Treat GET /detalles_inversion/{plan_id} as list of detalles for that plan
	mux.HandleFunc("/detalles_inversion/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.ListDetallesByPlan(a.DB, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Item-specific routes for detalles
	mux.HandleFunc("/detalles_inversion/item/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetDetalle(a.DB, w, r, id)
		case http.MethodPatch:
			controllers.UpdateDetallePatch(a.DB, w, r, id)
		case http.MethodDelete:
			controllers.DeleteDetalle(a.DB, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	return mux
}
