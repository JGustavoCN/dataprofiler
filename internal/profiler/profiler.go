package profiler

import (
	"fmt"
	"strings"
)

type ProfilerResult struct {
	NameFile     string
	TotalMaxRows int
	TotalColumns int
	Columns      []ColumnResult
}

func Profile(columns []Column, fileName string) (columnResult ProfilerResult) {
	setResultMetadata(columns, &columnResult, fileName)

	if len(columns) == 0 {
		return
	}

	for i, col := range columns {
		columnResult.Columns = append(columnResult.Columns, AnalyzeColumn(col))
		fmt.Printf("------ %d Coluna Analisada\n", (i + 1))
		columnCount := len(col.Values)
		if columnCount > columnResult.TotalMaxRows {
			columnResult.TotalMaxRows = columnCount
		}
	}
	fmt.Println("===== Retorno do Profile")
	return
}

func ProfileAsync(headers []string, dataChan <-chan []string, fileName string) (profilerResult ProfilerResult) {
	setResultMetadata(headers, &profilerResult, fileName)

	accumulators := make([]*ColumnAccumulator,profilerResult.TotalColumns)

	for i, name := range headers {
		accumulators[i] = NewColumnAccumulator(name)
	}

	rowCount := 0
	for record := range dataChan {
		rowCount++
		for i, value := range record {
			if i < len(accumulators) {
				accumulators[i].Add(value)
			}
		}
	}

	columnResults := make([]ColumnResult, len(headers))
	for i, acc := range accumulators{
		columnResults[i] = acc.Result()
		fmt.Printf("---- ✅ Coluna %d lida com sucesso!\n", i+1)
	}

	fmt.Printf("✅ Processamento Async concluído: %d linhas processadas.\n", rowCount)
	profilerResult.TotalMaxRows = rowCount
	profilerResult.Columns = columnResults
	return 
}

func setResultMetadata[T string | Column](columns []T, profilerResult *ProfilerResult, fileName string) {
	profilerResult.NameFile = strings.TrimSuffix(fileName, ".csv")
	profilerResult.TotalColumns = len(columns)
}
