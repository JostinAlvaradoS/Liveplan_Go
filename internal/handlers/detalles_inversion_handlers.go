package handlers

import (
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/controllers"
	"gorm.io/gorm"
)

func RegisterDetallesInversionRoutes(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("/detalles_inversion", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListDetalles(db, w, r)
		case http.MethodPost:
			controllers.CreateDetalle(db, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/detalles_inversion/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.ListDetallesByPlan(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/detalles_inversion/item/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetDetalle(db, w, r, id)
		case http.MethodPatch:
			controllers.UpdateDetallePatch(db, w, r, id)
		case http.MethodDelete:
			controllers.DeleteDetalle(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
