package handlers

import (
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/controllers"
	"gorm.io/gorm"
)

func RegisterComposicionFinanciamientoRoutes(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("/composicion_financiamiento", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListComposicionFinanciamiento(db, w, r)
		case http.MethodPost:
			controllers.CreateComposicionFinanciamiento(db, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/composicion_financiamiento/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.ListComposicionFinanciamientoByPlan(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/composicion_financiamiento/item/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetComposicionFinanciamiento(db, w, r, id)
		case http.MethodPatch:
			controllers.UpdateComposicionFinanciamientoPatch(db, w, r, id)
		case http.MethodDelete:
			controllers.DeleteComposicionFinanciamiento(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
