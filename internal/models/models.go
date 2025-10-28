package models


// PlanNegocio representa la tabla plan_negocio
type PlanNegocio struct {
	ID           uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Autor        string `json:"autor" gorm:"type:varchar(90);not null"`
	Problematica string `json:"problematica" gorm:"type:varchar(300);not null"`
	Descripcion  string `json:"descripcion" gorm:"type:text"`
}

// TipoInversionInicial representa la tabla tipo_inversion_inicial
type TipoInversionInicial struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Tipo string `json:"tipo" gorm:"type:varchar(100);not null"`
}

// InversionInicial representa la tabla inversion_inicial
type InversionInicial struct {
	ID            uint                  `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanNegocioID uint                  `json:"plan_negocio_id" gorm:"not null;index"`
	PlanNegocio   *PlanNegocio          `json:"plan_negocio,omitempty" gorm:"constraint:OnDelete:CASCADE;foreignKey:PlanNegocioID"`
	TipoID        uint                  `json:"tipo_id" gorm:"not null;index"`
	Tipo          *TipoInversionInicial `json:"tipo,omitempty" gorm:"foreignKey:TipoID"`
	Seccion       string                `json:"seccion" gorm:"type:varchar(100);not null"`
	Importe       float64               `json:"importe" gorm:"type:numeric(15,2);not null"`
}

// DetalleInversionInicial representa la tabla detalle_inversion_inicial
type DetalleInversionInicial struct {
	ID            uint                  `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanNegocioID uint                  `json:"plan_negocio_id" gorm:"not null;index"`
	PlanNegocio   *PlanNegocio          `json:"plan_negocio,omitempty" gorm:"constraint:OnDelete:CASCADE;foreignKey:PlanNegocioID"`
	InversionID   uint                  `json:"inversion_id" gorm:"not null;index"`
	Inversion     *InversionInicial     `json:"inversion,omitempty" gorm:"constraint:OnDelete:CASCADE;foreignKey:InversionID"`
	TipoID        uint                  `json:"tipo_id" gorm:"not null;index"`
	Tipo          *TipoInversionInicial `json:"tipo,omitempty" gorm:"foreignKey:TipoID"`
	Elemento      string                `json:"elemento" gorm:"type:varchar(100);not null"`
	Importe       float64               `json:"importe" gorm:"type:numeric(15,2);not null"`
	VidaUtil      int                   `json:"vida_util"`
}

// ProductoServicio representa un producto o servicio asociado a un plan de negocio
type ProductoServicio struct {
	ID            uint         `json:"id" gorm:"primaryKey;autoIncrement"`
	Nombre        string       `json:"nombre" gorm:"type:varchar(150);not null"`
	PlanNegocioID uint         `json:"plan_negocio_id" gorm:"not null;index"`
	PlanNegocio   *PlanNegocio `json:"plan_negocio,omitempty" gorm:"foreignKey:PlanNegocioID;constraint:OnDelete:CASCADE"`
}

// Supuesto representa supuestos financieros asociados a un plan de negocio
type Supuesto struct {
	ID                    uint         `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanNegocioID         uint         `json:"plan_negocio_id" gorm:"not null;index"`
	PorcenVentas          float64      `json:"porcen_ventas" gorm:"type:numeric(6,2)"`
	VariacionPorcenVentas float64      `json:"variacion_porcen_ventas" gorm:"type:numeric(6,2)"`
	PTU                   float64      `json:"ptu" gorm:"type:numeric(6,2)"`
	ISR                   float64      `json:"isr" gorm:"type:numeric(6,2)"`
	PlanNegocio           *PlanNegocio `json:"plan_negocio,omitempty" gorm:"foreignKey:PlanNegocioID;constraint:OnDelete:CASCADE"`
}

// VentaDiaria representa las ventas diarias por producto asociado a un plan de negocio
type VentaDiaria struct {
	ID                 uint              `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanNegocioID      uint              `json:"plan_negocio_id" gorm:"not null;index"`
	ProductoServicioID uint              `json:"producto_servicio_id" gorm:"not null;index"`
	VentaDia           *int              `json:"venta_dia" gorm:"column:venta_dia"`
	PlanNegocio        *PlanNegocio      `json:"plan_negocio,omitempty" gorm:"foreignKey:PlanNegocioID;constraint:OnDelete:CASCADE"`
	ProductoServicio   *ProductoServicio `json:"producto_servicio,omitempty" gorm:"foreignKey:ProductoServicioID;constraint:OnDelete:CASCADE"`
}

