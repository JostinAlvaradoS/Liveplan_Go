package procedimientos

import (
	"fmt"

	"gorm.io/gorm"
)

// (lockPlan removed) concurrency control for recalculations is handled elsewhere.

// Recalcular ejecuta los procedimientos de cálculo dependientes para un plan:
//   - CalcularDepreciaciones
//   - CalcularPresupuestos
//
// Se ejecutan en paralelo y si alguno falla, se devuelve un error que concatena
// los errores ocurridos en ambos procedimientos.
func Recalcular(db *gorm.DB, planID uint) error {
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
		func() error { return CalcularCostoMateriasPrimas(db, planID) },
		func() error { return CalcularVentas(db, planID) },
		func() error { return CalcularCostosVentas(db, planID) },
	}
	if err := runAdaptive(stage2Tasks); err != nil {
		return fmt.Errorf("recalcular (stage2): %w", err)
	}

	// Stage 3: calcular depreciaciones y presupuestos en paralelo
	// Stage 3: depreciaciones + presupuestos (also adaptive)
	stage3Tasks := []func() error{
		func() error { return CalcularEstadoResultados(db, planID) },
	}
	if err := runAdaptive(stage3Tasks); err != nil {
		return fmt.Errorf("recalcular (stage3): %w", err)
	}

	// Stage 3: calcular depreciaciones y presupuestos en paralelo
	// Stage 3: depreciaciones + presupuestos (also adaptive)
	stage4Tasks := []func() error{
		func() error { return CalcularFlujoEfectivo(db, planID) },
	}
	if err := runAdaptive(stage4Tasks); err != nil {
		return fmt.Errorf("recalcular (stage4): %w", err)
	}

	stage5Tasks := []func() error{
		func() error { return CalcularBalanceGeneral(db, planID) },
	}
	if err := runAdaptive(stage5Tasks); err != nil {
		return fmt.Errorf("recalcular (stage5): %w", err)
	}


	// Stage 3: depreciaciones + presupuestos (also adaptive)
	stage6Tasks := []func() error{
		func() error { return CalcularEvaluacion(db, planID) },
	}
	if err := runAdaptive(stage6Tasks); err != nil {
		return fmt.Errorf("recalcular (stage6): %w", err)
	}

	return nil
}
