package profiler

import "testing"

func TestStatsCalc(t *testing.T) {
	t.Run("Deve calcular estatísticas de números positivos", func(t *testing.T) {
		input := []float64{10.0, 20.0, 30.0}
		got := StatsCalc(input)

		expected := map[string]string{
			"Min":     "10.00",
			"Max":     "30.00",
			"Sum":     "60.00",
			"Average": "20.00",
		}

		checkStats(t, got, expected)
	})

	t.Run("Deve calcular estatísticas com números negativos", func(t *testing.T) {
		input := []float64{-10.0, -5.0, -20.0}
		got := StatsCalc(input)

		expected := map[string]string{
			"Min":     "-20.00",
			"Max":     "-5.00",
			"Sum":     "-35.00",
			"Average": "-11.67",
		}

		checkStats(t, got, expected)
	})
}

func checkStats(t *testing.T, got, expected map[string]string) {
	for k, v := range expected {
		if got[k] != v {
			t.Errorf("Para chave %s: recebeu %s, esperava %s", k, got[k], v)
		}
	}
}