type VariablesDeSensibilidad struct {
	ID               uint         `json:"id" gorm:"primaryKey;autoIncrement"`
	Cantidad_volumen float64      `json:"cantidad_volumen" gorm:"type:numeric(15,2);"`
	Precio           float64      `json:"precio" gorm:"type:numeric(15,2);"`
	Costo            float64      `json:"costo" gorm:"type:numeric(15,2);"`
	PlanNegocioID    uint         `json:"plan_negocio_id" gorm:"not null;index"`
	PlanNegocio      *PlanNegocio `json:"plan_negocio,omitempty" gorm:"foreignKey:PlanNegocioID;constraint:OnDelete:CASCADE"`
}

// VariacionAnual almacena el porcentaje de crecimiento de ventas por año (año 1..5) para un plan
type VariacionAnual struct {
	ID            uint         `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanNegocioID uint         `json:"plan_negocio_id" gorm:"not null;index"`
	Año1          float64      `json:"anio1" gorm:"column:anio1;type:numeric(6,2)"`
	Año2          float64      `json:"anio2" gorm:"column:anio2;type:numeric(6,2)"`
	Año3          float64      `json:"anio3" gorm:"column:anio3;type:numeric(6,2)"`
	Año4          float64      `json:"anio4" gorm:"column:anio4;type:numeric(6,2)"`
	Año5          float64      `json:"anio5" gorm:"column:anio5;type:numeric(6,2)"`
	PlanNegocio   *PlanNegocio `json:"plan_negocio,omitempty" gorm:"foreignKey:PlanNegocioID;constraint:OnDelete:CASCADE"`
}

// PreciosProdServ almacena precio por producto_servicio dentro de un plan y un precio calculado
type PreciosProdServ struct {
	ID                 uint              `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanNegocioID      uint              `json:"plan_negocio_id" gorm:"not null;index"`
	ProductoServicioID uint              `json:"producto_servicio_id" gorm:"not null;index"`
	Precio             *float64          `json:"precio" gorm:"type:numeric(15,2)"`
	PrecioCalc         *float64          `json:"precio_calc" gorm:"type:numeric(15,2)"`
	PlanNegocio        *PlanNegocio      `json:"plan_negocio,omitempty" gorm:"foreignKey:PlanNegocioID;constraint:OnDelete:CASCADE"`
	ProductoServicio   *ProductoServicio `json:"producto_servicio,omitempty" gorm:"foreignKey:ProductoServicioID;constraint:OnDelete:CASCADE"`
}

// CategoriaCosto es un catálogo de categorías de costo
type CategoriaCosto struct {
	ID     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Nombre string `json:"nombre" gorm:"type:varchar(150);not null"`
}

