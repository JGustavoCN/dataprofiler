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
	Name              string            `json:"name"`
	MainType          DataType          `json:"main_type"`
	Sensitivity       DataSensitivity   `json:"sensitivity_level"`
	SensitivityReason string            `json:"sensitivity_reason"`
	SLA               QualityScore      `json:"sla"`
	SlaReason         string            `json:"sla_reason"`
	BlankCount        int               `json:"blank_count"`
	CountFilled       int               `json:"count_filled"`
	Filled            float64           `json:"filled_ratio"`
	BlankRatio        float64           `json:"blank_ratio"`
	ConsistencyRatio  float64           `json:"consistency_ratio"`
	TypeCounts        map[DataType]int  `json:"type_counts"`
	Stats             map[string]string `json:"stats,omitempty"`
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
	result.Sensitivity, result.SensitivityReason = ClassifySensitivity(result.MainType)

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

		result.ConsistencyRatio = 1.0
		if filledCount > 0 {
			winnerCount := result.TypeCounts[result.MainType]
			result.ConsistencyRatio = float64(winnerCount) / float64(filledCount)
		}
		result.SLA, result.SlaReason = CalculateSLA(result.BlankRatio, result.ConsistencyRatio, result.MainType)

	} else {
		result.SLA = SlaGood
		result.ConsistencyRatio = 1.0
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
