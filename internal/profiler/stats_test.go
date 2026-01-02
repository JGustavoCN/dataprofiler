package profiler

import "testing"

func TestStatsCalc(t *testing.T) {
	t.Run("Deve calcular estatísticas de números positivos", func(t *testing.T) {
		input := []float64{10.0, 20.0, 30.0}
		got := StatsCalc(input)

		expected := map[StatKey]string{
			StatMin:     "10.00",
			StatMax:     "30.00",
			StatSum:     "60.00",
			StatAverage: "20.00",
		}

		checkStats(t, got, expected)
	})

	t.Run("Deve calcular estatísticas com números negativos", func(t *testing.T) {
		input := []float64{-10.0, -5.0, -20.0}
		got := StatsCalc(input)

		expected := map[StatKey]string{
			StatMin:     "-20.00",
			StatMax:     "-5.00",
			StatSum:     "-35.00",
			StatAverage: "-11.67",
		}

		checkStats(t, got, expected)
	})
}

func checkStats(t *testing.T, got, expected map[StatKey]string) {
	for k, v := range expected {
		if got[k] != v {
			t.Errorf("Para chave %s: recebeu %s, esperava %s", k, got[k], v)
		}
	}
}
