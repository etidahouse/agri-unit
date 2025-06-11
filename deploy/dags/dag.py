from airflow import DAG
from airflow.operators.bash import BashOperator
from datetime import datetime, timedelta

default_args = {
    'start_date': datetime(2025, 6, 11),
    'retries': 1,
    'retry_delay': timedelta(minutes=1),
}

with DAG('agreste_ingest', schedule_interval='* * * * *', default_args=default_args, catchup=False) as dag_weather:
    curl_weather = BashOperator(
        task_id='curl_weather',
        bash_command="""
        curl -X POST http://weather-ingestor:8080/ingest
        """
    )

with DAG('weather_ingest', schedule_interval='*/5 * * * *', default_args=default_args, catchup=False) as dag_other:
    curl_other = BashOperator(
        task_id='curl_other',
         bash_command="""
        curl -X POST http://agreste-ingestor:8080/ingest \
             -H "Content-Type: application/json" \
             -d '{"zipUrl": "https://agreste.agriculture.gouv.fr/agreste-web/download/service/SV-Accès micro données RICA/RicaMicrodonnées2023_v2.zip", "csvFileName": "Rica_France_micro_Donnees_ex2023.csv"}'
        """
    )
