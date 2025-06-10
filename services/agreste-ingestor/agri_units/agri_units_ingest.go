package agri_units

import (
	"agreste-ingestor/misc"
	"fmt"
	"strconv"
)

func HandleAgriUnitSurveyIngest(
	endpointURL string,
	fileName string,
	agriUnitStorage AgriUnitStorage,
	agriUnitSurveyStorage AgriculturalUnitSurveyStorage,
) error {
	fmt.Printf("Starting ingestion process for endpoint: %s\n", endpointURL)

	existingAgriUnits, err := agriUnitStorage.SelectAll()
	if err != nil {
		return fmt.Errorf("failed to select all agricultural units: %w", err)
	}
	fmt.Printf("Fetched %d existing agricultural units.\n", len(existingAgriUnits))

	existingAgriUnitSurveys, err := agriUnitSurveyStorage.SelectAll()
	if err != nil {
		return fmt.Errorf("failed to select all agricultural unit surveys: %w", err)
	}
	fmt.Printf("Fetched %d existing agricultural unit surveys.\n", len(existingAgriUnitSurveys))

	agriUnitSurveyCSVData, err := misc.DownloadZipAndReadSpecificCSV(endpointURL, fileName)
	if err != nil {
		return fmt.Errorf("failed to fetch csv survey: %w", err)
	}
	fmt.Printf("survey fetched")

	agriUnitsByIDNum := make(map[int]AgriculturalUnit)
	for _, unit := range existingAgriUnits {
		agriUnitsByIDNum[unit.IDNum] = unit
	}

	agriUnitSurveysByIDNumAndYear := make(map[string]AgriculturalUnitSurvey)
	for _, survey := range existingAgriUnitSurveys {
		compositeKey := fmt.Sprintf("%d_%d", survey.IDNum, survey.Year)
		agriUnitSurveysByIDNumAndYear[compositeKey] = survey
	}

	headerRow := agriUnitSurveyCSVData[0]
	headerMap := make(map[string]int)
	for idx, colName := range headerRow {
		headerMap[colName] = idx
	}

	const (
		HeaderIDNum = "IDNUM"
		HeaderYear  = "MILEX"
	)

	idNumColIdx, idNumExists := headerMap[HeaderIDNum]
	if !idNumExists {
		return fmt.Errorf("CSV header missing required column: %s", HeaderIDNum)
	}
	yearColIdx, yearExists := headerMap[HeaderYear]
	if !yearExists {
		return fmt.Errorf("CSV header missing required column: %s", HeaderYear)
	}

	var newAgriUnits []AgriculturalUnit
	var newAgriUnitSurveys []AgriculturalUnitSurvey
	seenIDNumsInCsv := make(map[int]struct{})

	for i, row := range agriUnitSurveyCSVData {
		if i == 0 {
			continue
		}

		idNumStr := row[idNumColIdx]
		idNum, err := strconv.Atoi(idNumStr)
		if err != nil {
			fmt.Printf("Skipping row %d: Invalid %s '%s' - %v\n", i, HeaderIDNum, idNumStr, err)
			continue
		}
		yearStr := row[yearColIdx]
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			fmt.Printf("Skipping row %d: Invalid %s '%s' - %v\n", i, HeaderYear, yearStr, err)
			continue
		}

		if _, seen := seenIDNumsInCsv[idNum]; !seen {
			seenIDNumsInCsv[idNum] = struct{}{}
			if _, exists := agriUnitsByIDNum[idNum]; !exists {
				lat, lon := misc.GenerateRandomCoordinates(0, 0, 0, 0)
				newUnit := CreateAgriculturalUnit(AgriculturalUnitValue{
					IDNum:     idNum,
					Latitude:  lat,
					Longitude: lon,
				})
				newAgriUnits = append(newAgriUnits, newUnit)
			}
		}

		compositeKey := fmt.Sprintf("%d_%d", idNum, year)

		if _, exists := agriUnitSurveysByIDNumAndYear[compositeKey]; !exists {
			surveyDataPayload := make(map[string]interface{})
			for colName, colIdx := range headerMap {
				if colIdx < len(row) {
					valStr := row[colIdx]
					var parsedValue interface{} = valStr
					if f, err := strconv.ParseFloat(valStr, 64); err == nil {
						parsedValue = f
					} else if b, err := strconv.ParseBool(valStr); err == nil {
						parsedValue = b
					}
					surveyDataPayload[colName] = parsedValue
				}
			}

			if otfDDVal, ok := surveyDataPayload["OTEFDD"]; ok {
				if otfDDNum, isFloat := otfDDVal.(float64); !isFloat || otfDDNum != 1500 {
					continue
				}
			} else {
				continue
			}

			newSurvey := CreateAgriculturalUnitSurvey(AgriculturalUnitSurveyValue{
				IDNum: idNum,
				Year:  year,
				Data:  surveyDataPayload,
			})
			newAgriUnitSurveys = append(newAgriUnitSurveys, newSurvey)
		}
	}

	fmt.Printf("Identified %d new Agricultural Units to potentially create.\n", len(newAgriUnits))
	fmt.Printf("Identified %d new Agricultural Unit Surveys to insert.\n", len(newAgriUnitSurveys))

	for _, unit := range newAgriUnits {
		err := agriUnitStorage.InsertOrUpdate(unit)
		if err != nil {
			fmt.Printf("Error inserting/updating AgriculturalUnit %d: %v\n", unit.IDNum, err)
			return err
		}
	}
	for _, survey := range newAgriUnitSurveys {
		err := agriUnitSurveyStorage.InsertOrUpdate(survey)
		if err != nil {
			fmt.Printf("Error inserting/updating AgriculturalUnitSurvey %d_%d: %v\n", survey.IDNum, survey.Year, err)
			return err
		}
	}

	return nil
}
