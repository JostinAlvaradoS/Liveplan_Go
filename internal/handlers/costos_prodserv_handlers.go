package handlers

import (
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/controllers"
	"gorm.io/gorm"
)

func RegisterCostosProdServRoutes(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("/costos_prodserv", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListCostosProdServ(db, w, r)
		case http.MethodPost:
			controllers.CreateCostosProdServ(db, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/costos_prodserv/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.ListCostosProdServByPlan(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/costos_prodserv/item/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetCostosProdServ(db, w, r, id)
		case http.MethodPatch:
			controllers.UpdateCostosProdServPatch(db, w, r, id)
		case http.MethodDelete:
			controllers.DeleteCostosProdServ(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/costos_prodserv/report_by_plan/", func(w http.ResponseWriter, r *http.Request) {
		// path expects plan id after the trailing slash
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid plan id", http.StatusBadRequest)
			return
		}
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		controllers.ReportCostosPorProducto(db, w, r, id)
	})
}