// CostosProdServ almacena costos asociados a un producto/servicio dentro de un plan
type CostosProdServ struct {
	ID                 uint              `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanNegocioID      uint              `json:"plan_negocio_id" gorm:"not null;index"`
	ProductoServicioID uint              `json:"producto_servicio_id" gorm:"not null;index"`
	CategoriaCostoID   uint              `json:"categoria_costo_id" gorm:"not null;index"`
	Costo              *float64          `json:"costo" gorm:"type:numeric(15,2)"`
	CostoCalc          *float64          `json:"costo_calc" gorm:"type:numeric(15,2)"`
	PlanNegocio        *PlanNegocio      `json:"plan_negocio,omitempty" gorm:"foreignKey:PlanNegocioID;constraint:OnDelete:CASCADE"`
	ProductoServicio   *ProductoServicio `json:"producto_servicio,omitempty" gorm:"foreignKey:ProductoServicioID;constraint:OnDelete:CASCADE"`
	CategoriaCosto     *CategoriaCosto   `json:"categoria_costo,omitempty" gorm:"foreignKey:CategoriaCostoID;constraint:OnDelete:CASCADE"`
}

// IndicadoresMacro almacena indicadores macroeconómicos asociados a un plan
type IndicadoresMacro struct {
	ID            uint         `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanNegocioID uint         `json:"plan_negocio_id" gorm:"not null;index"`
	TipoCambio    float64      `json:"tipo_cambio" gorm:"column:tipo_cambio;type:numeric(15,6)"`
	Inflacion     float64      `json:"inflacion" gorm:"type:numeric(6,2)"`
	TasaDeuda     float64      `json:"tasa_deuda" gorm:"column:tasa_deuda;type:numeric(6,2)"`
	TasaInteres   float64      `json:"tasa_interes" gorm:"column:tasa_interes;type:numeric(6,2)"`
	TasaImpuesto  float64      `json:"tasa_impuesto" gorm:"column:tasa_impuesto;type:numeric(6,2)"`
	PTU           float64      `json:"ptu" gorm:"type:numeric(6,2)"`
	DiasxMes      int          `json:"diasxmes" gorm:"column:diasxmes"`
	PlanNegocio   *PlanNegocio `json:"plan_negocio,omitempty" gorm:"foreignKey:PlanNegocioID;constraint:OnDelete:CASCADE"`
}

// ComposicionFinanciamiento representa la composición de financiamiento para un plan
type ComposicionFinanciamiento struct {
	ID                uint         `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanNegocioID     uint         `json:"plan_negocio_id" gorm:"not null;index"`
	CapitalPorcentaje float64      `json:"capital_porcentaje" gorm:"column:capital_porcentaje;type:numeric(6,2)"`
	DeudaPorcentaje   float64      `json:"deuda_porcentaje" gorm:"column:deuda_porcentaje;type:numeric(6,2)"`
	Total_Inversion   float64      `json:"total_inversion" gorm:"column:total_inversion;type:numeric(15,2)"`
	PlanNegocio       *PlanNegocio `json:"plan_negocio,omitempty" gorm:"foreignKey:PlanNegocioID;constraint:OnDelete:CASCADE"`
}

// Depreciacion representa las depreciaciones calculadas a partir de un detalle de inversión inicial
type Depreciacion struct {
	ID                  uint     `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanNegocioID       uint     `json:"plan_negocio_id" gorm:"not null;index"`
	DetalleInversionID  uint     `json:"detalle_inversion_id" gorm:"not null;index"`
	DepreciacionMensual *float64 `json:"depreciacion_mensual" gorm:"column:depreciacion_mensual;type:numeric(15,2)"`
	DepreciacionAnio1   *float64 `json:"depreciacion_anio1" gorm:"column:depreciacion_anio1;type:numeric(15,2)"`
	DepreciacionAnio2   *float64 `json:"depreciacion_anio2" gorm:"column:depreciacion_anio2;type:numeric(15,2)"`
	DepreciacionAnio3   *float64 `json:"depreciacion_anio3" gorm:"column:depreciacion_anio3;type:numeric(15,2)"`
	DepreciacionAnio4   *float64 `json:"depreciacion_anio4" gorm:"column:depreciacion_anio4;type:numeric(15,2)"`
	DepreciacionAnio5   *float64 `json:"depreciacion_anio5" gorm:"column:depreciacion_anio5;type:numeric(15,2)"`
	ValorRescate        *float64 `json:"valor_rescate" gorm:"column:valor_rescate;type:numeric(15,2)"`

	PlanNegocio      *PlanNegocio             `json:"plan_negocio,omitempty" gorm:"foreignKey:PlanNegocioID;constraint:OnDelete:CASCADE"`
	DetalleInversion *DetalleInversionInicial `json:"detalle_inversion,omitempty" gorm:"foreignKey:DetalleInversionID;constraint:OnDelete:CASCADE"`
}

