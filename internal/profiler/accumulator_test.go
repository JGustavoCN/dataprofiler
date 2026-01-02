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

		if count := result.TypeCounts[TypeInteger]; count != 2 {
			t.Errorf("Esperado 2 ints, recebeu %d", count)
		}

		if count := result.TypeCounts[TypeString]; count != 1 {
			t.Errorf("Esperado 1 string, recebeu %d", count)
		}
	})

	t.Run("Deve calcular estatísticas de números positivos", func(t *testing.T) {
		acc := NewColumnAccumulator("Precos")

		acc.Add("10")
		acc.Add("20")
		acc.Add("30")

		result := acc.Result()

		if result.MainType != TypeInteger {
			t.Errorf("Esperava MainType INTEGER, recebeu %s", result.MainType)
		}

		if result.Stats == nil {
			t.Fatal("Stats não deveria ser nil para colunas numéricas")
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

		if result.MainType != TypeFloat {
			t.Errorf("Esperava MainType FLOAT, recebeu %s", result.MainType)
		}

		expectedStats := map[string]string{
			"Min":     "-20.00",
			"Max":     "-5.00",
			"Sum":     "-35.50",
			"Average": "-11.83",
		}

		checkStatsMap(t, result.Stats, expectedStats)
	})

	t.Run("Deve suportar formato Brasileiro (Vírgula)", func(t *testing.T) {

		acc := NewColumnAccumulator("Frete")

		acc.Add("10,50")
		acc.Add("20,50")

		result := acc.Result()

		if result.MainType != TypeFloat {
			t.Errorf("Deveria ter detectado FLOAT mesmo com vírgula, recebeu: %s", result.MainType)
		}

		expectedStats := map[string]string{
			"Sum":     "31.00",
			"Average": "15.50",
		}

		checkStatsMap(t, result.Stats, expectedStats)
	})

	t.Run("NÃO deve calcular Stats para String", func(t *testing.T) {
		acc := NewColumnAccumulator("Nomes")
		acc.Add("Ana")
		acc.Add("Bia")

		result := acc.Result()

		if len(result.Stats) > 0 {
			t.Errorf("Não deveria ter estatísticas numéricas para texto! Recebeu: %v", result.Stats)
		}
	})
}

func checkStatsMap(t *testing.T, got, expected map[string]string) {
	t.Helper()
	for k, v := range expected {
		if got[k] != v {
			t.Errorf("Stats[%s]: recebeu %s, esperava %s", k, got[k], v)
		}
	}
}
