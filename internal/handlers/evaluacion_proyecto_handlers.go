package handlers

import (
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/controllers"
	"gorm.io/gorm"
)

func RegisterEvaluacionProyectoRoutes(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("/evaluacion_proyecto", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListEvaluacionProyecto(db, w, r)
		case http.MethodPost:
			controllers.CreateEvaluacionProyecto(db, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/evaluacion_proyecto/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetEvaluacionProyectoByPlan(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/evaluacion_proyecto/item/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetEvaluacionProyecto(db, w, r, id)
		case http.MethodPut:
			controllers.UpdateEvaluacionProyecto(db, w, r, id)
		case http.MethodDelete:
			controllers.DeleteEvaluacionProyecto(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
