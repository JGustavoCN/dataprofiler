package infra

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

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
Olá, João! © 2024;25;Lisboa`


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
		if got.Values[1] != "Olá, João! © 2024" {
			t.Errorf("Falha no Windows1252. O Sniffer não converteu. Esperado: %s, Recebido: %s", "Olá, João! © 2024", got.Values[1])
		}	
	})
}


func TestParseDataAsync(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	csvContent := "nome;idade\nOlá, João! © 2024;30"
	win1252Bytes := toWindows1252(csvContent)
	reader := bytes.NewReader(win1252Bytes)
	headers, dataChan, err := ParseDataAsync(logger, reader)

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
	if row1[1] != "30" {
		t.Errorf("Esperado dado '30', recebido '%s'", row1[1])
	}

	t.Run("Deve detectar o Windows-1252 e coverter para utf-8 corretamente", func(t *testing.T) {
		if  row1[0] != "Olá, João! © 2024" {
			t.Errorf("Falha no Windows1252. O Sniffer não converteu. Esperado: %s, Recebido: %s", "Olá, João! © 2024", row1[0])
		}	
	})
}


func toWindows1252(s string) []byte {
	encoder := charmap.Windows1252.NewEncoder()
	b, _, _ := transform.Bytes(encoder, []byte(s))
	return b
}

func TestSmartReader(t *testing.T){
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

func TestParseDataAsync_ShouldDetectSeparators(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	testCases := []struct {
		name      string
		content   string
		expected  int 
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
			
			headers, _, err := ParseDataAsync(logger, reader)
			
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
		name      string
		content   string
		expected  int 
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