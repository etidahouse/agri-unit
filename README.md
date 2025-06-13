# Agri-Unit Data Pipeline

This project demonstrates a data pipeline built to ingest, process, and store agricultural (Agreste) and weather (OpenWeather) data. It showcases a modern, scalable approach for handling diverse data streams.

---

## Project Highlights

* **Go (Golang):** Chosen for its performance, efficiency, and suitability for building robust data ingestion and transformation services.

* **Docker:** All services are containerized using **Docker**, ensuring consistent development and deployment environments.

* **PostgreSQL:** Serves as the data warehouse for storing the structured agricultural and weather data.

* **Streamlit:** An interactive dashboard has been added for **data visualization**, allowing for easy exploration and quick analysis of agricultural and weather information.

---

## What It Does

The pipeline performs the following key functions:

* **Ingests Agreste Data:** Processes agricultural production data, anonymizing farm locations with random geographical coordinates and filtering for cereal production.

* **Ingests OpenWeather Data:** Fetches and processes real-time weather information.

* **Transforms & Stores:** Processes the raw data and loads it into a **PostgreSQL** data warehouse, ready for analysis.

* **Data Visualization:** Offers an intuitive user interface via Streamlit to visualize farm weather on a map, view detailed weather history, and analyze cereal gains for each farm.

---

## Getting Started

To launch all pipeline services (ingestion, database, and Streamlit interface), follow these simple steps:

1.  Ensure Docker and Docker Compose are installed on your machine.

2.  Navigate to the `deploy/` directory at the project root in your terminal:

    ```bash
    cd deploy/
    ```

3.  Launch all services using Docker Compose:

    ```bash
    docker compose up --build -d
    ```

    (The `--build` option rebuilds images if code changes have been made. The `-d` option runs containers in the background.)

Once the services are launched:

* The Streamlit application will generally be accessible via `http://localhost:8501`.

* The ingestion services (Agreste and OpenWeather) run in the background and are ready to be triggered.

---

## Automation (Cron)

The pipeline is configured so that ingestion and transformation processes are automatically launched at regular intervals using cron jobs. This ensures your data remains up-to-date without constant manual intervention.

---

## Manually Triggering Ingestion (Optional)

If you wish to manually trigger the ingestion processes for Agreste or OpenWeather data, you can use the following `curl` commands (default ports are 8080 for Agreste and 8081 for OpenWeather):

* **Ingest Agreste (Rica) data:**

    ```bash
    curl -X POST http://localhost:8080/ingest \
      -H "Content-Type: application/json" \
      -d '{
        "zipUrl": "[https://agreste.agriculture.gouv.fr/agreste-web/download/service/SV-Accès](https://agreste.agriculture.gouv.fr/agreste-web/download/service/SV-Accès) micro données RICA/RicaMicrodonnées2023_v2.zip",
        "csvFileName": "Rica_France_micro_Donnees_ex2023.csv"
      }'
    ```

* **Ingest OpenWeather data:**

    ```bash
    curl -X POST http://localhost:8081/ingest
    ```

---

## Accessing the PostgreSQL Database

To connect directly to the PostgreSQL database inside the Docker container, you can use the `docker exec` command:

1.  Identify the name or ID of your PostgreSQL container (you can find it with `docker ps`).

2.  Execute the `psql` command inside the container:

    ```bash
    docker exec -it <postgresql_container_name_or_id> psql -U postgres -d mydatabase
    ```

---