package main

import (
	"log"
	"net/http"
	"os"

	"github.com/FraMan97/kairos-p2p-ledger/internal/api"
	"github.com/FraMan97/kairos-p2p-ledger/internal/config"
	"github.com/FraMan97/kairos-p2p-ledger/internal/database"
	"github.com/FraMan97/kairos-p2p-ledger/internal/worker"
)

func main() {
	config.InitConfig()

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "/data"
	}

	database.InitDB(dbPath)
	defer database.CloseDB()

	worker.StartUpgrader()

	mux := http.NewServeMux()
	mux.HandleFunc("/anchor", api.AnchorHash)
	mux.HandleFunc("/receipt", api.GetReceipt)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("[Main] Kairos Ledger Server listening on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("[Main] Server failed: %v", err)
	}
}
