package profiler

import "testing"

func TestAnalyze(t *testing.T) {

	t.Run("Analise inteiros com não preenchidos (Validando SLA)", func(t *testing.T) {
		input := Column{
			Name:   "Idades",
			Values: []string{"1", "5", "2", ""},
		}

		got := AnalyzeColumn(input)

		if got.SLA != SlaCritical {
			t.Errorf("Esperava SLA CRITICAL para 25%% de vazio em Integer, recebeu %s", got.SLA)
		}
		if got.ConsistencyRatio != 1.0 {
			t.Errorf("Esperava consistência 1.0 (todos são int), recebeu %f", got.ConsistencyRatio)
		}
	})

	t.Run("Dados heterogeneos (Teste de Consistência)", func(t *testing.T) {

		input := Column{
			Name:   "NumerosMisturados",
			Values: []string{"1", "2", "2.5", "texto", "5", "True", "5", "6", "7", ""},
		}

		got := AnalyzeColumn(input)

		expectedConsistency := 6.0 / 9.0

		if got.MainType != TypeInteger {
			t.Errorf("Esperava MainType INTEGER, veio %s", got.MainType)
		}

		if got.ConsistencyRatio != expectedConsistency {
			t.Errorf("Consistência errada. Esperava %f, veio %f", expectedConsistency, got.ConsistencyRatio)
		}

		if got.SLA != SlaCritical {
			t.Errorf("Esperava SLA CRITICAL por baixa consistência, recebeu %s", got.SLA)
		}

		if got.SlaReason == "" {
			t.Error("SlaReason não deveria vir vazio")
		}
	})

	t.Run("Analise inteiro totalmente prenchidos", func(t *testing.T) {
		input := Column{
			Name:   "Idades",
			Values: []string{"1", "2", "2"},
		}

		got := AnalyzeColumn(input)

		expected := ColumnResult{
			Name:        "Idades",
			MainType:    TypeInteger,
			BlankCount:  0,
			CountFilled: 3,
			Filled:      1,
			BlankRatio:  0,
		}

		checkAnalyse(t, got, expected)
	})

	t.Run("Analise inteiros com não preenchidos e espaços", func(t *testing.T) {
		input := Column{
			Name:   "Idades",
			Values: []string{"1   ", "   5  ", "2 ", "     "},
		}

		got := AnalyzeColumn(input)

		expected := ColumnResult{
			Name:        "Idades",
			MainType:    TypeInteger,
			BlankCount:  1,
			CountFilled: 3,
			Filled:      0.75,
			BlankRatio:  0.25,
		}

		checkAnalyse(t, got, expected)
	})

	t.Run("Analise float com não preenchidos e espaços", func(t *testing.T) {
		input := Column{
			Name:   "Idades",
			Values: []string{"1.5", "   5.4  ", "2.0 ", "     "},
		}

		got := AnalyzeColumn(input)

		expected := ColumnResult{
			Name:        "Idades",
			MainType:    TypeFloat,
			BlankCount:  1,
			CountFilled: 3,
			Filled:      0.75,
			BlankRatio:  0.25,
		}

		checkAnalyse(t, got, expected)
	})

	t.Run("Analise palavras com ruido nos dados", func(t *testing.T) {
		input := Column{
			Name:   "Animais",
			Values: []string{"cachorro", "  gato  ", "2.0 ", "     "},
		}

		got := AnalyzeColumn(input)

		expected := ColumnResult{
			Name:        "Animais",
			MainType:    TypeString,
			BlankCount:  1,
			CountFilled: 3,
			Filled:      0.75,
			BlankRatio:  0.25,
		}

		checkAnalyse(t, got, expected)
	})

	t.Run("Dados heterogeneos", func(t *testing.T) {
		input := Column{
			Name:   "Animais",
			Values: []string{"cachorro", "  2  ", "2.0 ", "     ", "5", "True", "5", "6", "7", "8"},
		}

		got := AnalyzeColumn(input)

		expected := ColumnResult{
			Name:        "Animais",
			MainType:    TypeInteger,
			BlankCount:  1,
			CountFilled: 9,
			Filled:      0.9,
			BlankRatio:  0.1,
		}
		checkAnalyse(t, got, expected)
	})

	t.Run("Contagem de tipos", func(t *testing.T) {
		input := Column{
			Name:   "Animais",
			Values: []string{"cachorro", "  2  ", "2.0 ", "     ", "5", "True", "5", "6", "7", "8"},
		}

		got := AnalyzeColumn(input).TypeCounts

		expected := map[DataType]int{
			TypeString:  1,
			TypeBoolean: 1,
			TypeInteger: 6,
			TypeFloat:   1,
		}
		checkTypeCounts(t, got, expected)
	})

	t.Run("Caso completo", func(t *testing.T) {
		input := Column{
			Name:   "Animais",
			Values: []string{"cachorro", "  2  ", "2.0 ", "    ", "5", "True", "5", "6", "7", "8"},
		}

		got := AnalyzeColumn(input)

		expected := ColumnResult{
			Name:        "Animais",
			MainType:    TypeInteger,
			BlankCount:  1,
			CountFilled: 9,
			Filled:      0.9,
			BlankRatio:  0.1,
			TypeCounts: map[DataType]int{
				TypeString:  1,
				TypeBoolean: 1,
				TypeInteger: 6,
				TypeFloat:   1,
			},
		}
		checkAnalyse(t, got, expected)
		checkTypeCounts(t, got.TypeCounts, expected.TypeCounts)
	})

	t.Run("Caso metade int e metade string, mantém o primeiro tipo que apareceu", func(t *testing.T) {
		input := Column{
			Name:   "Animais",
			Values: []string{"2", "  w ", "2.0 ", "    ", "5", "True", "jkjn", "kjnk", "7", ""},
		}

		got := AnalyzeColumn(input)

		expected := ColumnResult{
			Name:        "Animais",
			MainType:    TypeInteger,
			BlankCount:  2,
			CountFilled: 8,
			Filled:      0.8,
			BlankRatio:  0.2,
			TypeCounts: map[DataType]int{
				TypeString:  3,
				TypeBoolean: 1,
				TypeInteger: 3,
				TypeFloat:   1,
			},
		}
		checkAnalyse(t, got, expected)
		checkTypeCounts(t, got.TypeCounts, expected.TypeCounts)
	})

	t.Run("Caso da coluna vazia", func(t *testing.T) {
		input := Column{
			Name:   "Animais",
			Values: []string{},
		}

		got := AnalyzeColumn(input)

		expected := ColumnResult{
			Name:     "Animais",
			MainType: TypeEmpty,
		}
		checkAnalyse(t, got, expected)
	})

	t.Run("Caso da maioria vazia", func(t *testing.T) {
		input := Column{
			Name:   "Animais",
			Values: []string{"cachorro", "  2  ", "2.0 ", "     ", "5", "    ", " ", "", "", "8"},
		}

		got := AnalyzeColumn(input)

		expected := ColumnResult{
			TypeCounts: map[DataType]int{
				TypeString:  1,
				TypeInteger: 3,
				TypeFloat:   1,
			},
		}
		checkTypeCounts(t, got.TypeCounts, expected.TypeCounts)
	})

	t.Run("Integração com StatsCalc em coluna numérica", func(t *testing.T) {
		input := Column{
			Name:   "Precos",
			Values: []string{" 10", "20", "30"},
		}

		got := AnalyzeColumn(input)
		if got.Stats == nil {
			t.Fatal("Esperava estatísticas para coluna numérica, mas veio nil")
		}

		if got.Stats["Average"] != "20.00" {
			t.Errorf("Integração falhou: Esperava Average 20.00, recebeu %s", got.Stats["Average"])
		}
	})

	t.Run("NÃO deve calcular Stats para String", func(t *testing.T) {
		input := Column{
			Name:   "Nomes",
			Values: []string{"Ana", "Bia"},
		}

		got := AnalyzeColumn(input)

		if got.Stats != nil {
			t.Errorf("Não deveria ter calculado estatísticas para texto!")
		}
	})
}

