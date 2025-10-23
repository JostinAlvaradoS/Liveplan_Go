package procedimientos

import (
	"fmt"

	"gorm.io/gorm"
)

// Recalcular ejecuta los procedimientos de cálculo dependientes para un plan:
//   - CalcularDepreciaciones
//   - CalcularPresupuestos
//
// Se ejecutan en paralelo y si alguno falla, se devuelve un error que concatena
// los errores ocurridos en ambos procedimientos.
func Recalcular(db *gorm.DB, planID uint) error {
	print("INICIO DE RECALCULAR")
	// Stage 1: try to run precios+costos and composicion in parallel, fallback to sequential on error
	stage1Tasks := []func() error{
		func() error { return CalcularPreciosYCostosPorPlan(db, planID) },
		func() error { return CalcularComposicion(db, planID) },
	}
	if err := runAdaptive(stage1Tasks); err != nil {
		return fmt.Errorf("recalcular (stage1): %w", err)
	}

	// Stage 2: calcular depreciaciones y presupuestos en paralelo
	// Stage 2: depreciaciones + presupuestos (also adaptive)
	stage2Tasks := []func() error{
		func() error { return CalcularDepreciaciones(db, planID) },
		func() error { return CalcularPresupuestos(db, planID) },
		func() error { return CalcularPrestamo(db, planID) },
	}
	if err := runAdaptive(stage2Tasks); err != nil {
		return fmt.Errorf("recalcular (stage2): %w", err)
	}

	print("FINAL DE RECALCULAR")
	return nil
}
