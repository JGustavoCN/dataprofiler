package profiler

import (
	"strconv"
	"strings"
)

type Column struct {
	Name   string
	Values []string
}

type ColumnResult struct {
	Name        string
	MainType    string
	BlankCount  int
	CountFilled int
	Filled      float64
	BlankRatio  float64
	TypeCounts  map[string]int
	Stats       map[string]string
}

func AnalyzeColumn(column Column) (result ColumnResult) {
	if len(column.Values) == 0 {
		return
	}

	typeCounts := make(map[string]int)
	var numericValues []float64
	filledCount := 0
	blankCount := 0
	for _, v := range column.Values {
		trimmed := strings.TrimSpace(v)
		if trimmed == "" {
			blankCount++
			continue
		}

		inferredType := InferType(trimmed)
		typeCounts[inferredType]++
		if inferredType == "float" || inferredType == "int" {
			number, _ := strconv.ParseFloat(trimmed, 64)
			numericValues = append(numericValues, number)
		}

		filledCount++

	}
	counts := 0
	for k, v := range typeCounts {
		if v > counts {
			result.MainType = k
			counts = v
		}
	}

	result.Name = column.Name
	result.TypeCounts = typeCounts
	result.BlankCount = blankCount
	result.CountFilled = filledCount
	if result.MainType == "float" || result.MainType == "int" {
		result.Stats = StatsCalc(numericValues)
	}
	result.Filled = (float64(filledCount) / float64(len(column.Values)))
	result.BlankRatio = (float64(blankCount) / float64(len(column.Values)))
	return
}
