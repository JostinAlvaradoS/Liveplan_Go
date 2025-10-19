package handlers

import (
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/controllers"
	"gorm.io/gorm"
)

func RegisterInversionesRoutes(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("/inversiones", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListInversiones(db, w, r)
		case http.MethodPost:
			controllers.CreateInversion(db, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/inversiones/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.ListInversionesByPlan(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/inversiones/item/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetInversion(db, w, r, id)
		case http.MethodPatch:
			controllers.UpdateInversionPatch(db, w, r, id)
		case http.MethodDelete:
			controllers.DeleteInversion(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
