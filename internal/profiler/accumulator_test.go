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

	t.Run("Integração SLA: Alta Severidade (CPF) com Inconsistência", func(t *testing.T) {
		acc := NewColumnAccumulator("Documentos")

		acc.Add("123.456.789-00")
		acc.Add("111.222.333-44")
		acc.Add("Não Informado")

		result := acc.Result()

		if result.MainType != TypeCPF {
			t.Errorf("Falha na inferência predominante. Esperava CPF, deu %s", result.MainType)
		}

		expectedConsistency := 2.0 / 3.0
		if result.ConsistencyRatio != expectedConsistency {
			t.Errorf("Consistência errada. Esperava %f, deu %f", expectedConsistency, result.ConsistencyRatio)
		}

		if result.SLA != SlaCritical {
			t.Errorf("Esperava SLA CRITICAL para CPF sujo, recebeu %s. Motivo: %s", result.SLA, result.SlaReason)
		}
	})

	t.Run("Integração SLA: Média Severidade (Integer) com Warning", func(t *testing.T) {
		acc := NewColumnAccumulator("Idades")

		for i := 0; i < 98; i++ {
			acc.Add("10")
		}

		acc.Add("erro")
		acc.Add("n/a")

		result := acc.Result()

		if result.MainType != TypeInteger {
			t.Fatalf("Esperava MainType INTEGER, deu %s", result.MainType)
		}

		acc.Add("sujeira_extra")

		result = acc.Result()
		if result.SLA != SlaWarning {
			t.Errorf("Esperava SLA WARNING para Integer com consistência 0.97, recebeu %s", result.SLA)
		}
	})

	t.Run("Integração SLA: Baixa Severidade (String) Tolerante", func(t *testing.T) {
		acc := NewColumnAccumulator("Observacoes")

		acc.Add("Obs 1")
		acc.Add("Obs 2")
		acc.Add("")
		acc.Add("")

		result := acc.Result()

		if result.MainType != TypeString {
			t.Fatalf("Esperava MainType STRING, deu %s", result.MainType)
		}

		if result.SLA != SlaWarning {
			t.Errorf("Esperava SLA WARNING (não Critical) para String 50%% vazia, recebeu %s", result.SLA)
		}
	})

	t.Run("Integração Sensitivity: Deve classificar PII e Dados Internos corretamente", func(t *testing.T) {

		accPii := NewColumnAccumulator("Documento_Cliente")
		accPii.Add("123.456.789-00")
		accPii.Add("111.222.333-44")

		resPii := accPii.Result()

		if resPii.MainType != TypeCPF {
			t.Fatalf("Pré-requisito falhou: Deveria ter detectado CPF, mas detectou %s", resPii.MainType)
		}

		if resPii.Sensitivity != SensitivityConfidential {
			t.Errorf("Falha de Governança! CPF deveria ser CONFIDENTIAL, mas foi classificado como %s", resPii.Sensitivity)
		}

		accInternal := NewColumnAccumulator("Placa_Caminhao")
		accInternal.Add("ABC-1234")
		accInternal.Add("MER-1D23")

		resInternal := accInternal.Result()

		if resInternal.Sensitivity != SensitivityInternal {
			t.Errorf("Falha de Negócio! Placa deveria ser INTERNAL, mas foi classificada como %s", resInternal.Sensitivity)
		}

		accPublic := NewColumnAccumulator("Quantidade_Estoque")
		accPublic.Add("100")
		accPublic.Add("500")

		resPublic := accPublic.Result()

		if resPublic.Sensitivity != SensitivityPublic {
			t.Errorf("Falha de Ruído! Inteiro deveria ser PUBLIC, mas foi classificado como %s", resPublic.Sensitivity)
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
