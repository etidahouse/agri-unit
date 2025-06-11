package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"weather-ingestor/weather"
)

type App struct {
	AgriUnitStorage weather.AgriUnitStorage
	WeatherStorage  weather.WeatherStorage
	apiURL          string
	apiKey          string
}

func (a *App) IngestionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed. Only POST method is supported.", http.StatusMethodNotAllowed)
		return
	}

	fetcher := weather.NewWeatherFetcher(a.apiURL, a.apiKey, a.WeatherStorage, a.AgriUnitStorage)

	err := fetcher.HandleWeatherIngest()

	if err != nil {
		log.Printf("Error during ingestion: %v\n", err)
		http.Error(w, fmt.Sprintf("Error processing ingestion: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data weather ingestion started successfully.\n"))
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

	realAgriUnitStorage := weather.NewAgriUnitStorage(db)
	realWeatherStorage := weather.NewWeatherStorage(db)

	apiUrl := os.Getenv("API_URL")
	apiKey := os.Getenv("API_KEY")

	app := &App{
		AgriUnitStorage: realAgriUnitStorage,
		WeatherStorage:  realWeatherStorage,
		apiURL:          apiUrl,
		apiKey:          apiKey,
	}

	http.HandleFunc("/ingest", app.IngestionHandler)

	port := ":8080"
	log.Printf("Server started on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
