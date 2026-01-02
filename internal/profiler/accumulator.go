package profiler

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"
)

type ColumnAccumulator struct {
	Name          string
	TotalCount    int
	BlankCount    int
	CountFilled   int
	TypeCounts    map[DataType]int
	numericMin    *float64
	numericMax    *float64
	numericSum    float64
	numericCount  int
	numericSample []float64
	sampleSize    int
	rng           *rand.Rand
}

func NewColumnAccumulator(name string) *ColumnAccumulator {
	seed := uint64(time.Now().UnixNano())
	return &ColumnAccumulator{
		Name:          name,
		TypeCounts:    make(map[DataType]int),
		numericSample: make([]float64, 0, 1000),
		sampleSize:    1000,
		rng:           rand.New(rand.NewPCG(seed, seed+1)),
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

	if len(acc.numericSample) < acc.sampleSize {
		acc.numericSample = append(acc.numericSample, val)
	} else {
		randomIndex := acc.rng.IntN(acc.numericCount)
		if randomIndex < acc.sampleSize {
			acc.numericSample[randomIndex] = val
		}
	}
}

func (acc *ColumnAccumulator) Result() ColumnResult {
	mainType := acc.determineMainType()
	sensitivity, reasonSensitivity := ClassifySensitivity(mainType)
	stats := make(map[StatKey]string)
	var histogram map[string]int
	if mainType == TypeInteger || mainType == TypeFloat {
		if acc.numericCount > 0 && acc.numericMin != nil && acc.numericMax != nil {
			stats[StatMin] = strconv.FormatFloat(*acc.numericMin, 'f', 2, 64)
			stats[StatMax] = strconv.FormatFloat(*acc.numericMax, 'f', 2, 64)
			stats[StatSum] = strconv.FormatFloat(acc.numericSum, 'f', 2, 64)

			avg := acc.numericSum / float64(acc.numericCount)
			stats[StatAverage] = strconv.FormatFloat(avg, 'f', 2, 64)
			histogram = calculateHistogram(acc.numericSample)
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
		Histogram:         histogram,
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

func calculateHistogram(values []float64) map[string]int {
	if len(values) == 0 {
		return nil
	}

	minVal, maxVal := values[0], values[0]
	for _, v := range values {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}

	numBuckets := 10
	rangeVal := maxVal - minVal

	if rangeVal == 0 {
		return map[string]int{fmt.Sprintf("%.2f", minVal): len(values)}
	}

	step := rangeVal / float64(numBuckets)
	histogram := make(map[string]int)

	for _, v := range values {
		bucketIndex := int((v - minVal) / step)
		if bucketIndex >= numBuckets {
			bucketIndex = numBuckets - 1
		}

		bucketStart := minVal + (float64(bucketIndex) * step)
		bucketEnd := bucketStart + step
		label := fmt.Sprintf("%.2f-%.2f", bucketStart, bucketEnd)

		histogram[label]++
	}
	return histogram
}
