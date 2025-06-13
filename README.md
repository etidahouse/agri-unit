# Agri-Unit Data Pipeline

This project demonstrates a data pipeline built to ingest, process, and store agricultural (Agreste) and weather (OpenWeather) data. It showcases a modern, scalable approach for handling diverse data streams.

---

## Project Highlights

- **Go (Golang):** Chosen for its performance, efficiency, and suitability for building robust data ingestion and transformation services.
- **Docker:** All services are containerized using **Docker**, ensuring consistent development and deployment environments.
- **PostgreSQL:** Serves as the data warehouse for storing structured agricultural and weather data.
- **Streamlit:** An interactive dashboard for **data visualization**, enabling fast exploration and insight into agricultural and weather trends.
- **CI Integration:** Containers are **automatically built and released on each merge** to the main branch via a continuous integration pipeline (GitHub Actions).

---

## What It Does

The pipeline performs the following key functions:

- **Ingests Agreste Data:** Processes agricultural production data, anonymizes farm locations using random geographic coordinates, and filters for cereal production.
- **Ingests OpenWeather Data:** Fetches and processes real-time weather information.
- **Transforms & Stores:** Cleans and loads raw data into a **PostgreSQL** data warehouse for analysis.
- **Data Visualization:** Provides an intuitive **Streamlit UI** to view farm weather on a map, explore weather history, and analyze cereal yield by farm.

---

## Getting Started

To launch all pipeline services (ingestion, database, and Streamlit interface), follow these steps:

1. Ensure **Docker** and **Docker Compose** are installed on your machine.

2. Navigate to the `deploy/` directory at the root of the project:

    ```bash
    cd deploy/
    ```

3. Launch all services using Docker Compose:

    ```bash
    docker compose up
    ```

Once services are running:

- Access the **Streamlit app** via [http://localhost:8501](http://localhost:8501).
- The **Agreste** and **OpenWeather** ingestion services run in the background and are ready to be triggered.

---

## Automated Scheduling with Cron

The pipeline includes cron jobs to **automatically trigger ingestion and transformation tasks** at regular intervals, ensuring up-to-date data without manual intervention.

---

## Manually Triggering Ingestion (Optional)

You can also manually trigger the ingestion processes for Agreste and OpenWeather data using `curl`:

- **Trigger Agreste (Rica) data ingestion:**

    ```bash
    curl -X POST http://localhost:8080/ingest \
      -H "Content-Type: application/json" \
      -d '{
        "zipUrl": "https://agreste.agriculture.gouv.fr/agreste-web/download/service/SV-Accès micro données RICA/RicaMicrodonnées2023_v2.zip",
        "csvFileName": "Rica_France_micro_Donnees_ex2023.csv"
      }'
    ```

- **Trigger OpenWeather data ingestion:**

    ```bash
    curl -X POST http://localhost:8081/ingest
    ```

---

## Accessing the PostgreSQL Database

To connect directly to the PostgreSQL instance running in Docker:

1. Find the container name or ID:

    ```bash
    docker ps
    ```

2. Connect to the database using `psql`:

    ```bash
    docker exec -it <postgresql_container_name_or_id> psql -U postgres -d mydatabase
    ```

---

## 🔁 Continuous Integration & Deployment

Every time code is merged into the main branch:

- Containers are **automatically built and published** using CI (GitHub Actions).
- The latest version of each service is deployed consistently, ensuring fast iteration and reliable updates across environments.

---
