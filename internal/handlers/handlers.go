package handlers

import (
	"net/http"

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

	// Delegate registration to per-table handler registrars
	RegisterPlanNegocioRoutes(mux, a.DB)
	RegisterTiposInversionRoutes(mux, a.DB)
	RegisterProductoServicioRoutes(mux, a.DB)
	RegisterSupuestoRoutes(mux, a.DB)
	RegisterVentaDiariaRoutes(mux, a.DB)
	RegisterVariablesDeSensibilidadRoutes(mux, a.DB)
	RegisterVariacionAnualRoutes(mux, a.DB)
	RegisterPreciosProdServRoutes(mux, a.DB)
	RegisterCategoriaCostoRoutes(mux, a.DB)
	RegisterCostosProdServRoutes(mux, a.DB)
	RegisterInversionesRoutes(mux, a.DB)
	RegisterDetallesInversionRoutes(mux, a.DB)

	return mux
}
