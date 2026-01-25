package handlers

import (
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/controllers"
	"gorm.io/gorm"
)

func RegisterAnalisisSensibilidadStatusRoutes(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("/analisis_sensibilidad_status/", func(w http.ResponseWriter, r *http.Request) {
		planID, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid plan id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetAnalisisSensibilidadStatus(db, w, r, planID)
		case http.MethodPut:
			controllers.UpdateAnalisisSensibilidadStatus(db, w, r, planID)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
