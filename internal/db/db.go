package db

import (
	"fmt"
	"os"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func Connect() (*gorm.DB, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "liveplan")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
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
		&models.DatosPrestamo{},
		&models.PrestamoCuotas{},
		&models.VentasDinero{},
	)
}
