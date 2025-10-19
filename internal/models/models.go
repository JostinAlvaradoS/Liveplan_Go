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
