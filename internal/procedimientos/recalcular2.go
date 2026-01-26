package procedimientos

import (
	"fmt"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

// Recalcular2 ejecuta el análisis de sensibilidad completo:
// 1. Obtiene y guarda en memoria las variables de sensibilidad actuales
// 2. Recorre todos los registros de AnalisisSensibilidad con todas las combinaciones
// 3. Para cada combinación, actualiza las variables, ejecuta Recalcular y guarda el VAN
// 4. Restaura las variables de sensibilidad originales
// 5. Llama a Recalcular original para dejar todo en su estado final correcto
// 6. Establece el status en true indicando que el análisis está completo
func Recalcular2(db *gorm.DB, planID uint) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. Obtener y guardar en memoria las variables de sensibilidad originales
		var variablesOriginales models.VariablesDeSensibilidad
		if err := tx.Where("plan_negocio_id = ?", planID).First(&variablesOriginales).Error; err != nil {
			return fmt.Errorf("error al obtener VariablesDeSensibilidad originales: %w", err)
		}

		// Guardar una copia en memoria
		variablesBackup := variablesOriginales

		// 2. Obtener todos los registros de AnalisisSensibilidad para el plan
		var analisisList []models.AnalisisSensibilidad
		if err := tx.Where("plan_negocio_id = ?", planID).Find(&analisisList).Error; err != nil {
			return fmt.Errorf("error al obtener AnalisisSensibilidad: %w", err)
		}

		// 3. Recorrer todas las combinaciones de variables de sensibilidad
		for _, analisis := range analisisList {
			// Crear copia temporal de las variables
			variablesTemp := variablesOriginales

			// Aplicar las variaciones de volumen y costo
			variablesTemp.Cantidad_volumen = variablesOriginales.Cantidad_volumen * (1 + analisis.Volumen/100.0)
			variablesTemp.Costo = variablesOriginales.Costo * (1 + analisis.Costo/100.0)

			// Actualizar VariablesDeSensibilidad en la BD
			if err := tx.Model(&models.VariablesDeSensibilidad{}).
				Where("plan_negocio_id = ?", planID).
				Updates(map[string]interface{}{
					"cantidad_volumen": variablesTemp.Cantidad_volumen,
					"costo":            variablesTemp.Costo,
				}).Error; err != nil {
				return fmt.Errorf("error al actualizar VariablesDeSensibilidad: %w", err)
			}

			// Ejecutar Recalcular para actualizar todos los cálculos
			// Nota: Recalcular debe ser llamado sin transacción anidada
			if err := Recalcular(db, planID); err != nil {
				return fmt.Errorf("error ejecutando Recalcular para análisis sensibilidad: %w", err)
			}

			// Obtener el VAN de EvaluacionProyecto
			var evaluacion models.EvaluacionProyecto
			if err := tx.Where("plan_negocio_id = ?", planID).First(&evaluacion).Error; err != nil {
				return fmt.Errorf("error al obtener EvaluacionProyecto: %w", err)
			}

			// Actualizar el registro de AnalisisSensibilidad con el Valor (VAN)
			if err := tx.Model(&analisis).Update("valor", evaluacion.VAN).Error; err != nil {
				return fmt.Errorf("error al actualizar AnalisisSensibilidad con VAN: %w", err)
			}
		}

		// 4. Restaurar las variables de sensibilidad originales
		if err := tx.Model(&models.VariablesDeSensibilidad{}).
			Where("plan_negocio_id = ?", planID).
			Updates(map[string]interface{}{
				"cantidad_volumen": variablesBackup.Cantidad_volumen,
				"costo":            variablesBackup.Costo,
				"precio":           variablesBackup.Precio,
			}).Error; err != nil {
			return fmt.Errorf("error al restaurar VariablesDeSensibilidad originales: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 5. Ejecutar Recalcular final para dejar todo en su estado correcto
	if err := Recalcular(db, planID); err != nil {
		return fmt.Errorf("error ejecutando Recalcular final en recalcular2: %w", err)
	}

	// 6. Set status to true when sensitivity analysis is complete
	if err := SetStatusSensibilidad(db, planID, true); err != nil {
		return fmt.Errorf("recalcular2 (set status): %w", err)
	}

	return nil
}
