package main

import (
	"log"
	"net/http"

	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/db"
	"github.com/JostinAlvaradoS/liveplan_backend_go/internal/handlers"
)

func main() {
	dbconn, err := db.Connect()
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	// close underlying sql DB when program exits
	sqlDB, err := dbconn.DB()
	if err != nil {
		log.Fatalf("failed to get sql DB: %v", err)
	}
	defer sqlDB.Close()

	if err := db.Migrate(dbconn); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	app := &handlers.App{DB: dbconn}
	mux := app.Routes()

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
