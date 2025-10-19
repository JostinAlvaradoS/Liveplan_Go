package handlers

import (
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/controllers"
	"gorm.io/gorm"
)

func RegisterVariablesDeSensibilidadRoutes(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("/variables_sensibilidad", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListVariablesDeSensibilidad(db, w, r)
		case http.MethodPost:
			controllers.CreateVariablesDeSensibilidad(db, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/variables_sensibilidad/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.ListVariablesDeSensibilidadByPlan(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/variables_sensibilidad/item/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetVariablesDeSensibilidad(db, w, r, id)
		case http.MethodPatch:
			controllers.UpdateVariablesDeSensibilidadPatch(db, w, r, id)
		case http.MethodDelete:
			controllers.DeleteVariablesDeSensibilidad(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
