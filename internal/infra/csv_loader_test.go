package infra

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/JGustavoCN/dataprofiler/internal/profiler"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func FuzzDetectSeparator(f *testing.F) {
	f.Add("col1,col2\nval1,val2")
	f.Add("col1;col2;col3")
	f.Add("apenas um texto sem separador")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		reader := strings.NewReader(input)
		_, _ = DetectSeparator(bufio.NewReaderSize(reader, 1024*1024))
	})
}

func TestLoadCSV(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	content := `Nome;Idade;Cidade
Joao;30;Aracaju
Maria;25;Lisboa`

	tmpFile, err := os.CreateTemp("", "teste_*.csv")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpName := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpName)

	columns, _, err := LoadCSV(logger, tmpName)
	if err != nil {
		t.Fatalf("Erro inesperado ao ler CSV: %v", err)
	}

	t.Run("Deve ter 3 colunas", func(t *testing.T) {
		if len(columns) != 3 {
			t.Errorf("Esperava 3 colunas, recebeu %d", len(columns))
		}
	})

	t.Run("Primeira coluna deve ser Nome com 2 valores", func(t *testing.T) {
		col := columns[0]
		if col.Name != "Nome" {
			t.Errorf("Header errado. Esperava Nome, veio %s", col.Name)
		}
		if len(col.Values) != 2 {
			t.Errorf("Tamanho errado. Esperava 2 valores, veio %d", len(col.Values))
		}
		if col.Values[0] != "Joao" {
			t.Errorf("Valor errado. Esperava Joao, veio %s", col.Values[0])
		}
	})
}

func TestParseData(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	content := `Nome;Idade;Cidade
Joao;30;Aracaju
Olá, João! © 2024. É muito importante testar acentuação: Avó, Pão, Ações, Frequência.;25;Lisboa`

	win1252Bytes := toWindows1252(content)
	reader := bytes.NewReader(win1252Bytes)
	columns, err := ParseData(logger, reader)

	if err != nil {
		t.Fatalf("Erro inesperado ao ler CSV: %v", err)
	}

	t.Run("Deve ter 3 colunas", func(t *testing.T) {
		if len(columns) != 3 {
			t.Errorf("Esperava 3 colunas, recebeu %d", len(columns))
		}
	})

	t.Run("Primeira coluna deve ser Nome com 2 valores", func(t *testing.T) {
		col := columns[0]
		if col.Name != "Nome" {
			t.Errorf("Header errado. Esperava Nome, veio %s", col.Name)
		}
		if len(col.Values) != 2 {
			t.Errorf("Tamanho errado. Esperava 2 valores, veio %d", len(col.Values))
		}
		if col.Values[0] != "Joao" {
			t.Errorf("Valor errado. Esperava Joao, veio %s", col.Values[0])
		}
	})

	t.Run("Deve detectar o Windows-1252 e coverter para utf-8 corretamente", func(t *testing.T) {
		got := columns[0]
		if got.Values[1] != "Olá, João! © 2024. É muito importante testar acentuação: Avó, Pão, Ações, Frequência." {
			t.Errorf("Falha no Windows1252. O Sniffer não converteu. Esperado: %s, Recebido: %s", "Olá, João! © 2024", got.Values[1])
		}
	})
}

func TestParseDataAsync(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	csvContent := "nome;idade\nOlá, João! © 2024;30"
	win1252Bytes := toWindows1252(csvContent)
	reader := bytes.NewReader(win1252Bytes)
	headers, dataChan, err := ParseDataAsync(context.Background(), logger, reader)

	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}

	if headers[0] != "nome" {
		t.Errorf("Esperado header 'nome', recebido '%s'", headers[0])
	}

	row1, ok := <-dataChan
	if !ok {
		t.Fatal("Canal fechou antes de entregar os dados")
	}
	if row1.Row[1] != "30" {
		t.Errorf("Esperado dado '30', recebido '%s'", row1.Row[1])
	}

	t.Run("Deve detectar o Windows-1252 e coverter para utf-8 corretamente", func(t *testing.T) {
		if row1.Row[0] != "Olá, João! © 2024" {
			t.Errorf("Falha no Windows1252. O Sniffer não converteu. Esperado: %s, Recebido: %s", "Olá, João! © 2024", row1.Row[0])
		}
	})
}

