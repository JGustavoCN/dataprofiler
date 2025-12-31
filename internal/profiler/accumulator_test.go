package profiler

import (
	"testing"
)

func TestAccumulator(t *testing.T) {

	t.Run("Contagem e Inferência Simples", func(t *testing.T) {
		acc := NewColumnAccumulator("Idade")

		acc.Add("10")
		acc.Add("20")
		acc.Add("")
		acc.Add("Ola")

		result := acc.Result()

		if result.BlankCount != 1 {
			t.Errorf("Esperado 1 vazio, recebeu %d", result.BlankCount)
		}

		if result.TypeCounts["int"] != 2 {
			t.Errorf("Esperado 2 ints, recebeu %d", result.TypeCounts["int"])
		}

		if result.TypeCounts["string"] != 1 {
			t.Errorf("Esperado 1 string, recebeu %d", result.TypeCounts["string"])
		}
	})

	t.Run("Deve calcular estatísticas de números positivos", func(t *testing.T) {
		acc := NewColumnAccumulator("Precos")

		acc.Add("10")
		acc.Add("20")
		acc.Add("30")

		result := acc.Result()

		if result.MainType != "int" {
			t.Errorf("Esperava MainType int, recebeu %s", result.MainType)
		}

		if result.Stats == nil {
			t.Fatal("Stats não deveria ser nil")
		}

		expectedStats := map[string]string{
			"Min":     "10.00",
			"Max":     "30.00",
			"Sum":     "60.00",
			"Average": "20.00",
		}

		checkStatsMap(t, result.Stats, expectedStats)
	})

	t.Run("Deve calcular estatísticas com números negativos e float", func(t *testing.T) {
		acc := NewColumnAccumulator("Temperaturas")

		acc.Add("-10.5")
		acc.Add("-5.0")
		acc.Add("-20")

		result := acc.Result()

		expectedStats := map[string]string{
			"Min":     "-20.00",
			"Max":     "-5.00",
			"Sum":     "-35.50",
			"Average": "-11.83",
		}

		checkStatsMap(t, result.Stats, expectedStats)
	})

	t.Run("NÃO deve calcular Stats para String", func(t *testing.T) {
		acc := NewColumnAccumulator("Nomes")
		acc.Add("Ana")
		acc.Add("Bia")

		result := acc.Result()

		if len(result.Stats) > 2 {
			if _, ok := result.Stats["Average"]; ok {
				t.Errorf("Não deveria ter calculado Average para texto!")
			}
		}
	})
}

func checkStatsMap(t *testing.T, got, expected map[string]string) {
	for k, v := range expected {
		if got[k] != v {
			t.Errorf("Stats[%s]: recebeu %s, esperava %s", k, got[k], v)
		}
	}
}
