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
		// path is /plan/{idOrUid}
		seg := r.URL.Path[len("/plan/"):]
		if seg == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			// allow string UID or numeric id
			if id, err := controllers.ParseUintFromPath(r.URL.Path); err == nil {
				controllers.ListPlanesByUser(db, w, r, id)
				return
			}
			// fallback: treat as user UID string
			controllers.ListPlanesByUserUID(db, w, r, seg)
		case http.MethodPatch:
			// Patch requires numeric id
			id, err := controllers.ParseUintFromPath(r.URL.Path)
			if err != nil {
				http.Error(w, "invalid id", http.StatusBadRequest)
				return
			}
			controllers.UpdatePlanNegocioPatch(db, w, r, id)
		case http.MethodDelete:
			id, err := controllers.ParseUintFromPath(r.URL.Path)
			if err != nil {
				http.Error(w, "invalid id", http.StatusBadRequest)
				return
			}
			controllers.DeletePlanNegocio(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