func TestParseDataAsync_DirtyLines(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	csvContent := `nome;idade
Joao;30
Maria;25;Lisboa
Pedro;40`

	reader := strings.NewReader(csvContent)
	headers, dataChan, err := ParseDataAsync(context.Background(), logger, reader)

	if err != nil {
		t.Fatalf("Erro ao iniciar parser: %v", err)
	}

	if len(headers) != 2 {
		t.Fatalf("Esperava 2 headers, achou %d", len(headers))
	}

	var results []profiler.StreamData
	for item := range dataChan {
		results = append(results, item)
	}

	if len(results) != 3 {
		t.Fatalf("Esperava 3 itens no canal, recebeu %d", len(results))
	}

	if results[0].Err != nil {
		t.Errorf("Linha 1 não deveria ter erro")
	}
	if results[0].Row[0] != "Joao" {
		t.Errorf("Linha 1 deveria ser Joao")
	}

	if results[1].Err == nil {
		t.Errorf("Linha 2 deveria ser um erro (colunas extras), mas veio nil")
	} else {
		if !strings.Contains(results[1].Err.Error(), "wrong number of fields") {
			t.Logf("Aviso: mensagem de erro diferente do esperado: %v", results[1].Err)
		}
	}
	if results[1].Row != nil {
		t.Errorf("Linha 2 com erro não deveria ter Row preenchido")
	}

	if results[2].Err != nil {
		t.Errorf("Linha 3 não deveria ter erro, o parser deveria se recuperar")
	}
	if results[2].Row[0] != "Pedro" {
		t.Errorf("Linha 3 deveria ser Pedro")
	}
}

func toWindows1252(s string) []byte {
	encoder := charmap.Windows1252.NewEncoder()
	b, _, _ := transform.Bytes(encoder, []byte(s))
	return b
}

func TestSmartReader(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	t.Run("Deve detectar e ler UTF-8 corretamente", func(t *testing.T) {
		input := "Olá, João! © 2024"
		reader := bytes.NewBuffer([]byte(input))

		smartReader, err := NewSmartReader(logger, reader)
		if err != nil {
			t.Fatalf("Erro ao criar smart reader: %v", err)
		}

		content, _ := io.ReadAll(smartReader)
		got := string(content)

		if got != input {
			t.Errorf("Falha no UTF-8. Esperado: %s, Recebido: %s", input, got)
		}
	})

	t.Run("Deve detectar o Windows-1252 e coverter para utf-8 corretamente", func(t *testing.T) {
		input := "Olá, João! © 2024"
		win1252Bytes := toWindows1252(input)
		reader := bytes.NewReader(win1252Bytes)

		smartReader, err := NewSmartReader(logger, reader)
		if err != nil {
			t.Fatalf("Erro ao criar smart reader: %v", err)
		}
		content, _ := io.ReadAll(smartReader)
		got := string(content)
		if got != input {
			t.Errorf("Falha no Windows1252. O Sniffer não converteu. Esperado: %s, Recebido: %s", input, got)
		}
	})
}

func TestSniffer(t *testing.T) {
	csvContent := "nome,idade\nJoao,30"
	readerCSV := bufio.NewReaderSize(strings.NewReader(csvContent), 1024*1024)
	isJSON, _ := sniffJSON(readerCSV)
	if isJSON {
		t.Error("Detectou CSV como JSON incorretamente")
	}

	jsonContent := `{"nome": "Joao", "idade": 30}`
	readerJSON := bufio.NewReaderSize(strings.NewReader(jsonContent), 1024*1024)
	isJSON, _ = sniffJSON(readerJSON)
	if !isJSON {
		t.Error("Falhou ao detectar JSONL válido")
	}

	jsonSpace := `   
      {"nome": "Maria"}`
	readerSpace := bufio.NewReaderSize(strings.NewReader(jsonSpace), 1024*1024)
	isJSON, _ = sniffJSON(readerSpace)
	if !isJSON {
		t.Error("Falhou ao detectar JSONL com espaços iniciais")
	}
}

