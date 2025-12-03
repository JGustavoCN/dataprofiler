package infra

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"

	"github.com/JGustavoCN/dataprofiler/internal/profiler"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func LoadCSV(filePath string) ([]profiler.Column, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	nameFile := filepath.Base(file.Name())
	column, err := ParseData(file)
	return column, nameFile, err
}

func ParseData(file io.Reader) ([]profiler.Column, error) {
	deco := charmap.Windows1252.NewDecoder()
	file = transform.NewReader(file, deco)
	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.LazyQuotes = true
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return []profiler.Column{}, nil
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

	return columns, nil
}
