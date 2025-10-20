package handlers

import (
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/controllers"
	"gorm.io/gorm"
)

func RegisterPresupuestoVentaRoutes(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("/presupuestos_venta", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListPresupuestosVenta(db, w, r)
		case http.MethodPost:
			controllers.CreatePresupuestoVenta(db, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/presupuestos_venta/", func(w http.ResponseWriter, r *http.Request) {
		pid, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		controllers.ListPresupuestosVentaByPlan(db, w, r, pid)
	})

	mux.HandleFunc("/presupuestos_venta/item/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetPresupuestoVenta(db, w, r, id)
		case http.MethodPatch:
			controllers.UpdatePresupuestoVentaPatch(db, w, r, id)
		case http.MethodDelete:
			controllers.DeletePresupuestoVenta(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
