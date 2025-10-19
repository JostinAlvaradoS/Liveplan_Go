package handlers

import (
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/controllers"
	"gorm.io/gorm"
)

func RegisterTiposInversionRoutes(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("/tipos_inversion", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListTipos(db, w, r)
		case http.MethodPost:
			controllers.CreateTipo(db, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}

func RegisterPlanNegocioRoutes(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("/plan", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListPlanNegocios(db, w, r)
		case http.MethodPost:
			controllers.CreatePlanNegocio(db, w, r)
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
			controllers.GetPlanNegocio(db, w, r, id)
		case http.MethodPatch:
			controllers.UpdatePlanNegocioPatch(db, w, r, id)
		case http.MethodDelete:
			controllers.DeletePlanNegocio(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
