package profiler

import (
	"strconv"
	"strings"
)

type ColumnAccumulator struct {
	Name        string
	TotalCount  int
	BlankCount  int
	CountFilled int
	TypeCounts  map[string]int

	numericMin   *float64
	numericMax   *float64
	numericSum   float64
	numericCount int
}

func NewColumnAccumulator(name string) *ColumnAccumulator {
	return &ColumnAccumulator{
		Name:       name,
		TypeCounts: make(map[string]int),
	}
}

func (acc *ColumnAccumulator) Add(value string) {
	acc.TotalCount++
	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" {
		acc.BlankCount++
		return
	}

	inferredType := InferType(trimmedValue)
	acc.TypeCounts[inferredType]++
	acc.CountFilled++

	if inferredType == "int" || inferredType == "float" {
		val, err := strconv.ParseFloat(trimmedValue, 64)
		if err == nil {
			acc.updateNumericStats(val)
		}
	}
}

func (acc *ColumnAccumulator) updateNumericStats(val float64) {
	acc.numericCount++
	acc.numericSum += val

	if acc.numericMin == nil || val < *acc.numericMin {
		v := val
		acc.numericMin = &v
	}

	if acc.numericMax == nil || val > *acc.numericMax {
		v := val
		acc.numericMax = &v
	}
}

func (acc *ColumnAccumulator) Result() (result ColumnResult) {
	winnerType := "string"
	maxCount := 0

	for typeName, count := range acc.TypeCounts {
		if count > maxCount {
			maxCount = count
			winnerType = typeName
		}
	}

	stats := make(map[string]string)

	if winnerType == "int" || winnerType == "float" {
		if acc.numericCount > 0 && acc.numericMin != nil && acc.numericMax != nil {
			stats["Min"] = strconv.FormatFloat(*acc.numericMin, 'f', 2, 64)
			stats["Max"] = strconv.FormatFloat(*acc.numericMax, 'f', 2, 64)
			stats["Sum"] = strconv.FormatFloat(acc.numericSum, 'f', 2, 64)

			avg := acc.numericSum / float64(acc.numericCount)
			stats["Average"] = strconv.FormatFloat(avg, 'f', 2, 64)
		}
	}

	return ColumnResult{
		Name:        acc.Name,
		MainType:    winnerType,
		CountFilled: acc.CountFilled,
		BlankCount:  acc.BlankCount,
		TypeCounts:  acc.TypeCounts,
		Filled:      (float64(acc.CountFilled) / float64(acc.TotalCount)),
		BlankRatio:  (float64(acc.BlankCount) / float64(acc.TotalCount)),
		Stats: stats,
	}
}
