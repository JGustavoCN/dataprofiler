package profiler

import "testing"

func TestInferType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Deve identificar Inteiro", "123", "int"},
		{"Deve identificar Inteiro negativo", "-50", "int"},
		{"Deve identificar Float", "12.50", "float"},
		{"Deve identificar Float com muitos decimais", "0.0009", "float"},
		{"Deve identificar Boolean True", "true", "bool"},
		{"Deve identificar Boolean False", "false", "bool"},
		{"Deve identificar Boolean f = False", "f", "bool"},
		{"Deve identificar Boolean t = True", "t", "bool"},
		{"Deve identificar String comum", "Garrafa de Agua", "string"},
		{"Deve identificar String vazia", "", "string"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InferType(tt.input)
			if got != tt.expected {
				t.Errorf("InferType(%q) = %v, esperava %v", tt.input, got, tt.expected)
			}
		})
	}

}

func FuzzInferType(f *testing.F) {
	f.Add("123")
	f.Add("12.50")
	f.Add("true")
	f.Add("2025-01-01")
	f.Add("Texto Comum")
	f.Add("") 
	f.Fuzz(func(t *testing.T, input string) {
		_ = InferType(input)
	})
}
