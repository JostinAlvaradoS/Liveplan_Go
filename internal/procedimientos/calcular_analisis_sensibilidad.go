package procedimientos

import (
	"fmt"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

// CalcularAnalisisSensibilidad calcula el análisis de sensibilidad para un plan:
// 1. Obtiene todos los registros de AnalisisSensibilidad para el plan
// 2. Para cada combinación volumen/costo:
//   - Guarda los valores originales de VariablesDeSensibilidad
//   - Actualiza VariablesDeSensibilidad con los valores de la combinación
//   - Ejecuta Recalcular para actualizar todos los cálculos
//   - Obtiene el VAN de EvaluacionProyecto
//   - Actualiza el registro de AnalisisSensibilidad con el VAN
//   - Restaura los valores originales de VariablesDeSensibilidad
func CalcularAnalisisSensibilidad(db *gorm.DB, planID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. Obtener todos los registros de AnalisisSensibilidad para el plan
		var analisisList []models.AnalisisSensibilidad
		if err := tx.Where("plan_negocio_id = ?", planID).Find(&analisisList).Error; err != nil {
			return fmt.Errorf("error al obtener AnalisisSensibilidad: %w", err)
		}

		// 2. Obtener las variables de sensibilidad originales
		var variablesOriginales models.VariablesDeSensibilidad
		if err := tx.Where("plan_negocio_id = ?", planID).First(&variablesOriginales).Error; err != nil {
			return fmt.Errorf("error al obtener VariablesDeSensibilidad originales: %w", err)
		}

		// 3. Para cada combinación volumen/costo, calcular el VAN
		for i, analisis := range analisisList {
			// Actualizar VariablesDeSensibilidad con los valores de la combinación
			variablesTemp := variablesOriginales

			// Aplicar las variaciones de volumen y costo
			variablesTemp.Cantidad_volumen = variablesOriginales.Cantidad_volumen * (1 + analisis.Volumen/100.0)
			variablesTemp.Costo = variablesOriginales.Costo * (1 + analisis.Costo/100.0)

			// Guardar las variables temporales (sin usar transacción anidada)
			if err := tx.Model(&models.VariablesDeSensibilidad{}).
				Where("plan_negocio_id = ?", planID).
				Updates(map[string]interface{}{
					"cantidad_volumen": variablesTemp.Cantidad_volumen,
					"costo":            variablesTemp.Costo,
				}).Error; err != nil {
				return fmt.Errorf("error al actualizar VariablesDeSensibilidad temporalmente: %w", err)
			}

			// Ejecutar recálculo completo
			if err := Recalcular(tx, planID); err != nil {
				return fmt.Errorf("error en recálculo para volumen %.2f%%, costo %.2f%%: %w",
					analisis.Volumen, analisis.Costo, err)
			}

			// Obtener el VAN calculado
			var evaluacion models.EvaluacionProyecto
			if err := tx.Where("plan_negocio_id = ?", planID).First(&evaluacion).Error; err != nil {
				return fmt.Errorf("error al obtener EvaluacionProyecto: %w", err)
			}

			// Actualizar el registro de AnalisisSensibilidad con el VAN
			analisisList[i].Valor = evaluacion.VAN
			if err := tx.Save(&analisisList[i]).Error; err != nil {
				return fmt.Errorf("error al actualizar AnalisisSensibilidad: %w", err)
			}
		}

		// 4. Restaurar los valores originales de VariablesDeSensibilidad
		if err := tx.Model(&models.VariablesDeSensibilidad{}).
			Where("plan_negocio_id = ?", planID).
			Updates(map[string]interface{}{
				"cantidad_volumen": variablesOriginales.Cantidad_volumen,
				"costo":            variablesOriginales.Costo,
			}).Error; err != nil {
			return fmt.Errorf("error al restaurar VariablesDeSensibilidad originales: %w", err)
		}

		return nil
	})
}
