package infra

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"unicode/utf8"

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
	smartReader, err := NewSmartReader(file)
	if err != nil {
		return nil, err
	}

	bufferedSmartReader := bufio.NewReader(smartReader)
	
	separator, err := DetectSeparator(bufferedSmartReader)
	if err != nil {
		separator = ';' 
	}

	reader := csv.NewReader(bufferedSmartReader)
	reader.Comma = separator
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

func ParseDataAsync(r io.Reader) ([]string, <-chan []string, error) {
	out := make(chan []string, 100)

	smartReader, err := NewSmartReader(r)
	if err != nil {
		close(out)
		return nil, nil, err
	}

	bufferedSmartReader := bufio.NewReader(smartReader)
	separator, err := DetectSeparator(bufferedSmartReader)
	if err != nil {
		separator = ';' 
	}

	reader := csv.NewReader(bufferedSmartReader)
	reader.Comma = separator
	reader.LazyQuotes = true

	
	headers, err := reader.Read()
	if err != nil {
		close(out)           
		return nil, nil, err 
	}
	fmt.Println("✅ Header lido com sucesso!")
	go func() {
		defer close(out)
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				continue
			}
			out <- record
		}
	}()

	return headers, out, nil
}


func NewSmartReader(r io.Reader) (io.Reader, error) {
	br := bufio.NewReader(r)
	bytesTosample , err := br.Peek(1024)
	if err != nil && err != io.EOF {
		return nil, err
	}

	if utf8.Valid(bytesTosample) {
		return br, nil
	}	

	decoderReader := transform.NewReader(br, charmap.Windows1252.NewDecoder())
	return decoderReader, nil
}



func DetectSeparator(r *bufio.Reader) (rune, error) {
	
	bytesToPeek, err := r.Peek(2048) 
	if err != nil && err != io.EOF {
		return ';', err 
	}

	semicolonCount := 0 // ;
	commaCount := 0     // ,
	pipeCount := 0      // |
	tabCount := 0       // \t

	for _, b := range bytesToPeek {
		if b == '\n' || b == '\r' {
			break
		}
		switch b {
		case ';':
			semicolonCount++
		case ',':
			commaCount++
		case '|':
			pipeCount++
		case '\t':
			tabCount++
		}
	}

	separator := ';'
	maxCount := semicolonCount

	if commaCount > maxCount {
		maxCount = commaCount
		separator = ','
	}
	if pipeCount > maxCount {
		maxCount = pipeCount
		separator = '|'
	}
	if tabCount > maxCount {
		separator = '\t'
	}

	return separator, nil
}