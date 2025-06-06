# Agri-Unit Data Pipeline

This project demonstrates a data pipeline built to ingest, process, and store agricultural (Agreste) and weather (OpenWeather) data. It showcases a modern, scalable approach for handling diverse data streams.

---

## Project Highlights

* **Go (Golang):** Chosen for its performance, efficiency, and suitability for building robust data ingestion and transformation services.
* **Docker:** All services are containerized using **Docker**, ensuring consistent development and deployment environments.
* **PostgreSQL:** Serves as the data warehouse for storing the structured agricultural and weather data.

---

## What It Does

The pipeline performs the following key functions:

* **Ingests Agreste Data:** Processes agricultural production data, anonymizing farm locations with random geographical coordinates and filtering for cereal production.
* **Ingests OpenWeather Data:** Fetches and processes real-time weather information.
* **Transforms & Stores:** Processes the raw data and loads it into a **PostgreSQL** data warehouse, ready for analysis.

---