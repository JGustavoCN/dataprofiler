package profiler

import (
	"io"
	"log/slog"
	"math/rand/v2"
	"strings"
	"time"
)

type StreamData struct {
	Row        []string
	LineNumber int
	Err        error
}
type DirtyLine struct {
	Line   int    `json:"line"`
	Reason string `json:"reason"`
}

type ProfilerResult struct {
	NameFile        string         `json:"name_file"`
	TotalMaxRows    int            `json:"total_max_rows"`
	TotalColumns    int            `json:"total_columns"`
	DirtyLinesCount int            `json:"dirty_lines_count"`
	Columns         []ColumnResult `json:"columns"`
	SampleRows      [][]string     `json:"sample_rows"`
	DirtyLines      []DirtyLine    `json:"dirty_lines"`
}

func Profile(logger *slog.Logger, columns []Column, fileName string) (columnResult ProfilerResult) {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(io.Discard, nil))
	}
	setResultMetadata(columns, &columnResult, fileName)

	if len(columns) == 0 {
		logger.Warn("Profile chamado com colunas vazias", "filename", fileName)
		return
	}

	logger.Info("Iniciando análise estatística (Síncrona)",
		"total_columns", len(columns),
		"filename", fileName,
	)

	for i, col := range columns {
		columnResult.Columns = append(columnResult.Columns, AnalyzeColumn(col))
		logger.Debug("Coluna analisada",
			"index", i+1,
			"column_name", col.Name,
			"rows", len(col.Values),
		)
		columnCount := len(col.Values)
		if columnCount > columnResult.TotalMaxRows {
			columnResult.TotalMaxRows = columnCount
		}
	}
	logger.Info("Análise estatística concluída", "total_columns_analyzed", len(columnResult.Columns))
	return
}

func ProfileAsync(logger *slog.Logger, headers []string, dataChan <-chan StreamData, fileName string) (profilerResult ProfilerResult) {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(io.Discard, nil))
	}
	setResultMetadata(headers, &profilerResult, fileName)

	accumulators := make([]*ColumnAccumulator, profilerResult.TotalColumns)
	for i, name := range headers {
		accumulators[i] = NewColumnAccumulator(name)
	}

	const previewSize = 50
	sampleRows := make([][]string, 0, previewSize)
	seed := uint64(time.Now().UnixNano())
	rng := rand.New(rand.NewPCG(seed, seed+1))

	rowCount := 0
	dirtyLines := []DirtyLine{}
	for msg := range dataChan {

		if msg.Err != nil {
			if len(dirtyLines) < 1000 {
				dirtyLines = append(dirtyLines, DirtyLine{
					Line:   msg.LineNumber,
					Reason: msg.Err.Error(),
				})
			}
			continue
		}
		record := msg.Row
		rowCount++

		for i, value := range record {
			if i < len(accumulators) {
				accumulators[i].Add(value)
			}
		}

		if len(sampleRows) < previewSize {
			rowCopy := make([]string, len(record))
			copy(rowCopy, record)
			sampleRows = append(sampleRows, rowCopy)
		} else {
			k := rng.IntN(rowCount)
			if k < previewSize {
				rowCopy := make([]string, len(record))
				copy(rowCopy, record)
				sampleRows[k] = rowCopy
			}
		}

		PutRowSlice(record)

		if rowCount%200000 == 0 {
			logger.Info("Processamento em andamento", "rows_processed", rowCount)
		}
	}

	columnResults := make([]ColumnResult, len(headers))
	for i, acc := range accumulators {
		columnResults[i] = acc.Result()
		logger.Debug("Coluna finalizada", "index", i+1, "column", headers[i])
	}

	logger.Info("Processamento Async concluído",
		"total_rows", rowCount,
		"total_columns", len(headers),
		"filename", fileName,
		"dirty_lines", len(dirtyLines),
	)
	profilerResult.DirtyLines = dirtyLines
	profilerResult.DirtyLinesCount = len(dirtyLines)
	profilerResult.TotalMaxRows = rowCount
	profilerResult.Columns = columnResults
	profilerResult.SampleRows = sampleRows
	return
}

func setResultMetadata[T string | Column](columns []T, profilerResult *ProfilerResult, fileName string) {
	profilerResult.NameFile = strings.TrimSuffix(fileName, ".csv")
	profilerResult.TotalColumns = len(columns)
}
