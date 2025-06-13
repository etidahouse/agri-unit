import os
import streamlit as st
import pandas as pd
import psycopg2
from dotenv import load_dotenv
import folium
from streamlit_folium import st_folium
import plotly.express as px

load_dotenv()

DB_HOST = os.getenv("DB_HOST", "localhost")
DB_PORT = os.getenv("DB_PORT", "5432")
DB_USER = os.getenv("DB_USER", "postgres")
DB_PASSWORD = os.getenv("DB_PASSWORD", "password")
DB_NAME = os.getenv("DB_NAME", "mydatabase")

@st.cache_resource
def get_connection():
    return psycopg2.connect(
        host=DB_HOST,
        port=DB_PORT,
        user=DB_USER,
        password=DB_PASSWORD,
        dbname=DB_NAME
    )

@st.cache_data(ttl=300)
def load_units():
    conn = get_connection()
    df = pd.read_sql("SELECT id, id_num, latitude, longitude FROM agricultural_units", conn)
    return df

@st.cache_data(ttl=300)
def load_weather_history(unit_id):
    conn = get_connection()
    query = """
        SELECT created_at, temperature, humidity, clouds, weather_main
        FROM weather
        WHERE agricultural_unit_id = %s
        ORDER BY created_at
    """
    df = pd.read_sql(query, conn, params=(unit_id,))
    return df

@st.cache_data(ttl=300)
def load_latest_weather():
    conn = get_connection()
    query = """
        SELECT DISTINCT ON (agricultural_unit_id)
               agricultural_unit_id, temperature, humidity, clouds, weather_main, created_at
        FROM weather
        ORDER BY agricultural_unit_id, created_at DESC
    """
    df = pd.read_sql(query, conn)
    return df

@st.cache_data(ttl=300)
def load_units():
    conn = get_connection()
    query = """
        SELECT
            au.id,
            au.id_num,
            au.latitude,
            au.longitude,
            aus.data AS exploitation_data
        FROM
            agricultural_units au
        LEFT JOIN
            (SELECT DISTINCT ON (id_num) id_num, data, year FROM agricultural_unit_surveys ORDER BY id_num, year DESC) aus
        ON
            au.id_num = aus.id_num;
    """
    df = pd.read_sql(query, conn)
    return df

exploitation_variable_labels = {
    "PBV3COLZ": "Produit Brut : colza (€)",
    "PBV3BLED": "Produit Brut : blé dur (€)",
    "PBV3BLET": "Produit Brut : blé tendre et épeautre (€)",
}

st.title("🌤️ Analyse météo des exploitations agricoles")

tabs = st.tabs(["🗺️ Carte météo", "📈 Historique météo", "📊 Analyse Exploitation"])

with tabs[0]:
    st.header("Carte des exploitations avec météo actuelle")

    units = load_units()
    weather = load_latest_weather()

    merged = units.merge(weather, left_on="id", right_on="agricultural_unit_id", how="left")

    map_center = [merged["latitude"].mean(), merged["longitude"].mean()]
    m = folium.Map(location=map_center, zoom_start=6)

    for _, row in merged.iterrows():
        popup = f"""
        <b>ID:</b> {row['id_num']}<br>
        <b>Température:</b> {row['temperature']}°C<br>
        <b>Humidité:</b> {row['humidity']}%<br>
        <b>Nuages:</b> {row['clouds']}%<br>
        <b>Condition:</b> {row['weather_main']}
        """
        folium.Marker(
            location=[row["latitude"], row["longitude"]],
            popup=popup,
            icon=folium.Icon(color="blue", icon="cloud")
        ).add_to(m)

    st_folium(m, width=700, height=500)

with tabs[1]:
    st.header("Évolution météo d'une exploitation")

    units = load_units()
    selected_id = st.selectbox("Choisir une exploitation", units["id_num"])

    selected_uuid = units[units["id_num"] == selected_id]["id"].values[0]
    df_hist = load_weather_history(selected_uuid)

    if df_hist.empty:
        st.warning("Pas de données météo pour cette exploitation.")
    else:
        st.subheader("📊 Température")
        st.plotly_chart(px.line(df_hist, x="created_at", y="temperature", title="Température dans le temps"))

        st.subheader("🌫️ Humidité")
        st.plotly_chart(px.line(df_hist, x="created_at", y="humidity", title="Humidité dans le temps"))

        st.subheader("☁️ Nuages")
        st.plotly_chart(px.line(df_hist, x="created_at", y="clouds", title="Couverture nuageuse"))

        st.subheader("🌦️ Conditions météo")
        st.dataframe(df_hist[["created_at", "weather_main"]])

with tabs[2]:
    st.header("📊 Analyse des Gains en Céréales par Exploitation")
    st.write("Sélectionnez une exploitation pour visualiser ses **Produits Bruts pour le Colza, le Blé Dur et le Blé Tendre/Épeautre**.")

    units = load_units()
    if not units.empty:
        selected_unit_id_num = st.selectbox(
            "Choisir une exploitation",
            units["id_num"],
            key='cereals_analysis_select'
        )
    
        selected_unit_data = units[units["id_num"] == selected_unit_id_num].iloc[0]
        exploitation_data_jsonb = selected_unit_data.get('exploitation_data')

        if exploitation_data_jsonb is None:
            st.warning(f"Pas de données d'enquête (jsonb) disponibles pour l'exploitation {selected_unit_id_num}.")
        else:
            colza_value = exploitation_data_jsonb.get("PBV3COLZ", 0)
            bledur_value = exploitation_data_jsonb.get("PBV3BLED", 0)
            bletendre_value = exploitation_data_jsonb.get("PBV3BLET", 0)

            colza_value = colza_value if isinstance(colza_value, (int, float)) else 0
            bledur_value = bledur_value if isinstance(bledur_value, (int, float)) else 0
            bletendre_value = bletendre_value if isinstance(bletendre_value, (int, float)) else 0

            cereal_data_for_display = {
                exploitation_variable_labels["PBV3COLZ"]: colza_value,
                exploitation_variable_labels["PBV3BLED"]: bledur_value,
                exploitation_variable_labels["PBV3BLET"]: bletendre_value,
            }

            st.write(f"**Produits Bruts pour l'exploitation {selected_unit_id_num}**")
            st.dataframe(pd.DataFrame(list(cereal_data_for_display.items()), columns=["Céréale", "Produit Brut (€)"]).set_index("Céréale"))

            fig_selected_cereals = px.bar(
                pd.DataFrame(list(cereal_data_for_display.items()), columns=["Céréale", "Produit Brut (€)"]),
                x="Céréale",
                y="Produit Brut (€)",
                title=f"Produits Bruts de l'exploitation {selected_unit_id_num}",
                labels={"Produit Brut (€)": "Valeur en Euros (€)"},
                color="Céréale"
            )
            fig_selected_cereals.update_layout(xaxis_title="Type de Céréale", yaxis_title="Produit Brut en Euros (€)")
            st.plotly_chart(fig_selected_cereals)

    else:
        st.info("Aucune exploitation disponible pour l'analyse des céréales.")
