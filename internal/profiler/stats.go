package profiler

import "strconv"

type StatKey string

const (
	StatMin     StatKey = "min"
	StatMax     StatKey = "max"
	StatSum     StatKey = "sum"
	StatAverage StatKey = "average"
)

func StatsCalc(v []float64) map[StatKey]string {
	if len(v) == 0 {
		return nil
	}

	min, max := v[0], v[0]
	sum := 0.0

	for _, valor := range v {
		if valor < min {
			min = valor
		}
		if valor > max {
			max = valor
		}
		sum += valor
	}

	avg := sum / float64(len(v))

	return map[StatKey]string{
		StatMin:     strconv.FormatFloat(min, 'f', 2, 64),
		StatMax:     strconv.FormatFloat(max, 'f', 2, 64),
		StatSum:     strconv.FormatFloat(sum, 'f', 2, 64),
		StatAverage: strconv.FormatFloat(avg, 'f', 2, 64),
	}
}
