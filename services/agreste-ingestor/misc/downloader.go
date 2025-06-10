package misc

import (
	"archive/zip"
	"crypto/tls"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func DownloadZipAndReadSpecificCSV(zipURL, csvFileNameInZip string) ([][]string, error) {
	fmt.Printf("Starting ingestion process for endpoint: %s\n", zipURL)
	tempZipFile, err := os.CreateTemp("", "downloaded-*.zip")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary zip file: %w", err)
	}
	defer os.Remove(tempZipFile.Name())
	defer tempZipFile.Close()

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := httpClient.Get(zipURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download ZIP file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download ZIP, status: %d %s", resp.StatusCode, resp.Status)
	}

	_, err = io.Copy(tempZipFile, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to write ZIP content to temp file: %w", err)
	}

	if err := tempZipFile.Sync(); err != nil {
		return nil, fmt.Errorf("failed to sync temp zip file: %w", err)
	}

	zipReader, err := zip.OpenReader(tempZipFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to open downloaded ZIP file: %w", err)
	}
	defer zipReader.Close()

	var csvFile *zip.File
	for _, file := range zipReader.File {
		if strings.EqualFold(filepath.Base(file.Name), csvFileNameInZip) {
			csvFile = file
			break
		}
	}

	if csvFile == nil {
		return nil, fmt.Errorf("CSV file '%s' not found in the ZIP archive", csvFileNameInZip)
	}

	rc, err := csvFile.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file '%s' from ZIP: %w", csvFile.Name, err)
	}
	defer rc.Close()

	reader := csv.NewReader(rc)

	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV records from '%s': %w", csvFileNameInZip, err)
	}

	return records, nil
}
