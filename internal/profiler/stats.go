package profiler

import "strconv"

func StatsCalc(v []float64) map[string]string {
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

	return map[string]string{
		"Min":     strconv.FormatFloat(min, 'f', 2, 64),
		"Max":     strconv.FormatFloat(max, 'f', 2, 64),
		"Sum":     strconv.FormatFloat(sum, 'f', 2, 64),
		"Average": strconv.FormatFloat(avg, 'f', 2, 64),
	}
}
