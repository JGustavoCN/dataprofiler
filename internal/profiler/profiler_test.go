package profiler

import (
	"testing"
)

func TestProfile(t *testing.T) {
	t.Run("Caminho feliz", func(t *testing.T) {
		inputColumns := []Column{
			{
				"Animais",
				[]string{"cachorro", "  gato  ", "pato", "camelo"},
			},
			{
				"Idades",
				[]string{"1", "2", "2", "4"},
			},
			{
				"Dono",
				[]string{"Joao", "  Alfreado  ", "Lucas ", " Gustavo"},
			},
		}

		inputName := "balanco.csv"

		got := Profile(inputColumns, inputName)

		expected := ProfilerResult{
			NameFile:     "balanco",
			TotalMaxRows:    4,
			TotalColumns: 3,
		}

		checkProfiler(t, got, expected)
	})

	t.Run("Colunas com tamanhos diferentes", func(t *testing.T) {
		inputColumns := []Column{
			{
				"Animais",
				[]string{"cachorro", "  gato  "},
			},
			{
				"Idades",
				[]string{"1", "2"},
			},
			{
				"Dono",
				[]string{"Joao", "  Alfreado  ", "Lucas ", " Gustavo"},
			},
		}

		inputName := "balanco.csv"

		got := Profile(inputColumns, inputName)

		expected := ProfilerResult{
			NameFile:     "balanco",
			TotalMaxRows:    4,
			TotalColumns: 3,
		}

		checkProfiler(t, got, expected)
	})

	t.Run("Colunas vazias", func(t *testing.T) {
		inputColumns := []Column{}

		inputName := ""

		got := Profile(inputColumns, inputName)

		expected := ProfilerResult{}

		checkProfiler(t, got, expected)
	})

	t.Run("Integração com AnalyseColumn em coluna numérica", func(t *testing.T) {
		input := []Column{
			{
				"Animais",
				[]string{"cachorro", "  gato  ", "pato", "camelo"},
			},
			{
				"Idades",
				[]string{"1", "2", "2", "4"},
			},
			{
				"Dono",
				[]string{"Joao", "  Alfreado  ", "Lucas ", " Gustavo"},
			},
		}

		got := Profile(input, "")
		if got.Columns == nil {
			t.Fatal("Esperava analise das colunas, mas veio nil")
		}

		if got.Columns[0].MainType != "string" {
			t.Errorf("Integração falhou: Esperava no tipo principal \"string\", recebeu %s", got.Columns[0].MainType)
		}
	})

}

func checkProfiler(t *testing.T, got, expected ProfilerResult) {

	t.Run("Nome correto", func(t *testing.T) {
		if got.NameFile != expected.NameFile {
			t.Errorf("Erro no nome: esperado [%s] - recebido [%s]", expected.NameFile, got.NameFile)
		}
	})
	t.Run("Contagem de linhas correto", func(t *testing.T) {
		if got.TotalMaxRows != expected.TotalMaxRows {
			t.Errorf("Erro na contagem de linhas: esperado [%d] - recebido [%d]", expected.TotalMaxRows, got.TotalMaxRows)
		}
	})
	t.Run("Contagem de colunas correto", func(t *testing.T) {
		if got.TotalColumns != expected.TotalColumns {
			t.Errorf("Erro na contagem de colunas: esperado [%d] - recebido [%d]", expected.TotalColumns, got.TotalColumns)
		}
	})
}
