package handlers

import (
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/controllers"
	"gorm.io/gorm"
)

func RegisterPreciosProdServRoutes(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("/precios_prodserv", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListPreciosProdServ(db, w, r)
		case http.MethodPost:
			controllers.CreatePreciosProdServ(db, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/precios_prodserv/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.ListPreciosProdServByPlan(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/precios_prodserv/item/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetPreciosProdServ(db, w, r, id)
		case http.MethodPatch:
			controllers.UpdatePreciosProdServPatch(db, w, r, id)
		case http.MethodDelete:
			controllers.DeletePreciosProdServ(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