func TestParseDataAsync_JSONL(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	jsonContent := `{"time":"2023-01-01", "level":"INFO", "msg":"Teste 1"}
{"msg":"Teste 2", "level":"WARN", "time":"2023-01-02"}
{"msg":"Teste 3", "level":"ERROR", "time":"2023-01-03", "extra":"ignorado"}
`

	reader := strings.NewReader(jsonContent)

	headers, dataChan, err := ParseDataAsync(context.Background(), logger, reader)

	if err != nil {
		t.Fatalf("Erro ao iniciar parser: %v", err)
	}

	expectedHeaders := []string{"level", "msg", "time"}
	if len(headers) != 3 {
		t.Fatalf("Esperava 3 headers, recebeu %d: %v", len(headers), headers)
	}
	for i, h := range headers {
		if h != expectedHeaders[i] {
			t.Errorf("Header na posição %d incorreto. Esperado %s, veio %s", i, expectedHeaders[i], h)
		}
	}

	var rows []profiler.StreamData
	for row := range dataChan {
		rows = append(rows, row)
	}

	if len(rows) != 3 {
		t.Errorf("Esperava 3 linhas de dados, recebeu %d", len(rows))
	}

	row2 := rows[1].Row
	if row2[0] != "WARN" { // level
		t.Errorf("Mapeamento incorreto. Coluna 0 (level) deveria ser WARN, foi %s", row2[0])
	}
	if row2[1] != "Teste 2" { // msg
		t.Errorf("Mapeamento incorreto. Coluna 1 (msg) deveria ser Teste 2, foi %s", row2[1])
	}
}

func TestParseDataAsync_ShouldDetectSeparators(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	testCases := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name:     "Separado por Ponto e Vírgula (Padrão Excel BR)",
			content:  "nome;idade\nJoao;30",
			expected: 2,
		},
		{
			name:     "Separado por Vírgula (Padrão US)",
			content:  "nome,idade\nJoao,30",
			expected: 2,
		},
		{
			name:     "Separado por Pipe",
			content:  "nome|idade\nJoao|30",
			expected: 2,
		},
		{
			name:     "Separado por Tab",
			content:  "nome\tidade\nJoao\t30",
			expected: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.content)

			headers, _, err := ParseDataAsync(context.Background(), logger, reader)

			if err != nil {
				t.Fatalf("Erro inesperado: %v", err)
			}

			if len(headers) != tc.expected {
				t.Errorf("Falha na detecção. Esperado %d colunas, mas detectou %d. Headers: %v",
					tc.expected, len(headers), headers)
			}
		})
	}
}

func TestParseData_ShouldDetectSeparators(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	testCases := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name:     "Separado por Ponto e Vírgula (Padrão Excel BR)",
			content:  "nome;idade\nJoao;30",
			expected: 2,
		},
		{
			name:     "Separado por Vírgula (Padrão US)",
			content:  "nome,idade\nJoao,30",
			expected: 2,
		},
		{
			name:     "Separado por Pipe",
			content:  "nome|idade\nJoao|30",
			expected: 2,
		},
		{
			name:     "Separado por Tab",
			content:  "nome\tidade\nJoao\t30",
			expected: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.content)

			headers, err := ParseData(logger, reader)

			if err != nil {
				t.Fatalf("Erro inesperado: %v", err)
			}

			if len(headers) != tc.expected {
				t.Errorf("Falha na detecção. Esperado %d colunas, mas detectou %d. Headers: %v",
					tc.expected, len(headers), headers)
			}
		})
	}
}
