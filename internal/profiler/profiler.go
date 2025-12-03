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
	columnResult.NameFile = strings.TrimSuffix(fileName, ".csv")
	columnResult.TotalColumns = len(columns)
	if len(columns) == 0 {
		return
	}

	for i, col := range columns {
		columnResult.Columns = append(columnResult.Columns, AnalyzeColumn(col))
		fmt.Printf("------ %d Coluna Analisada\n",(i+1))
		lengthCol := len(col.Values)
		if lengthCol > columnResult.TotalMaxRows {
			columnResult.TotalMaxRows = lengthCol
		}
	}
	fmt.Println("===== Retorno do Profile")
	return
}
