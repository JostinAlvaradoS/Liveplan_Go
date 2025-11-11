package handlers

import (
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/controllers"
	"gorm.io/gorm"
)

func RegisterGastosOperacionRoutes(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("/gastos_operacion", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListGastosOperacion(db, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/gastos_operacion/", func(w http.ResponseWriter, r *http.Request) {
		// this endpoint expects a plan id (e.g. /gastos_operacion/123) and returns
		// gastos operacion for that plan
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.ListGastosOperacionByPlan(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/gastos_operacion/item/", func(w http.ResponseWriter, r *http.Request) {
		// item-level operations: GET, PATCH, DELETE
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetGastosOperacion(db, w, r, id)
		case http.MethodPatch:
			controllers.UpdateGastosOperacionPatch(db, w, r, id)
		case http.MethodDelete:
			controllers.DeleteGastosOperacion(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
