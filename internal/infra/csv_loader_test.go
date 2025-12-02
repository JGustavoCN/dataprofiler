package infra

import (
	"os"
	"strings"
	"testing"
)

func TestLoadCSV(t *testing.T) {
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

	columns, _, err := LoadCSV(tmpName)
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
	content := `Nome;Idade;Cidade
Joao;30;Aracaju
Maria;25;Lisboa`

	reader := strings.NewReader(content)
	columns, err := ParseData(reader)

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