func checkAnalyse(t *testing.T, got, expected ColumnResult) {
	t.Helper()
	t.Run("Nome correto", func(t *testing.T) {
		if got.Name != expected.Name {
			t.Errorf("Nome errado [esperado: %s - recebeu: %s]", expected.Name, got.Name)
		}
	})
	t.Run("Tipo correto", func(t *testing.T) {
		if got.MainType != expected.MainType {
			t.Errorf("Tipo correto [esperado: %s - recebeu: %s]", expected.MainType, got.MainType)
		}
	})
	t.Run("Valores de preenchimento", func(t *testing.T) {
		if got.CountFilled != expected.CountFilled {
			t.Errorf("CountFilled [esperado: %d - recebeu: %d]", expected.CountFilled, got.CountFilled)
		}
		if got.BlankCount != expected.BlankCount {
			t.Errorf("BlankCount [esperado: %d - recebeu: %d]", expected.BlankCount, got.BlankCount)
		}
		if got.Filled != expected.Filled {
			t.Errorf("Filled [esperado: %f - recebeu: %f]", expected.Filled, got.Filled)
		}
		if got.BlankRatio != expected.BlankRatio {
			t.Errorf("BlankRatio [esperado: %f - recebeu: %f]", expected.BlankRatio, got.BlankRatio)
		}
	})
}

func checkTypeCounts(t *testing.T, got, expected map[DataType]int) {
	t.Helper()
	t.Run("Contagem de tipos", func(t *testing.T) {
		for k, v := range expected {
			if got[k] != v {
				t.Errorf("Contagem errada no Tipo: %s [esperado: %d - recebeu: %d]", k, v, got[k])
			}
		}
	})
}
