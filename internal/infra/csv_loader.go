package infra

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"unicode"
	"unicode/utf8"

	"github.com/JGustavoCN/dataprofiler/internal/profiler"
	"golang.org/x/text/encoding/charmap"
	unicodeenc "golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func LoadCSV(logger *slog.Logger, filePath string) ([]profiler.Column, string, error) {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(io.Discard, nil))
	}
	file, err := os.Open(filePath)
	if err != nil {
		logger.Error("Falha ao abrir arquivo", "path", filePath, "error", err)
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

func ParseDataAsync(ctx context.Context, logger *slog.Logger, r io.Reader) ([]string, <-chan profiler.StreamData, error) {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(io.Discard, nil))
	}

	smartReader, err := NewSmartReader(logger, r)
	if err != nil {
		return nil, nil, err
	}

	bufferedSmartReader := bufio.NewReader(smartReader)
	isJson, err := sniffJSON(bufferedSmartReader)
	if err != nil {
		return nil, nil, fmt.Errorf("erro ao detectar formato: %w", err)
	}

	if isJson {
		logger.Info("Formato detectado: JSONL (Logs/NoSQL)")
		return parseJSONLAsync(ctx, logger, bufferedSmartReader)
	}
	logger.Info("Formato detectado: CSV (Tabular)")
	return parseCSVAsync(ctx, logger, bufferedSmartReader)

}

func sniffJSON(r *bufio.Reader) (bool, error) {
	bytesToPeek, err := r.Peek(50)
	if err != nil && err != io.EOF {
		return false, err
	}

	for _, b := range bytesToPeek {
		if unicode.IsSpace(rune(b)) {
			continue
		}
		if b == '{' {
			return true, nil
		}
		return false, nil
	}
	return false, nil
}

func parseCSVAsync(ctx context.Context, logger *slog.Logger, reader *bufio.Reader) ([]string, <-chan profiler.StreamData, error) {
	out := make(chan profiler.StreamData, 100)

	separator, err := DetectSeparator(reader)
	if err != nil {
		separator = ';'
		logger.Warn("Falha na detecção de separador, usando fallback", "error", err, "fallback", separator)
	} else {
		logger.Info("Separador detectado", "separator", string(separator))
	}

	csvReader := csv.NewReader(reader)
	csvReader.Comma = separator
	csvReader.LazyQuotes = true
	csvReader.ReuseRecord = true

	headersRef, err := csvReader.Read()
	if err != nil {
		close(out)
		return nil, nil, err
	}
	headers := make([]string, len(headersRef))
	copy(headers, headersRef)
	csvReader.FieldsPerRecord = len(headers)
	logger.Info("Início do streaming",
		"columns_count", len(headers),
		"headers", headers,
	)

	go func() {
		defer close(out)
		count := 0
		lineNum := 1
		errorCount := 0
		for {
			select {
			case <-ctx.Done():
				logger.Warn("Leitura cancelada pelo contexto")
				return
			default:
				record, err := csvReader.Read()
				lineNum++
				if err == io.EOF {
					goto EndProcessing
				}
				if err != nil {
					errorCount++
					out <- profiler.StreamData{
						Row:        nil,
						LineNumber: lineNum,
						Err:        err,
					}
					continue
				}

				rowCopy := profiler.GetRowSlice()
				rowCopy = append(rowCopy, record...)

				out <- profiler.StreamData{
					Row:        rowCopy,
					LineNumber: lineNum,
					Err:        nil,
				}
				count++

			}
		}
	EndProcessing:
		logger.Info("Streaming CSV finalizado",
			"total_rows_read", lineNum-1,
			"total_errors", errorCount,
		)
	}()

	return headers, out, nil
}

func parseJSONLAsync(ctx context.Context, logger *slog.Logger, reader *bufio.Reader) ([]string, <-chan profiler.StreamData, error) {
	out := make(chan profiler.StreamData, 100)

	scanner := bufio.NewScanner(reader)

	const maxCapacity = 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, nil, fmt.Errorf("erro lendo primeira linha JSON: %w", err)
		}
		return nil, nil, errors.New("arquivo JSONL vazio")
	}

	firstLine := scanner.Bytes()
	var firstMap map[string]interface{}
	if err := json.Unmarshal(firstLine, &firstMap); err != nil {
		return nil, nil, fmt.Errorf("erro de parsing na primeira linha (não é JSON válido?): %w", err)
	}

	headers := make([]string, 0, len(firstMap))
	for k := range firstMap {
		headers = append(headers, k)
	}

	sort.Strings(headers)

	logger.Info("Schema JSONL inferido", "headers", headers)

	go func() {
		defer close(out)

		processMap := func(m map[string]interface{}, lineNum int) {
			row := profiler.GetRowSlice()

			for _, header := range headers {
				val, exists := m[header]
				if !exists || val == nil {
					row = append(row, "")
				} else {
					row = append(row, fmt.Sprintf("%v", val))
				}
			}

			out <- profiler.StreamData{
				Row:        row,
				LineNumber: lineNum,
				Err:        nil,
			}
		}

		processMap(firstMap, 1)

		lineNum := 1
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
			}
			lineNum++

			if len(scanner.Bytes()) == 0 {
				continue
			}

			var currentMap map[string]interface{}
			if err := json.Unmarshal(scanner.Bytes(), &currentMap); err != nil {
				out <- profiler.StreamData{
					LineNumber: lineNum,
					Err:        fmt.Errorf("json malformado: %w", err),
				}
				continue
			}

			processMap(currentMap, lineNum)
		}

		if err := scanner.Err(); err != nil {
			logger.Error("Erro fatal no scanner JSONL", "error", err)
			out <- profiler.StreamData{
				LineNumber: lineNum,
				Err:        fmt.Errorf("erro de I/O: %w", err),
			}
		}
	}()

	return headers, out, nil
}

func NewSmartReader(logger *slog.Logger, r io.Reader) (io.Reader, error) {
	br := bufio.NewReader(r)

	bomCheck, err := br.Peek(4)
	if err != nil && err != io.EOF && len(bomCheck) < 2 {
		return br, nil
	}

	if len(bomCheck) >= 2 && bomCheck[0] == 0xFF && bomCheck[1] == 0xFE {
		logger.Info("Encoding detectado: UTF-16 LE (Convertendo para UTF-8)")
		win16le := unicodeenc.UTF16(unicodeenc.LittleEndian, unicodeenc.UseBOM)
		return transform.NewReader(br, win16le.NewDecoder()), nil
	}
	if len(bomCheck) >= 2 && bomCheck[0] == 0xFE && bomCheck[1] == 0xFF {
		logger.Info("Encoding detectado: UTF-16 BE (Convertendo para UTF-8)")
		win16be := unicodeenc.UTF16(unicodeenc.BigEndian, unicodeenc.UseBOM)
		return transform.NewReader(br, win16be.NewDecoder()), nil
	}

	const sampleSize = 2048
	sample, err := br.Peek(sampleSize)

	if err != nil && err != io.EOF {
		return nil, err
	}

	if utf8.Valid(sample) {
		logger.Debug("Encoding detectado: UTF-8 (Nativo)")
		return br, nil
	}

	logger.Warn("Encoding UTF-8 inválido detectado na amostra. Aplicando fallback Windows1252 -> UTF-8")
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
