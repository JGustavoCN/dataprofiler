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
				{Name: "name"},
				{Name: "idade"},
			},
			TotalMaxRows: 2,
			TotalColumns: 2,
		}

		checkProfiler(t, got, expected)

		t.Run("Cabeçalho correto", func(t *testing.T) {
			if len(got.Columns) != 2 {
				t.Fatalf("Esperava 2 colunas, recebeu %d", len(got.Columns))
			}
			if got.Columns[0].Name != expected.Columns[0].Name {
				t.Errorf("Erro no nome do cabeçalho: esperado [%s] - recebido [%s]", expected.Columns[0].Name, got.Columns[0].Name)
			}
			if got.Columns[1].Name != expected.Columns[1].Name {
				t.Errorf("Erro no nome do cabeçalho: esperado [%s] - recebido [%s]", expected.Columns[1].Name, got.Columns[1].Name)
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

		if len(result.Columns) < 2 {
			t.Fatalf("Esperava 2 colunas analisadas, recebeu %d", len(result.Columns))
		}

		colProduto := result.Columns[0]

		if colProduto.MainType != TypeString {
			t.Errorf("Coluna Produto deveria ser STRING, foi %s", colProduto.MainType)
		}

		colPreco := result.Columns[1]

		if colPreco.MainType != TypeFloat && colPreco.MainType != TypeInteger {
			t.Errorf("Coluna Preco deveria ser numérica (FLOAT/INT), foi %s", colPreco.MainType)
		}

		stats := colPreco.Stats
		if stats == nil {
			t.Fatal("Stats da coluna Preco não deveria ser nil")
		}

		if stats[StatSum] != "3050.00" {
			t.Errorf("Soma incorreta. Esperado 3050.00, recebeu %s", stats[StatSum])
		}

		if stats[StatAverage] != "762.50" {
			t.Errorf("Média incorreta. Esperado 762.50, recebeu %s", stats[StatAverage])
		}
		if len(result.SampleRows) != 4 {
			t.Errorf("Esperava 4 linhas de amostra (total do arquivo < 50), recebeu %d", len(result.SampleRows))
		}

		// Verifica se copiou os dados corretamente (Deep Check)
		if result.SampleRows[0][0] != "TV" {
			t.Errorf("Amostra corrompida ou fora de ordem. Esperava 'TV', recebeu '%s'", result.SampleRows[0][0])
		}
	})
}

func TestProfile(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	t.Run("Caminho feliz", func(t *testing.T) {
		inputColumns := []Column{
			{Name: "Animais", Values: []string{"cachorro", "  gato  ", "pato", "camelo"}},
			{Name: "Idades", Values: []string{"1", "2", "2", "4"}},
			{Name: "Dono", Values: []string{"Joao", "  Alfreado  ", "Lucas ", " Gustavo"}},
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
			{Name: "Animais", Values: []string{"cachorro", "  gato  "}},
			{Name: "Idades", Values: []string{"1", "2"}},

			{Name: "Dono", Values: []string{"Joao", "  Alfreado  ", "Lucas ", " Gustavo"}},
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
			{Name: "Animais", Values: []string{"cachorro", "  gato  "}},
			{Name: "Idades", Values: []string{"1", "2", "2", "4"}},
		}

		got := Profile(logger, input, "")
		if got.Columns == nil {
			t.Fatal("Esperava analise das colunas, mas veio nil")
		}

		if got.Columns[0].MainType != TypeString {
			t.Errorf("Integração falhou: Esperava no tipo principal STRING, recebeu %s", got.Columns[0].MainType)
		}

		if got.Columns[1].MainType != TypeInteger {
			t.Errorf("Integração falhou: Esperava no tipo principal INTEGER, recebeu %s", got.Columns[1].MainType)
		}
	})

	t.Run("Deve gerar Histograma Exato para Coluna Numérica (Modo Síncrono)", func(t *testing.T) {

		values := []string{"10", "20", "30", "40", "50"}
		input := Column{
			Name:   "Produtos",
			Values: values,
		}

		got := AnalyzeColumn(input)

		if got.Histogram == nil {
			t.Fatal("Histograma não foi gerado para coluna numérica")
		}

		totalCounts := 0
		for _, v := range got.Histogram {
			totalCounts += v
		}

		if totalCounts != 5 {
			t.Errorf("Histograma perdeu dados. Esperava total 5, contou %d", totalCounts)
		}
	})
}

func checkProfiler(t *testing.T, got, expected ProfilerResult) {
	t.Helper()
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
