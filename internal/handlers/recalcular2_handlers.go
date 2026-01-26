package handlers

import (
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/controllers"
	"gorm.io/gorm"
)

func RegisterRecalcular2Routes(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("/recalcular2/", func(w http.ResponseWriter, r *http.Request) {
		planID, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid plan id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodPost:
			controllers.ExecuteRecalcular2(db, w, r, planID)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
