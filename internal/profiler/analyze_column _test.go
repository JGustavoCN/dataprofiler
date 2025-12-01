package profiler

import "testing"

func TestAnalyze(t *testing.T) {
	t.Run("Analise inteiro totalmente prenchidos", func(t *testing.T) {
		input := Column{
			"Idades",
			[]string{"1", "2", "2"},
		}

		got := AnalyzeColumn(input)

		expected := ColumnResult{
			Name:     "Idades",
			MainType: "int",
			Filled:   1,
		}

		checkAnalyse(t, got, expected)
	})

	t.Run("Analise inteiros com não preenchidos e espaços", func(t *testing.T) {
		input := Column{
			"Idades",
			[]string{"1   ", "   5  ", "2 ", "     "},
		}

		got := AnalyzeColumn(input)

		expected := ColumnResult{
			Name:     "Idades",
			MainType: "int",
			Filled:   0.75,
		}

		checkAnalyse(t, got, expected)
	})

	t.Run("Analise float com não preenchidos e espaços", func(t *testing.T) {
		input := Column{
			"Idades",
			[]string{"1.5", "   5.4  ", "2.0 ", "     "},
		}

		got := AnalyzeColumn(input)

		expected := ColumnResult{
			Name:     "Idades",
			MainType: "float",
			Filled:   0.75,
		}

		checkAnalyse(t, got, expected)
	})

	t.Run("Analise palavras com ruido nos dados", func(t *testing.T) {
		input := Column{
			"Animais",
			[]string{"cachorro", "  gato  ", "2.0 ", "     "},
		}

		got := AnalyzeColumn(input)

		expected := ColumnResult{
			Name:     "Animais",
			MainType: "string",
			Filled:   0.75,
		}

		checkAnalyse(t, got, expected)
	})

	t.Run("Dados heterogeneos", func(t *testing.T) {
		input := Column{
			"Animais",
			[]string{"cachorro", "  2  ", "2.0 ", "     ", "5", "True", "5", "6", "7", "8"},
		}

		got := AnalyzeColumn(input)

		expected := ColumnResult{
			Name:     "Animais",
			MainType: "int",
			Filled:   0.9,
		}
		checkAnalyse(t, got, expected)
	})

	t.Run("Contagem de tipos", func(t *testing.T) {
		input := Column{
			"Animais",
			[]string{"cachorro", "  2  ", "2.0 ", "     ", "5", "True", "5", "6", "7", "8"},
		}

		got := AnalyzeColumn(input).TypeCounts

		expected := map[string]int{
			"string": 1,
			"bool":   1,
			"int":    6,
			"float":  1,
		}
		checkTypeCounts(t, got, expected)
	})

	t.Run("Caso completo", func(t *testing.T) {
		input := Column{
			"Animais",
			[]string{"cachorro", "  2  ", "2.0 ", "    ", "5", "True", "5", "6", "7", "8"},
		}

		got := AnalyzeColumn(input)

		expected := ColumnResult{
			Name:     "Animais",
			MainType: "int",
			Filled:   0.9,
			TypeCounts: map[string]int{
				"string": 1,
				"bool":   1,
				"int":    6,
				"float":  1,
			},
		}
		checkAnalyse(t, got, expected)
		checkTypeCounts(t, got.TypeCounts, expected.TypeCounts)
	})

	t.Run("Caso da coluna vazia", func(t *testing.T) {
		input := Column{
			"Animais",
			[]string{},
		}

		got := AnalyzeColumn(input)

		expected := ColumnResult{}
		checkAnalyse(t, got, expected)
		checkTypeCounts(t, got.TypeCounts, expected.TypeCounts)
	})

	t.Run("Caso da maioria vazia", func(t *testing.T) {
		input := Column{
			"Animais",
			[]string{"cachorro", "  2  ", "2.0 ", "     ", "5", "    ", " ", "", "", "8"},
		}

		got := AnalyzeColumn(input)

		expected := ColumnResult{
			TypeCounts: map[string]int{
				"string": 1,
				"int":    3,
				"float":  1,
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
			t.Errorf("Integração falhou: Esperava Max 20.00, recebeu %s", got.Stats["Max"])
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
	t.Run("Valor de preenchimento correto", func(t *testing.T) {
		if got.Filled != expected.Filled {
			t.Errorf("Valor de preenchimento correto [esperado: %f - recebeu: %f]", expected.Filled, got.Filled)
		}
	})

}

func checkTypeCounts(t *testing.T, got, expected map[string]int) {
	t.Run("Contagem de tipos", func(t *testing.T) {
		for k, v := range expected {
			if got[k] != v {
				t.Errorf("Contagem errada no Tipo: %s [esperado: %d - recebeu: %d]", k, v, got[k])
			}
		}
	})
}
