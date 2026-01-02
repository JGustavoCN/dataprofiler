package profiler

import (
	"strconv"
	"strings"
)

type ColumnAccumulator struct {
	Name         string
	TotalCount   int
	BlankCount   int
	CountFilled  int
	TypeCounts   map[DataType]int
	numericMin   *float64
	numericMax   *float64
	numericSum   float64
	numericCount int
}

func NewColumnAccumulator(name string) *ColumnAccumulator {
	return &ColumnAccumulator{
		Name:       name,
		TypeCounts: make(map[DataType]int),
	}
}

func (acc *ColumnAccumulator) Add(value string) {
	acc.TotalCount++

	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" {
		acc.BlankCount++
		return
	}

	acc.CountFilled++

	inferredType := InferType(trimmedValue, acc.Name)
	acc.TypeCounts[inferredType]++

	if inferredType == TypeInteger || inferredType == TypeFloat {
		valClean := strings.Replace(trimmedValue, ",", ".", 1)

		val, err := strconv.ParseFloat(valClean, 64)
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

func (acc *ColumnAccumulator) Result() ColumnResult {
	mainType := acc.determineMainType()
	sensitivity, reasonSensitivity := ClassifySensitivity(mainType)
	stats := make(map[string]string)

	if mainType == TypeInteger || mainType == TypeFloat {
		if acc.numericCount > 0 && acc.numericMin != nil && acc.numericMax != nil {
			stats["Min"] = strconv.FormatFloat(*acc.numericMin, 'f', 2, 64)
			stats["Max"] = strconv.FormatFloat(*acc.numericMax, 'f', 2, 64)
			stats["Sum"] = strconv.FormatFloat(acc.numericSum, 'f', 2, 64)

			avg := acc.numericSum / float64(acc.numericCount)
			stats["Average"] = strconv.FormatFloat(avg, 'f', 2, 64)
		}
	}

	var filledRatio, blankRatio float64
	if acc.TotalCount > 0 {
		filledRatio = float64(acc.CountFilled) / float64(acc.TotalCount)
		blankRatio = float64(acc.BlankCount) / float64(acc.TotalCount)
	}

	consistencyRatio := 1.0
	if acc.CountFilled > 0 {
		winnerCount := acc.TypeCounts[mainType]
		consistencyRatio = float64(winnerCount) / float64(acc.CountFilled)
	}

	sla, reasonSLA := CalculateSLA(blankRatio, consistencyRatio, mainType)
	return ColumnResult{
		Name:              acc.Name,
		MainType:          mainType,
		Sensitivity:       sensitivity,
		SensitivityReason: reasonSensitivity,
		CountFilled:       acc.CountFilled,
		BlankCount:        acc.BlankCount,
		TypeCounts:        acc.TypeCounts,
		Filled:            filledRatio,
		BlankRatio:        blankRatio,
		SLA:               sla,
		SlaReason:         reasonSLA,
		ConsistencyRatio:  consistencyRatio,
		Stats:             stats,
	}
}

func (acc *ColumnAccumulator) determineMainType() DataType {
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

	for dtype, count := range acc.TypeCounts {
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