// PresupuestoVenta representa el presupuesto de ventas para un producto dentro de un plan
type PresupuestoVenta struct {
	ID            uint              `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanNegocioID uint              `json:"plan_negocio_id" gorm:"not null;index"`
	ProductoID    uint              `json:"producto_id" gorm:"not null;index"`
	Anio          int               `json:"anio" gorm:"not null;index"`
	Crecimiento   *float64          `json:"crecimiento" gorm:"type:numeric(6,2)"`
	Mensual       *float64          `json:"mensual" gorm:"type:numeric(15,2)"`
	Anual         *float64          `json:"anual" gorm:"type:numeric(15,2)"`
	PlanNegocio   *PlanNegocio      `json:"plan_negocio,omitempty" gorm:"foreignKey:PlanNegocioID;constraint:OnDelete:CASCADE"`
	Producto      *ProductoServicio `json:"producto,omitempty" gorm:"foreignKey:ProductoID;constraint:OnDelete:CASCADE"`
}

type DatosPrestamo struct {
	ID                     uint         `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanNegocioID          uint         `json:"plan_negocio_id" gorm:"not null;index"`
	Monto                  float64      `json:"monto" gorm:"type:numeric(15,2);not null"`
	TasaAnual              float64      `json:"tasa_anual" gorm:"column:tasa_anual;type:numeric(6,2);not null"`
	PeriodosCapitalizacion int          `json:"periodos_capitalizacion" gorm:"column:periodos_capitalizacion;not null"`
	TasaMensual            float64      `json:"tasa_mensual" gorm:"column:tasa_mensual;type:numeric(6,4);not null"`
	Cuota                  float64      `json:"cuota" gorm:"column:cuota;type:numeric(15,2);not null"`
	PeriodosAmortizacion   int          `json:"periodos_amortizacion" gorm:"column:periodos_amortizacion;not null"`
	PlanNegocio            *PlanNegocio `json:"plan_negocio,omitempty" gorm:"foreignKey:PlanNegocioID;constraint:OnDelete:CASCADE"`
}

type PrestamoCuotas struct {
	ID             uint         `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanNegocioID  uint         `json:"plan_negocio_id" gorm:"not null;index"`
	SaldoInicial   float64      `json:"saldo_inicial" gorm:"column:saldo_inicial;type:numeric(15,2);not null"`
	PeriodoMes     int          `json:"periodo_mes" gorm:"column:periodo_mes;not null"`
	Anio           int          `json:"anio" gorm:"column:anio;not null"`
	Mes            int          `json:"mes" gorm:"column:mes;not null"`
	Interes        float64      `json:"interes" gorm:"type:numeric(15,2);not null"`
	Amortizacion   float64      `json:"amortizacion" gorm:"type:numeric(15,2);not null"`
	CuotaTotal     float64      `json:"cuota_total" gorm:"column:cuota_total;type:numeric(15,2);not null"`
	SaldoPendiente float64      `json:"saldo_pendiente" gorm:"column:saldo_pendiente;type:numeric(15,2);not null"`
	PlanNegocio    *PlanNegocio `json:"plan_negocio,omitempty" gorm:"foreignKey:PlanNegocioID;constraint:OnDelete:CASCADE"`
}

type VentasDinero struct {
	ID            uint              `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanNegocioID uint              `json:"plan_negocio_id" gorm:"not null;index"`
	ProductoID    uint              `json:"producto_id" gorm:"not null;index"`
	Anio          int               `json:"anio" gorm:"not null;index"`
	Mensual       float64           `json:"mensual" gorm:"not null;index"`
	Anual         float64           `json:"anual" gorm:"not null;index"`
	PlanNegocio   *PlanNegocio      `json:"plan_negocio,omitempty" gorm:"foreignKey:PlanNegocioID;constraint:OnDelete:CASCADE"`
	Producto      *ProductoServicio `json:"producto,omitempty" gorm:"foreignKey:ProductoID;constraint:OnDelete:CASCADE"`
}

