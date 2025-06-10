package main

import (
	"agreste-ingestor/agri_units"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type IngestionRequest struct {
	ZipURL      string `json:"zipUrl"`
	CSVFileName string `json:"csvFileName"`
}

type App struct {
	AgriUnitStorage       agri_units.AgriUnitStorage
	AgriUnitSurveyStorage agri_units.AgriculturalUnitSurveyStorage
}

func (a *App) IngestionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed. Only POST method is supported.", http.StatusMethodNotAllowed)
		return
	}

	var req IngestionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading JSON request body: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.ZipURL == "" || req.CSVFileName == "" {
		http.Error(w, "Fields 'zipUrl' and 'csvFileName' are required.", http.StatusBadRequest)
		return
	}

	log.Printf("Ingestion request received: ZipURL='%s', CSVFileName='%s'\n", req.ZipURL, req.CSVFileName)

	err = agri_units.HandleAgriUnitSurveyIngest(
		req.ZipURL,
		req.CSVFileName,
		a.AgriUnitStorage,
		a.AgriUnitSurveyStorage,
	)

	if err != nil {
		log.Printf("Error during ingestion: %v\n", err)
		http.Error(w, fmt.Sprintf("Error processing ingestion: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data ingestion started successfully.\n"))
	log.Println("Ingestion completed successfully.")
}

func main() {

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatalf("Error: One or more database connection environment variables are missing. Ensure DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, and DB_NAME are defined.")
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open database connection with URL '%s': %v", dbURL, err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("Database connection failed (ping failed) with URL '%s': %v", dbURL, err)
	}
	log.Println("PostgreSQL database connection established successfully.")

	realAgriUnitStorage := agri_units.NewAgriUnitStorage(db)
	realAgriUnitSurveyStorage := agri_units.NewAgriculturalUnitSurveyStorage(db)

	app := &App{
		AgriUnitStorage:       realAgriUnitStorage,
		AgriUnitSurveyStorage: realAgriUnitSurveyStorage,
	}

	http.HandleFunc("/ingest", app.IngestionHandler)

	port := ":8080"
	log.Printf("Server started on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
