package infra

import (
	"bufio"
	"encoding/csv"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/JGustavoCN/dataprofiler/internal/profiler"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func LoadCSV(logger *slog.Logger, filePath string) ([]profiler.Column, string, error) {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(io.Discard, nil))
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	nameFile := filepath.Base(file.Name())
	column, err := ParseData(logger, file)
	return column, nameFile, err
}

func ParseData(logger *slog.Logger, file io.Reader) ([]profiler.Column, error) {
	if logger == nil { 
		logger = slog.New(slog.NewJSONHandler(io.Discard, nil)) 
	}
	smartReader, err := NewSmartReader(logger, file)
	if err != nil {
		return nil, err
	}

	bufferedSmartReader := bufio.NewReader(smartReader)
	
	separator, err := DetectSeparator(bufferedSmartReader)
	if err != nil {
		separator = ';'
		logger.Warn("Falha na detecção de separador, usando fallback", "error", err, "fallback", separator)
	} else {
		logger.Info("Separador detectado", "separator", string(separator))
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

	logger.Info("Estrutura carregada", 
		"total_rows", len(records),
		"columns_count", len(headers),
		"headers", headers,
	)
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

func ParseDataAsync(logger *slog.Logger, r io.Reader) ([]string, <-chan []string, error) {
	if logger == nil { 
		logger = slog.New(slog.NewJSONHandler(io.Discard, nil)) 
	}
	out := make(chan []string, 100)

	smartReader, err := NewSmartReader(logger, r)
	if err != nil {
		close(out)
		return nil, nil, err
	}

	bufferedSmartReader := bufio.NewReader(smartReader)
	separator, err := DetectSeparator(bufferedSmartReader)
	if err != nil {
		separator = ';'
		logger.Warn("Falha na detecção de separador, usando fallback", "error", err, "fallback", separator)
	} else {
		logger.Info("Separador detectado", "separator", string(separator))
	}

	reader := csv.NewReader(bufferedSmartReader)
	reader.Comma = separator
	reader.LazyQuotes = true
	reader.ReuseRecord = true
	
	headers, err := reader.Read()
	if err != nil {
		close(out)           
		return nil, nil, err 
	}
	
	logger.Info("Início do streaming", 
		"columns_count", len(headers), 
		"headers", headers,
	)

	go func() {
		defer close(out)
		count := 0
		errorCount := 0
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				logger.Warn("Erro ao ler linha CSV", 
						"line_attempt", count + errorCount,
						"error", err,
				)
				errorCount++
				continue
			}

			rowCopy := profiler.GetRowSlice()
			rowCopy = append(rowCopy, record...)

			out <- rowCopy
			count++
		}
		logger.Info("Streaming finalizado", 
			"total_rows_read", count-1,
			"total_errors", errorCount,
		)
	}()

	return headers, out, nil
}


func NewSmartReader(logger *slog.Logger, r io.Reader) (io.Reader, error) {
	br := bufio.NewReader(r)
	bytesTosample , err := br.Peek(1024)
	if err != nil && err != io.EOF {
		return nil, err
	}

	if utf8.Valid(bytesTosample) {
		logger.Debug("Encoding detectado: UTF-8")
		return br, nil
	}	
	logger.Warn("Encoding UTF-8 inválido detectado. Convertendo Windows1252 -> UTF-8")
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