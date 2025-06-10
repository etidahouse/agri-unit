package misc

import (
	"testing"
)

func TestDownloadZipAndReadSpecificCSV_AgresteData(t *testing.T) {
	agresteZipURL := "https://agreste.agriculture.gouv.fr/agreste-web/download/service/SV-Accès micro données RICA/RicaMicrodonnées2023_v2.zip"
	agresteCsvFileName := "Rica_France_micro_Donnees_ex2023.csv"

	csvData, err := DownloadZipAndReadSpecificCSV(agresteZipURL, agresteCsvFileName)

	if err != nil {
		t.Fatalf("DownloadZipAndReadSpecificCSV a retourné une erreur inattendue : %v", err)
	}

	if len(csvData) == 0 {
		t.Error("DownloadZipAndReadSpecificCSV n'a retourné aucune donnée (0 lignes), des enregistrements étaient attendus.")
	}

	t.Logf("Lecture réussie de %d lignes depuis '%s'. Première ligne : %v", len(csvData), agresteCsvFileName, csvData[0])
}
