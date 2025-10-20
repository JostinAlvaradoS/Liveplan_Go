package procedimientos

import (
	"fmt"
	"strings"
	"sync"

	"gorm.io/gorm"
)

// Recalcular ejecuta los procedimientos de cálculo dependientes para un plan:
//   - CalcularDepreciaciones
//   - CalcularPresupuestos
//
// Se ejecutan en paralelo y si alguno falla, se devuelve un error que concatena
// los errores ocurridos en ambos procedimientos.
func Recalcular(db *gorm.DB, planID uint) error {
	var wg sync.WaitGroup
	errs := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		errs <- CalcularDepreciaciones(db, planID)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		errs <- CalcularPresupuestos(db, planID)
	}()

	wg.Wait()
	close(errs)

	var parts []string
	for e := range errs {
		if e != nil {
			parts = append(parts, e.Error())
		}
	}
	if len(parts) > 0 {
		return fmt.Errorf("recalcular: %s", strings.Join(parts, "; "))
	}
	return nil
}