type Ventas struct {
	ID            uint              `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanNegocioID uint              `json:"plan_negocio_id" gorm:"not null;index"`
	ProductoID    uint              `json:"producto_id" gorm:"not null;index"`
	Anio          int               `json:"anio" gorm:"not null;index"`
	Venta         float64           `json:"venta" gorm:"not null;index"`
	PlanNegocio   *PlanNegocio      `json:"plan_negocio,omitempty" gorm:"foreignKey:PlanNegocioID;constraint:OnDelete:CASCADE"`
	Producto      *ProductoServicio `json:"producto,omitempty" gorm:"foreignKey:ProductoID;constraint:OnDelete:CASCADE"`
}
type CostosVentas struct {
	ID            uint              `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanNegocioID uint              `json:"plan_negocio_id" gorm:"not null;index"`
	ProductoID    uint              `json:"producto_id" gorm:"not null;index"`
	Anio          int               `json:"anio" gorm:"not null;index"`
	Mes           int               `json:"mes" gorm:"not null;index"`
	Costo         float64           `json:"costo" gorm:"not null;index"`
	PlanNegocio   *PlanNegocio      `json:"plan_negocio,omitempty" gorm:"foreignKey:PlanNegocioID;constraint:OnDelete:CASCADE"`
	Producto      *ProductoServicio `json:"producto,omitempty" gorm:"foreignKey:ProductoID;constraint:OnDelete:CASCADE"`
}

type CostoMateriasPrimas struct {
	ID            uint              `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanNegocioID uint              `json:"plan_negocio_id" gorm:"not null;index"`
	ProductoID    uint              `json:"producto_id" gorm:"not null;index"`
	Anio          int               `json:"anio" gorm:"not null;index"`
	CostoMensual  float64           `json:"costo_mensual" gorm:"not null;index"`
	CostoAnual    float64           `json:"costo_anual" gorm:"not null;index"`
	PlanNegocio   *PlanNegocio      `json:"plan_negocio,omitempty" gorm:"foreignKey:PlanNegocioID;constraint:OnDelete:CASCADE"`
	Producto      *ProductoServicio `json:"producto,omitempty" gorm:"foreignKey:ProductoID;constraint:OnDelete:CASCADE"`
}

type GastosOperacionBase struct {
	ID            uint         `json:"id" gorm:"primaryKey;autoIncrement"`
	Descripcion   string       `json:"descripcion" gorm:"type:varchar(200);not null"`
	Valor         float64      `json:"valor" gorm:"not null;index"`
}
type GastosOperacion struct {
	ID            uint         `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanNegocioID uint         `json:"plan_negocio_id" gorm:"not null;index"`
	Descripcion   string       `json:"descripcion" gorm:"type:varchar(200);not null"`
	Mensual      float64      `json:"mensual" gorm:"not null;index"`
	Anual        float64      `json:"anual" gorm:"not null;index"`
	PlanNegocio   *PlanNegocio `json:"plan_negocio,omitempty" gorm:"foreignKey:PlanNegocioID;constraint:OnDelete:CASCADE"`
}

type PoliticasVenta 		struct {
	ID            uint         `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanNegocioID uint         `json:"plan_negocio_id" gorm:"not null;index"`
	PorcentajeCredito float64      `json:"porcentaje_credito" gorm:"not null;index"`
	PorcentajeContado float64      `json:"porcentaje_contado" gorm:"not null;index"`
}

type PoliticasCompra		struct {
	ID            uint         `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanNegocioID uint         `json:"plan_negocio_id" gorm:"not null;index"`
	PorcentajeCredito float64      `json:"porcentaje_credito" gorm:"not null;index"`
	PorcentajeContado float64      `json:"porcentaje_contado" gorm:"not null;index"`
}

