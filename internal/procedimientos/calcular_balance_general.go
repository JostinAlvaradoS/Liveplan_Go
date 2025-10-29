package procedimientos

import (
	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/gorm"
)

// CalcularBalanceGeneral calcula y actualiza los registros de BalanceGeneral para un plan
func CalcularBalanceGeneral(db *gorm.DB, planID uint) error {
	// Obtener datos de Supuesto para el cálculo de inventarios
	var supuesto models.Supuesto
	if err := db.Where("plan_negocio_id = ?", planID).First(&supuesto).Error; err != nil {
		return err
	}

	// Calcular para cada año y mes
	for anio := 1; anio <= 5; anio++ {
		// Calcular mes 0 (solo año 1)
		if anio == 1 {
			if err := calcularBalanceMes(db, planID, anio, 0, supuesto); err != nil {
				return err
			}
		}

		// Calcular meses 1-12
		for mes := 1; mes <= 12; mes++ {
			if err := calcularBalanceMes(db, planID, anio, mes, supuesto); err != nil {
				return err
			}
		}
	}

	return nil
}

func calcularBalanceMes(db *gorm.DB, planID uint, anio, mes int, supuesto models.Supuesto) error {
	// Buscar el registro de balance existente
	var balance models.BalanceGeneral
	if err := db.Where("plan_negocio_id = ? AND anio = ? AND mes = ?", planID, anio, mes).First(&balance).Error; err != nil {
		return err
	}

	// 1. Calcular Corrientes_Efectivo
	var efectivo float64
	if anio == 1 && mes == 0 {
		// Mes 0 año 1: buscar en detallesInversion donde tipo=3 y elemento="Efectivo"
		var detalleEfectivo models.DetalleInversionInicial
		err := db.Joins("JOIN tipo_inversion_inicials ON detalle_inversion_inicials.tipo_id = tipo_inversion_inicials.id").
			Where("detalle_inversion_inicials.plan_negocio_id = ? AND tipo_inversion_inicials.id = 3 AND detalle_inversion_inicials.elemento = ?", planID, "Efectivo").
			First(&detalleEfectivo).Error
		if err == nil {
			efectivo = detalleEfectivo.Importe
		}
	} else {
		// Mes 1 en adelante: tomar efectivo_final del flujo de efectivo del mes anterior
		var flujoAnterior models.FlujoEfectivo
		mesAnterior := mes - 1
		anioAnterior := anio
		if mes == 1 {
			if anio == 1 {
				// Si es mes 1 del año 1, buscar mes 0 del año 1
				mesAnterior = 0
				anioAnterior = 1
			} else {
				// Si es mes 1 de años 2-5, buscar mes 12 del año anterior
				mesAnterior = 12
				anioAnterior = anio - 1
			}
		}

		err := db.Where("plan_negocio_id = ? AND anio = ? AND mes = ?", planID, anioAnterior, mesAnterior).
			First(&flujoAnterior).Error
		if err == nil {
			efectivo = flujoAnterior.EfectivoFinal
		}
	}

	// 2. Calcular Corrientes_CuentasxCobrar
	var cuentasPorCobrar float64
	if anio == 1 && mes == 0 {
		cuentasPorCobrar = 0 // Inicializa en 0 para mes 0 año 1
	} else {
		// Obtener cuentas por cobrar del mes anterior
		var balanceAnterior models.BalanceGeneral
		mesAnterior := mes - 1
		anioAnterior := anio
		if mes == 1 {
			if anio == 1 {
				// Si es mes 1 del año 1, tomar mes 0 del año 1
				mesAnterior = 0
				anioAnterior = 1
			} else {
				// Si es mes 1 de años 2-5, tomar mes 12 del año anterior
				mesAnterior = 12
				anioAnterior = anio - 1
			}
		}

		err := db.Where("plan_negocio_id = ? AND anio = ? AND mes = ?", planID, anioAnterior, mesAnterior).
			First(&balanceAnterior).Error
		if err == nil {
			cuentasPorCobrar = balanceAnterior.Corrientes_CuentasxCobrar
		}

		// Obtener ventas del estado de resultados del mes actual
		var estadoResultados models.EstadoResultados
		err = db.Where("plan_negocio_id = ? AND anio = ? AND mes = ?", planID, anio, mes).
			First(&estadoResultados).Error
		if err == nil {
			// Obtener política de venta del mes actual
			var politicaVenta models.PoliticasVenta
			err = db.Where("plan_negocio_id = ? AND anio = ? AND mes = ?", planID, anio, mes).
				First(&politicaVenta).Error
			if err == nil {
				// Agregar ventas a crédito
				ventasCredito := estadoResultados.Ventas * (politicaVenta.PorcentajeCredito / 100.0)
				cuentasPorCobrar += ventasCredito
			}
		}

		// Restar cobros de ventas a crédito del flujo de efectivo
		var flujoEfectivo models.FlujoEfectivo
		err = db.Where("plan_negocio_id = ? AND anio = ? AND mes = ?", planID, anio, mes).
			First(&flujoEfectivo).Error
		if err == nil {
			cuentasPorCobrar -= flujoEfectivo.Ingresos_CobrosVentasCredito
		}
	}

	// 3. Calcular Corrientes_Inventarios
	var inventarios float64
	if anio == 1 && mes == 0 {
		// Mes 0: buscar en detallesInversion donde elemento="Inventario de materias primas"
		var detalleInventario models.DetalleInversionInicial
		err := db.Where("plan_negocio_id = ? AND elemento = ?", planID, "Inventario de materias primas").
			First(&detalleInventario).Error
		if err == nil {
			inventarios = detalleInventario.Importe
		}
	} else {
		// Mes 1 en adelante: ventas * supuestos.porcenventas/100
		var estadoResultados models.EstadoResultados
		err := db.Where("plan_negocio_id = ? AND anio = ? AND mes = ?", planID, anio, mes).
			First(&estadoResultados).Error
		if err == nil {
			inventarios = estadoResultados.Ventas * (supuesto.PorcenVentas / 100.0)
		}
	}

	// 4. Corrientes_Otros = 0 en todos los meses
	corrientesOtros := 0.0

	// 5. Calcular CorrientesSuma
	corrientesSuma := efectivo + cuentasPorCobrar + inventarios + corrientesOtros

	// 6. Calcular NoCorrientes_Suma
	var noCorrientesSuma float64
	if anio == 1 && mes == 0 {
		// Mes 0: suma de todos los detalles de inversión cuyo tipo es 1 o 2
		var detallesInversion []models.DetalleInversionInicial
		err := db.Where("plan_negocio_id = ? AND (tipo_id = 1 OR tipo_id = 2)", planID).
			Find(&detallesInversion).Error
		if err == nil {
			for _, detalle := range detallesInversion {
				noCorrientesSuma += detalle.Importe
			}
		}
	} else {
		// Mes 1 en adelante: NoCorrientes_suma(mes anterior) - depreciaciones mensuales acumuladas
		var balanceAnterior models.BalanceGeneral
		mesAnterior := mes - 1
		anioAnterior := anio
		if mes == 1 {
			if anio == 1 {
				// Si es mes 1 del año 1, buscar mes 0 del año 1
				mesAnterior = 0
				anioAnterior = 1
			} else {
				// Si es mes 1 de años 2-5, buscar mes 12 del año anterior
				mesAnterior = 12
				anioAnterior = anio - 1
			}
		}

		err := db.Where("plan_negocio_id = ? AND anio = ? AND mes = ?", planID, anioAnterior, mesAnterior).
			First(&balanceAnterior).Error
		if err == nil {
			noCorrientesSuma = balanceAnterior.NoCorrientes_Suma
		}

		// Restar depreciaciones mensuales del mes actual
		var depreciaciones []models.Depreciacion
		err = db.Where("plan_negocio_id = ?", planID).Find(&depreciaciones).Error
		if err == nil {
			for _, depreciacion := range depreciaciones {
				if depreciacion.DepreciacionMensual != nil {
					noCorrientesSuma -= *depreciacion.DepreciacionMensual
				}
			}
		}
	}

	// 7. Calcular TotalActivo
	totalActivo := corrientesSuma + noCorrientesSuma

	// 8. Calcular PasivoProveedoresCortoPlazo
	var pasivoProveedores float64
	if anio == 1 && mes == 0 {
		pasivoProveedores = 0 // Inicializa en 0 para mes 0 año 1
	} else {
		// Obtener pasivo proveedores del mes anterior
		var balanceAnterior models.BalanceGeneral
		mesAnterior := mes - 1
		anioAnterior := anio
		if mes == 1 {
			if anio == 1 {
				// Si es mes 1 del año 1, tomar mes 0 del año 1
				mesAnterior = 0
				anioAnterior = 1
			} else {
				// Si es mes 1 de años 2-5, tomar mes 12 del año anterior
				mesAnterior = 12
				anioAnterior = anio - 1
			}
		}

		err := db.Where("plan_negocio_id = ? AND anio = ? AND mes = ?", planID, anioAnterior, mesAnterior).
			First(&balanceAnterior).Error
		if err == nil {
			pasivoProveedores = balanceAnterior.PasivoProveedoresCortoPlazo
		}

		// Sumar costos de materia prima del mes actual * porcentaje a crédito
		var costosMateriasPrimas []models.CostoMateriasPrimas
		err = db.Where("plan_negocio_id = ? AND anio = ?", planID, anio).Find(&costosMateriasPrimas).Error
		if err == nil {
			// Obtener política de compra del mes actual
			var politicaCompra models.PoliticasCompra
			err = db.Where("plan_negocio_id = ? AND anio = ? AND mes = ?", planID, anio, mes).
				First(&politicaCompra).Error
			if err == nil {
				for _, costo := range costosMateriasPrimas {
					costosCredito := costo.CostoMensual * (politicaCompra.PorcentajeCredito / 100.0)
					pasivoProveedores += costosCredito
				}
			}
		}

		// Restar pagos de compras a crédito del flujo de efectivo
		var flujoEfectivo models.FlujoEfectivo
		err = db.Where("plan_negocio_id = ? AND anio = ? AND mes = ?", planID, anio, mes).
			First(&flujoEfectivo).Error
		if err == nil {
			pasivoProveedores -= flujoEfectivo.Egresos_ComprasCostosCredito
		}
	}

	// 9. Calcular PasivoPrestamosCortoPlazo
	var pasivoPrestamos float64
	if anio == 1 && mes == 0 {
		// Mes 0: suma de todas las amortizaciones del primer año (mes 1-12 del año 1)
		var prestamoCuotas []models.PrestamoCuotas
		err := db.Where("plan_negocio_id = ? AND anio = 1", planID).Find(&prestamoCuotas).Error
		if err == nil {
			for _, cuota := range prestamoCuotas {
				pasivoPrestamos += cuota.Amortizacion
			}
		}
	} else {
		// Mes 1 en adelante: PasivoPrestamosCortoPlazo(mes anterior) - amortización del mes actual
		var balanceAnterior models.BalanceGeneral
		mesAnterior := mes - 1
		anioAnterior := anio
		if mes == 1 {
			if anio == 1 {
				// Si es mes 1 del año 1, tomar mes 0 del año 1
				mesAnterior = 0
				anioAnterior = 1
			} else {
				// Si es mes 1 de años 2-5, tomar mes 12 del año anterior
				mesAnterior = 12
				anioAnterior = anio - 1
			}
		}

		err := db.Where("plan_negocio_id = ? AND anio = ? AND mes = ?", planID, anioAnterior, mesAnterior).
			First(&balanceAnterior).Error
		if err == nil {
			pasivoPrestamos = balanceAnterior.PasivoPrestamosCortoPlazo
		}

		// Restar amortización del mes actual
		var prestamoCuota models.PrestamoCuotas
		err = db.Where("plan_negocio_id = ? AND anio = ? AND mes = ?", planID, anio, mes).
			First(&prestamoCuota).Error
		if err == nil {
			pasivoPrestamos -= prestamoCuota.Amortizacion
		}
	}

	// 10. Calcular PasivoCuentasxPagarCortoPlazo
	var pasivoCuentasPorPagar float64
	if anio == 1 && mes == 0 {
		// Mes 0: inicializar en 0
		pasivoCuentasPorPagar = 0
	} else {
		// Mes 1 en adelante: PasivoCuentasxPagarCortoPlazo(mes anterior) + ISR(mes actual) - Egresos_PagosSRI(mes actual)
		var balanceAnterior models.BalanceGeneral
		mesAnterior := mes - 1
		anioAnterior := anio
		if mes == 1 {
			if anio == 1 {
				// Si es mes 1 del año 1, tomar mes 0 del año 1
				mesAnterior = 0
				anioAnterior = 1
			} else {
				// Si es mes 1 de años 2-5, tomar mes 12 del año anterior
				mesAnterior = 12
				anioAnterior = anio - 1
			}
		}

		err := db.Where("plan_negocio_id = ? AND anio = ? AND mes = ?", planID, anioAnterior, mesAnterior).
			First(&balanceAnterior).Error
		if err == nil {
			pasivoCuentasPorPagar = balanceAnterior.PasivoCuentasxPagarCortoPlazo
		}

		// Sumar ISR del mes actual
		var estadoResultados models.EstadoResultados
		err = db.Where("plan_negocio_id = ? AND anio = ? AND mes = ?", planID, anio, mes).
			First(&estadoResultados).Error
		if err == nil {
			pasivoCuentasPorPagar += estadoResultados.ISR
		}

		// Restar pagos SRI del mes actual
		var flujoEfectivo models.FlujoEfectivo
		err = db.Where("plan_negocio_id = ? AND anio = ? AND mes = ?", planID, anio, mes).
			First(&flujoEfectivo).Error
		if err == nil {
			pasivoCuentasPorPagar -= flujoEfectivo.Egresos_PagosSRI
		}
	}

	// 11. Calcular PasivoOtrosCortoPlazo
	var pasivoOtrosCortoPlazo float64
	if anio == 1 && mes == 0 {
		// Mes 0: inicializar en 0
		pasivoOtrosCortoPlazo = 0
	} else {
		// Mes 1 en adelante: PasivoOtrosCortoPlazo(mes anterior) + PTU(mes actual) - Egresos_PagoPTU(mes actual)
		var balanceAnterior models.BalanceGeneral
		mesAnterior := mes - 1
		anioAnterior := anio
		if mes == 1 {
			if anio == 1 {
				// Si es mes 1 del año 1, tomar mes 0 del año 1
				mesAnterior = 0
				anioAnterior = 1
			} else {
				// Si es mes 1 de años 2-5, tomar mes 12 del año anterior
				mesAnterior = 12
				anioAnterior = anio - 1
			}
		}

		err := db.Where("plan_negocio_id = ? AND anio = ? AND mes = ?", planID, anioAnterior, mesAnterior).
			First(&balanceAnterior).Error
		if err == nil {
			pasivoOtrosCortoPlazo = balanceAnterior.PasivoOtrosCortoPlazo
		}

		// Sumar PTU del mes actual
		var estadoResultados models.EstadoResultados
		err = db.Where("plan_negocio_id = ? AND anio = ? AND mes = ?", planID, anio, mes).
			First(&estadoResultados).Error
		if err == nil {
			pasivoOtrosCortoPlazo += estadoResultados.PTU
		}

		// Restar pagos PTU del mes actual
		var flujoEfectivo models.FlujoEfectivo
		err = db.Where("plan_negocio_id = ? AND anio = ? AND mes = ?", planID, anio, mes).
			First(&flujoEfectivo).Error
		if err == nil {
			pasivoOtrosCortoPlazo -= flujoEfectivo.Egresos_PagoPTU
		}
	}

	// 12. Calcular PasivoCortoPlazo_Suma
	pasivoCortoPlazoSuma := pasivoProveedores + pasivoPrestamos + pasivoCuentasPorPagar + pasivoOtrosCortoPlazo

	// Actualizar el registro de balance
	updates := map[string]interface{}{
		"corrientes_efectivo":               efectivo,
		"corrientes_cuentasx_cobrar":        cuentasPorCobrar,
		"corrientes_inventarios":            inventarios,
		"corrientes_otros":                  corrientesOtros,
		"corrientes_suma":                   corrientesSuma,
		"no_corrientes_suma":                noCorrientesSuma,
		"total_activo":                      totalActivo,
		"pasivo_proveedores_corto_plazo":    pasivoProveedores,
		"pasivo_prestamos_corto_plazo":      pasivoPrestamos,
		"pasivo_cuentasx_pagar_corto_plazo": pasivoCuentasPorPagar,
		"pasivo_otros_corto_plazo":          pasivoOtrosCortoPlazo,
		"pasivo_corto_plazo_suma":           pasivoCortoPlazoSuma,
	}

	if err := db.Model(&balance).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}
