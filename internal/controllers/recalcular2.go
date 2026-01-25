package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/procedimientos"
	"gorm.io/gorm"
)

// ExecuteRecalcular2 ejecuta el análisis de sensibilidad completo
func ExecuteRecalcular2(db *gorm.DB, w http.ResponseWriter, r *http.Request, planID uint) {
	w.Header().Set("Content-Type", "application/json")

	// Ejecutar Recalcular2
	if err := procedimientos.Recalcular2(db, planID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("error ejecutando recalcular 2: %v", err)})
		return
	}

	// Retornar respuesta exitosa
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Recalcular 2 ejecutado exitosamente"})
}
