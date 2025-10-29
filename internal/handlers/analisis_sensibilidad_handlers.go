package handlers

import (
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/controllers"
	"gorm.io/gorm"
)

func RegisterAnalisisSensibilidadRoutes(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("/analisis_sensibilidad", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListAnalisisSensibilidad(db, w, r)
		case http.MethodPost:
			controllers.CreateAnalisisSensibilidad(db, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/analisis_sensibilidad/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.ListAnalisisSensibilidadByPlan(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/analisis_sensibilidad/item/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetAnalisisSensibilidad(db, w, r, id)
		case http.MethodPut:
			controllers.UpdateAnalisisSensibilidad(db, w, r, id)
		case http.MethodDelete:
			controllers.DeleteAnalisisSensibilidad(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
