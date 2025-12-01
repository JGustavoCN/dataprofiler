package infra

import (
	"encoding/csv"
	"os"
	"path/filepath"

	"github.com/JGustavoCN/dataprofiler/internal/profiler"
)

func LoadCSV(filePath string) ([]profiler.Column, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", err 
	}
	defer file.Close()

	nameFile := filepath.Base(file.Name())

	reader := csv.NewReader(file)
	reader.Comma = ';'       
	reader.LazyQuotes = true 

	records, err := reader.ReadAll()
	if err != nil {
		return nil, "", err
	}

	if len(records) == 0 {
		return []profiler.Column{}, nameFile, nil
	}

	headers := records[0]
	columns := make([]profiler.Column, len(headers))
	for i, name := range headers {
		columns[i] = profiler.Column{
			Name:   name,
			Values: make([]string, 0, len(records)-1), 
		}
	}

	for _, row := range records[1:] {
		for i, value := range row {
			if i < len(columns) {
				columns[i].Values = append(columns[i].Values, value)
			}
		}
	}

	return columns, nameFile, nil
}
