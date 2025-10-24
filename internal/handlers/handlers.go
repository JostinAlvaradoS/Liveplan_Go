package handlers

import (
	"net/http"

	"gorm.io/gorm"
)

type App struct {
	DB *gorm.DB
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := "http://localhost:4200,https://liveplan-frontend.web.app"
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-Custom-Header")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
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
	RegisterCostoMateriasPrimasRoutes(mux, a.DB)
	RegisterIndicadoresMacroRoutes(mux, a.DB)
	RegisterComposicionFinanciamientoRoutes(mux, a.DB)
	RegisterDepreciacionesRoutes(mux, a.DB)
	RegisterPresupuestoVentaRoutes(mux, a.DB)
	RegisterInversionesRoutes(mux, a.DB)
	RegisterDetallesInversionRoutes(mux, a.DB)
	RegisterVentasDineroRoutes(mux, a.DB)
	RegisterPrestamoRoutes(mux, a.DB)

	return corsMiddleware(mux)
}