type EstadoResultados struct {
	ID                     uint         `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanNegocioID          uint         `json:"plan_negocio_id" gorm:"not null;index"`
	Anio                   int          `json:"anio" gorm:"not null;index"`
	Mes                    int          `json:"mes" gorm:"not null;index"`
	Ventas                 float64      `json:"ventas" gorm:"not null;index"`
	CostosVentas           float64      `json:"costos_ventas" gorm:"not null;index"`
	UtilidadBruta          float64      `json:"utilidad_bruta" gorm:"not null;index"`
	GastosVentaAdm         float64      `json:"gastos_venta_adm" gorm:"not null;index"`
	Depreciacion           float64      `json:"depreciacion" gorm:"not null;index"`
	Amortizacion           float64      `json:"amortizacion" gorm:"not null;index"`
	UtilidadprevioIntImp   float64      `json:"utilidad_previo_int_imp" gorm:"not null;index"`
	GastosFinancieros      float64      `json:"gastos_financieros" gorm:"not null;index"`
	UtilidadAntesPTU       float64      `json:"utilidad_antes_ptu" gorm:"not null;index"`
	PTU                    float64      `json:"ptu" gorm:"not null;index"`
	UtilidadAntesImpuestos float64      `json:"utilidad_antes_impuestos" gorm:"not null;index"`
	ISR                    float64      `json:"isr" gorm:"not null;index"`
	UtilidadNeta           float64      `json:"utilidad_neta" gorm:"not null;index"`
	PlanNegocio            *PlanNegocio `json:"plan_negocio,omitempty" gorm:"foreignKey:PlanNegocioID;constraint:OnDelete:CASCADE"`
}


type FlujoEfectivo struct {
	ID                           uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanNegocioID                uint    `json:"plan_negocio_id" gorm:"not null;index"`
	Anio                         int     `json:"anio" gorm:"not null;index"`
	Mes                          int     `json:"mes" gorm:"not null;index"`
	Ingresos_VentaContado        float64 `json:"ingresos_venta_contado" gorm:"not null;index"`
	Ingresos_CobrosVentasCredito float64 `json:"ingresos_cobros_ventas_credito" gorm:"not null;index"`
	Ingresos_OtrosIngresos       float64 `json:"ingresos_otros_ingresos" gorm:"not null;index"`
	Ingresos_Prestamos           float64 `json:"ingresos_prestamos" gorm:"not null;index"`
	Ingresos_AportesCapital      float64 `json:"ingresos_aportes_capital" gorm:"not null;index"`
	Egresos_ComprasCostosContado float64 `json:"egresos_compras_costos_contado" gorm:"not null;index"`
	Egresos_ComprasCostosCredito float64 `json:"egresos_compras_costos_credito" gorm:"not null;index"`
	Egresos_Intereses            float64 `json:"egresos_intereses" gorm:"not null;index"`
	Egresos_PagosPrestamos        float64 `json:"egresos_pagos_prestamos" gorm:"not null;index"`
	Egresos_PagosSRI			float64 `json:"egresos_pagos_sri" gorm:"not null;index"`
	Egresos_PagoPTU			float64 `json:"egresos_pago_ptu" gorm:"not null;index"`
	AumentoInventarios           float64 `json:"aumento_inventarios" gorm:"not null;index"`
	FlujoCaja                    float64 `json:"flujo_caja" gorm:"not null;index"`
	EfectivoInicial             float64 `json:"efectivo_inicial" gorm:"not null;index"`
	EfectivoFinal               float64 `json:"efectivo_final" gorm:"not null;index"`
	PlanNegocio                  *PlanNegocio `json:"plan_negocio,omitempty" gorm:"foreignKey:PlanNegocioID;constraint:OnDelete:CASCADE"`
}

type BalanceGeneral struct {
	ID                 uint         `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanNegocioID      uint         `json:"plan_negocio_id" gorm:"not null;index"`
	Anio               int          `json:"anio" gorm:"not null;index"`
	Mes                int          `json:"mes" gorm:"not null;index"`
	Corrientes_Efectivo    float64      `json:"corrientes_efectivo" gorm:"not null;index"`
	Corrientes_CuentasxCobrar float64      `json:"corrientes_cuentasx_cobrar" gorm:"not null;index"`
	Corrientes_Inventarios    float64      `json:"corrientes_inventarios" gorm:"not null;index"`
	Corrientes_Otros         float64      `json:"corrientes_otros_activos" gorm:"not null;index"`
	Corrientes_Suma		  float64      `json:"corrientes_suma" gorm:"not null;index"`
	NoCorrientes_Suma	   float64      `json:"no_corrientes_suma" gorm:"not null;index"`
	
	PlanNegocio              *PlanNegocio `json:"plan_negocio,omitempty" gorm:"foreignKey:PlanNegocioID;constraint:OnDelete:CASCADE"`
}