package profiler

import "testing"

func TestInferType(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		header   string
		expected DataType // Mudou de string para DataType
	}{
		// --- 1. PRIMITIVOS BÁSICOS ---
		{"Vazio", "", "qualquer_coisa", TypeEmpty},
		{"Inteiro Simples", "12345", "id", TypeInteger},
		{"Inteiro Negativo", "-50", "temp", TypeInteger},
		{"Float Ponto", "12.50", "valor", TypeFloat},
		{"Float Virgula (BR)", "12,50", "preco", TypeFloat},
		{"Boolean True", "true", "ativo", TypeBoolean},
		{"Boolean S (Sim)", "s", "flag", TypeBoolean},
		{"String Comum", "Garrafa de Agua", "desc", TypeString},

		// --- 2. PADRÕES "FORTES" ---
		{"Chave NFe", "35230912345678000190550010000000011000000000", "obs", TypeFiscalKey44},
		{"Email", "contato@empresa.com.br", "email_contato", TypeEmail},
		{"Placa Mercosul", "ABC1D23", "placa_veiculo", TypePlaca},
		{"Placa Antiga", "ABC1234", "veiculo", TypePlaca},
		{"Container ISO", "MSKU1234567", "container", TypeContainer},
		{"Celular BR", "11 91234-5678", "tel", TypeMobile},
		{"CEP com Traço", "01310-100", "end_cep", TypeCEP},
		{"CNPJ Formatado", "12.345.678/0001-90", "doc", TypeCNPJ},
		{"CPF Formatado", "123.456.789-00", "doc", TypeCPF},

		// --- 3. DATAS ---
		{"Data BR", "25/12/2023", "data", TypeDate},
		{"Data ISO", "2023-12-25", "dt_nasc", TypeDate},

		// --- 4. ZONA DE CONFLITO: 8 DÍGITOS ---
		{"8 Digitos -> NCM (Header)", "12345678", "ncm_produto", TypeNCM},
		{"8 Digitos -> RNTRC (Header)", "12345678", "rntrc_motorista", TypeRNTRC},
		{"8 Digitos -> CEP (Header)", "12345678", "cep_origem", TypeCEP},
		{"8 Digitos -> Data Compacta", "20231225", "dt_emissao", TypeDateCompact},
		{"8 Digitos -> Inteiro (Sem Contexto)", "99999999", "codigo_x", TypeInteger},

		// --- 5. ZONA DE CONFLITO: 11 DÍGITOS ---
		{"11 Digitos -> CPF (Header)", "12345678901", "cpf_motorista", TypeCPF},
		{"11 Digitos -> Inteiro (Sem Contexto)", "12345678901", "id_transacao", TypeInteger},

		// --- 6. ZONA DE CONFLITO: EAN vs CNPJ ---
		{"EAN (Header Produto)", "7891234567890", "ean_produto", TypeEAN},
		{"EAN (Header SKU)", "7891234567890", "sku", TypeEAN},
		{"CNPJ Limpo (Header Empresa)", "12345678000190", "cnpj_emitente", TypeCNPJ},

		// TESTE CRÍTICO: CNPJ/EAN sem header deve cair para INTEGER para não estragar cálculo
		{"CNPJ Limpo (Padrão Regex)", "12345678000190", "coluna_x", TypeInteger},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InferType(tt.value, tt.header)

			if got != tt.expected {
				t.Errorf("InferType(%q, %q) = %v; esperava %v", tt.value, tt.header, got, tt.expected)
			}
		})
	}
}

func FuzzInferType(f *testing.F) {
	f.Add("123", "id")
	f.Add("abc", "nome")
	f.Add("2023-01-01", "data")
	f.Add("", "")

	f.Fuzz(func(t *testing.T, value string, header string) {
		_ = InferType(value, header)
	})
}
