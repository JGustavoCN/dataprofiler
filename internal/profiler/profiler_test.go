package profiler

import (
	"errors"
	"io"
	"log/slog"
	"testing"
)

func TestProfileAsync(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	t.Run("Caminho Feliz - Dados Válidos", func(t *testing.T) {
		inputHeaders := []string{"name", "idade"}
		inputDataChan := make(chan StreamData)
		inputName := "balanco.csv"

		go func() {
			defer close(inputDataChan)
			inputDataChan <- StreamData{Row: []string{"Joao", "10"}}
			inputDataChan <- StreamData{Row: []string{"Gustavo", "21"}}
		}()

		got := ProfileAsync(logger, inputHeaders, inputDataChan, inputName)
		expected := ProfilerResult{
			NameFile: "balanco",
			Columns: []ColumnResult{
				{
					Name: "name",
				},
				{
					Name: "idade",
				},
			},
			TotalMaxRows: 2,
			TotalColumns: 2,
		}

		checkProfiler(t, got, expected)
		t.Run("Cabeçalho correto", func(t *testing.T) {
			if got.Columns[0].Name != expected.Columns[0].Name {
				t.Errorf("Erro no nome do cabeçalho: esperado [%s] - recebido [%s]", expected.Columns[0].Name, got.Columns[0].Name)
			}
			if got.Columns[1].Name != expected.Columns[1].Name {
				t.Errorf("Erro no nome do cabeçalho: esperado [%s] - recebido [%s]", expected.Columns[0].Name, got.Columns[0].Name)
			}

		})
		if got.TotalMaxRows != 2 {
			t.Errorf("Esperava 2 linhas processadas, obteve %d", got.TotalMaxRows)
		}
		if got.DirtyLinesCount != 0 {
			t.Errorf("Não esperava DirtyLines, obteve %d", got.DirtyLinesCount)
		}
	})
	t.Run("Deve registrar Dirty Lines e continuar processando", func(t *testing.T) {
		headers := []string{"nome", "email"}
		dataChan := make(chan StreamData)
		fileName := "sujo.csv"

		go func() {
			defer close(dataChan)
			dataChan <- StreamData{Row: []string{"Ana", "ana@teste.com"}, LineNumber: 2}
			dataChan <- StreamData{Err: errors.New("record on line 3: wrong number of fields"), LineNumber: 3}
			dataChan <- StreamData{Row: []string{"Bia", "bia@teste.com"}, LineNumber: 4}
		}()

		result := ProfileAsync(logger, headers, dataChan, fileName)
		if result.TotalMaxRows != 2 {
			t.Errorf("Deveria ter processado apenas as 2 linhas válidas, contou %d", result.TotalMaxRows)
		}

		if result.DirtyLinesCount != 1 {
			t.Errorf("Deveria ter encontrado 1 DirtyLine, encontrou %d", result.DirtyLinesCount)
		}

		if len(result.DirtyLines) > 0 {
			dirty := result.DirtyLines[0]
			if dirty.Line != 3 {
				t.Errorf("A linha do erro deveria ser 3, foi %d", dirty.Line)
			}
			if dirty.Reason == "" {
				t.Error("O motivo do erro (Reason) não foi preenchido")
			}
		}
	})
}

func TestProfileAsync_Integration(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	t.Run("Deve calcular estatísticas completas via streaming", func(t *testing.T) {

		headers := []string{"Produto", "Preco"}
		fileName := "vendas.csv"

		dataChan := make(chan StreamData)

		go func() {
			defer close(dataChan)
			dataChan <- StreamData{Row: []string{"TV", "1000.00"}}
			dataChan <- StreamData{Row: []string{"Radio", "200.00"}}
			dataChan <- StreamData{Row: []string{"Celular", "1800.00"}}
			dataChan <- StreamData{Row: []string{"Cabo", "50.00"}}
		}()

		result := ProfileAsync(logger, headers, dataChan, fileName)

		if result.TotalMaxRows != 4 {
			t.Errorf("Esperado 4 linhas, recebeu %d", result.TotalMaxRows)
		}
		colProduto := result.Columns[0]
		if colProduto.MainType != "string" {
			t.Errorf("Coluna Produto deveria ser string, foi %s", colProduto.MainType)
		}

		colPreco := result.Columns[1]
		if colPreco.MainType != "float" && colPreco.MainType != "int" {
			t.Errorf("Coluna Preco deveria ser numérica, foi %s", colPreco.MainType)
		}

		stats := colPreco.Stats
		if stats == nil {
			t.Fatal("Stats da coluna Preco não deveria ser nil")
		}

		if stats["Sum"] != "3050.00" {
			t.Errorf("Soma incorreta. Esperado 3050.00, recebeu %s", stats["Sum"])
		}

		if stats["Average"] != "762.50" {
			t.Errorf("Média incorreta. Esperado 762.50, recebeu %s", stats["Average"])
		}

		if stats["Min"] != "50.00" {
			t.Errorf("Min incorreto. Esperado 50.00, recebeu %s", stats["Min"])
		}
	})
}

func TestProfile(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
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

		got := Profile(logger, inputColumns, inputName)

		expected := ProfilerResult{
			NameFile:     "balanco",
			TotalMaxRows: 4,
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

		got := Profile(logger, inputColumns, inputName)

		expected := ProfilerResult{
			NameFile:     "balanco",
			TotalMaxRows: 4,
			TotalColumns: 3,
		}

		checkProfiler(t, got, expected)
	})

	t.Run("Colunas vazias", func(t *testing.T) {
		inputColumns := []Column{}

		inputName := ""

		got := Profile(logger, inputColumns, inputName)

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

		got := Profile(logger, input, "")
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
	t.Run("Contagem de colunas correto", func(t *testing.T) {
		if got.TotalColumns != expected.TotalColumns {
			t.Errorf("Erro na contagem de colunas: esperado [%d] - recebido [%d]", expected.TotalColumns, got.TotalColumns)
		}
	})
	t.Run("Contagem de linhas correto", func(t *testing.T) {
		if got.TotalMaxRows != expected.TotalMaxRows {
			t.Errorf("Erro na contagem de linhas: esperado [%d] - recebido [%d]", expected.TotalMaxRows, got.TotalMaxRows)
		}
	})
}
