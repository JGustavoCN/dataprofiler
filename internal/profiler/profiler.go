	package profiler

	import "strings"

	type ProfilerResult struct {
		NameFile     string
		TotalMaxRows    int
		TotalColumns int
		Columns      []ColumnResult
	}

	func Profile(columns []Column, fileName string) (columnResult ProfilerResult) {
		columnResult.NameFile = strings.TrimSuffix(fileName,".csv")
		columnResult.TotalColumns = len(columns)
		if len(columns) == 0 {
			return
		}
		
		for _,col := range columns{
			columnResult.Columns = append(columnResult.Columns, AnalyzeColumn(col))
			lengthCol := len(col.Values)
			if lengthCol > columnResult.TotalMaxRows {
				columnResult.TotalMaxRows = lengthCol
			}
		}
		return
	}