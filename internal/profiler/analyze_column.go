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
	MainType    DataType
	BlankCount  int
	CountFilled int
	Filled      float64
	BlankRatio  float64
	TypeCounts  map[DataType]int
	Stats       map[string]string
}

func AnalyzeColumn(column Column) (result ColumnResult) {

	if len(column.Values) == 0 {
		return ColumnResult{Name: column.Name, MainType: TypeEmpty}
	}

	result.TypeCounts = make(map[DataType]int)
	var numericValues []float64

	filledCount := 0
	blankCount := 0

	for _, v := range column.Values {
		trimmed := strings.TrimSpace(v)

		if trimmed == "" {
			blankCount++
			continue
		}

		inferredType := InferType(trimmed, column.Name)
		result.TypeCounts[inferredType]++
		filledCount++

		if inferredType == TypeInteger || inferredType == TypeFloat {

			valClean := strings.Replace(trimmed, ",", ".", 1)

			if number, err := strconv.ParseFloat(valClean, 64); err == nil {
				numericValues = append(numericValues, number)
			}
		}
	}

	result.MainType = determineMainType(result.TypeCounts)

	if result.MainType == TypeInteger || result.MainType == TypeFloat {

		result.Stats = StatsCalc(numericValues)
	}

	result.Name = column.Name
	result.BlankCount = blankCount
	result.CountFilled = filledCount

	total := float64(len(column.Values))
	if total > 0 {
		result.Filled = float64(filledCount) / total
		result.BlankRatio = float64(blankCount) / total
	}

	return result
}

func determineMainType(counts map[DataType]int) DataType {
	var winner DataType = TypeString
	maxCount := 0

	priority := map[DataType]int{
		TypeEmpty:       0,
		TypeString:      1,
		TypeBoolean:     2,
		TypeInteger:     3,
		TypeFloat:       4,
		TypeDate:        5,
		TypeDateCompact: 6,
		TypeEmail:       7,
		TypePlaca:       8,
		TypeCEP:         9,
		TypeCPF:         10,
		TypeCNPJ:        11,
		TypeFiscalKey44: 12,
	}

	for dtype, count := range counts {
		if count > maxCount {
			maxCount = count
			winner = dtype
		} else if count == maxCount {

			if priority[dtype] > priority[winner] {
				winner = dtype
			}
		}
	}

	return winner
}
