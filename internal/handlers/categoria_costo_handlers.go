package handlers

import (
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/controllers"
	"gorm.io/gorm"
)

func RegisterCategoriaCostoRoutes(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("/categoria_costo", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListCategoriaCosto(db, w, r)
		case http.MethodPost:
			controllers.CreateCategoriaCosto(db, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/categoria_costo/item/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetCategoriaCosto(db, w, r, id)
		case http.MethodPatch:
			controllers.UpdateCategoriaCostoPatch(db, w, r, id)
		case http.MethodDelete:
			controllers.DeleteCategoriaCosto(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
