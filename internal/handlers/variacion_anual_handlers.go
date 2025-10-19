package handlers

import (
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/controllers"
	"gorm.io/gorm"
)

func RegisterVariacionAnualRoutes(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("/variacion_anual", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListVariacionAnual(db, w, r)
		case http.MethodPost:
			controllers.CreateVariacionAnual(db, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/variacion_anual/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.ListVariacionAnualByPlan(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/variacion_anual/item/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetVariacionAnual(db, w, r, id)
		case http.MethodPatch:
			controllers.UpdateVariacionAnualPatch(db, w, r, id)
		case http.MethodDelete:
			controllers.DeleteVariacionAnual(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
