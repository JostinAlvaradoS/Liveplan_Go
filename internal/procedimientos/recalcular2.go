package procedimientos

import (
	"fmt"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

// Recalcular2 ejecuta el análisis de sensibilidad completo:
// 1. Obtiene y guarda las variables de sensibilidad originales
// 2. Recorre todos los registros de AnalisisSensibilidad
// 3. Para cada combinación, actualiza variables, ejecuta cálculos y guarda VAN
// 4. Restaura las variables originales y ejecuta Recalcular final
// 5. Establece status en true
func Recalcular2(db *gorm.DB, planID uint) error {
	// 1. Obtener variables originales
	var variablesOriginales models.VariablesDeSensibilidad
	if err := db.Where("plan_negocio_id = ?", planID).First(&variablesOriginales).Error; err != nil {
		return fmt.Errorf("error al obtener VariablesDeSensibilidad originales: %w", err)
	}

	variablesBackup := variablesOriginales

	// 2. Obtener todas las combinaciones de análisis de sensibilidad
	var analisisList []models.AnalisisSensibilidad
	if err := db.Where("plan_negocio_id = ?", planID).Find(&analisisList).Error; err != nil {
		return fmt.Errorf("error al obtener AnalisisSensibilidad: %w", err)
	}

	// 3. Recorrer cada combinación
	for idx, analisis := range analisisList {
		fmt.Printf("[%d] Vol: %+.2f%% | Costo: %+.2f%%\n", idx+1, analisis.Volumen, analisis.Costo)

		// Actualizar variables de sensibilidad directamente con los porcentajes
		if err := db.Model(&models.VariablesDeSensibilidad{}).
			Where("plan_negocio_id = ?", planID).
			Updates(map[string]interface{}{
				"cantidad_volumen": analisis.Volumen,
				"costo":            analisis.Costo,
			}).Error; err != nil {
			return fmt.Errorf("error actualizando variables [%d]: %w", idx+1, err)
		}

		// Ejecutar cálculos completos para esta combinación
		if err := Recalcular(db, planID); err != nil {
			return fmt.Errorf("error en Recalcular [%d]: %w", idx+1, err)
		}

		// Obtener VAN actualizado
		var evaluacion models.EvaluacionProyecto
		if err := db.Where("plan_negocio_id = ?", planID).First(&evaluacion).Error; err != nil {
			return fmt.Errorf("error obteniendo VAN [%d]: %w", idx+1, err)
		}

		fmt.Printf("[%d] VAN: %.2f\n", idx+1, evaluacion.VAN)

		// Guardar VAN en tabla de análisis
		if err := db.Model(&analisis).Update("valor", evaluacion.VAN).Error; err != nil {
			return fmt.Errorf("error guardando VAN [%d]: %w", idx+1, err)
		}
	}

	// 4. Restaurar variables originales
	if err := db.Model(&models.VariablesDeSensibilidad{}).
		Where("plan_negocio_id = ?", planID).
		Updates(map[string]interface{}{
			"cantidad_volumen": variablesBackup.Cantidad_volumen,
			"costo":            variablesBackup.Costo,
			"precio":           variablesBackup.Precio,
		}).Error; err != nil {
		return fmt.Errorf("error al restaurar VariablesDeSensibilidad originales: %w", err)
	}

	// Ejecutar Recalcular final con variables originales
	if err := Recalcular(db, planID); err != nil {
		return fmt.Errorf("error en Recalcular final: %w", err)
	}

	// 5. Establecer status en true
	if err := SetStatusSensibilidad(db, planID, true); err != nil {
		return fmt.Errorf("recalcular2 (set status): %w", err)
	}

	return nil
}
