package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/stockyard-dev/stockyard-notebook/internal/server"
	"github.com/stockyard-dev/stockyard-notebook/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9370"
	}
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./notebook-data"
	}

	db, err := store.Open(dataDir)
	if err != nil {
		log.Fatalf("notebook: open database: %v", err)
	}
	defer db.Close()

	srv := server.New(db, server.DefaultLimits())

	fmt.Printf("\n  Notebook — Self-hosted personal notes\n")
	fmt.Printf("  ─────────────────────────────────\n")
	fmt.Printf("  Dashboard:  http://localhost:%s/ui\n", port)
	fmt.Printf("  API:        http://localhost:%s/api\n", port)
	fmt.Printf("  Data:       %s\n", dataDir)
	fmt.Printf("  ─────────────────────────────────\n\n")

	log.Printf("notebook: listening on :%s", port)
	if err := http.ListenAndServe(":"+port, srv); err != nil {
		log.Fatalf("notebook: %v", err)
	}
}
