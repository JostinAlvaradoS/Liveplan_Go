package db

import (
	"fmt"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "liveplan"
)

func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return gdb, nil
}

// AutoMigrate models
func Migrate(gdb *gorm.DB) error {
	return gdb.AutoMigrate(
		&models.PlanNegocio{},
		&models.TipoInversionInicial{},
		&models.InversionInicial{},
		&models.DetalleInversionInicial{},
		&models.ProductoServicio{},
		&models.Supuesto{},
		&models.VentaDiaria{},
		&models.VariablesDeSensibilidad{},
		&models.VariacionAnual{},
		&models.PreciosProdServ{},
		&models.CategoriaCosto{},
		&models.CostosProdServ{},
		&models.IndicadoresMacro{},
		&models.ComposicionFinanciamiento{},
		&models.Depreciacion{},
		&models.PresupuestoVenta{},
	)
}
