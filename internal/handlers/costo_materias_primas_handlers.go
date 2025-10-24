package handlers

import (
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/controllers"
	"gorm.io/gorm"
)

func RegisterCostoMateriasPrimasRoutes(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("/costo_materias_primas", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListCostoMateriasPrimas(db, w, r)
		case http.MethodPost:
			controllers.CreateCostoMateriasPrimas(db, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/costo_materias_primas/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.ListCostoMateriasPrimasByPlan(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/costo_materias_primas/item/", func(w http.ResponseWriter, r *http.Request) {
		id, err := controllers.ParseUintFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			controllers.GetCostoMateriasPrimas(db, w, r, id)
		case http.MethodPatch:
			controllers.UpdateCostoMateriasPrimasPatch(db, w, r, id)
		case http.MethodDelete:
			controllers.DeleteCostoMateriasPrimas(db, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
