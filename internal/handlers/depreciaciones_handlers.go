package handlers

import (
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/controllers"
	"gorm.io/gorm"
)

func RegisterDepreciacionesRoutes(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("/depreciaciones", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListDepreciaciones(db, w, r)
		case http.MethodPost:
			controllers.CreateDepreciacion(db, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/depreciaciones/plan/", func(w http.ResponseWriter, r *http.Request) {
		// expects /depreciaciones/plan/{plan_id}
		pid, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		controllers.ListDepreciacionesByPlan(db, w, r, pid)
	})

	mux.HandleFunc("/depreciaciones/item/", func(w http.ResponseWriter, r *http.Request) {
		// expects /depreciaciones/item/{id}
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetDepreciacion(db, w, r, id)
		case http.MethodPatch:
			controllers.UpdateDepreciacionPatch(db, w, r, id)
		case http.MethodDelete:
			controllers.DeleteDepreciacion(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
