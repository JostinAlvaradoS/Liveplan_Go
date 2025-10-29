package handlers

import (
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/controllers"
	"gorm.io/gorm"
)

func RegisterConceptosEvaluacionRoutes(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("/conceptos_evaluacion", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListConceptosEvaluacion(db, w, r)
		case http.MethodPost:
			controllers.CreateConceptoEvaluacion(db, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/conceptos_evaluacion/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.ListConceptosEvaluacionByPlan(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/conceptos_evaluacion/item/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetConceptoEvaluacion(db, w, r, id)
		case http.MethodPut:
			controllers.UpdateConceptoEvaluacion(db, w, r, id)
		case http.MethodDelete:
			controllers.DeleteConceptoEvaluacion(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
