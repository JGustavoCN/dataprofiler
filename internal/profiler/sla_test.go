package profiler

import (
	"strings"
	"testing"
)

func TestCalculateSLA(t *testing.T) {

	tests := []struct {
		name             string
		blankRatio       float64
		consistencyRatio float64
		dtype            DataType
		expectedScore    QualityScore
		expectedReason   string
	}{

		{
			name:             "CPF Perfeito",
			blankRatio:       0.0,
			consistencyRatio: 1.0,
			dtype:            TypeCPF,
			expectedScore:    SlaGood,
			expectedReason:   "Dados saudáveis",
		},
		{
			name:             "CPF com 0.5% vazio (Tolerância Zero)",
			blankRatio:       0.005,
			consistencyRatio: 1.0,
			dtype:            TypeCPF,
			expectedScore:    SlaWarning,
			expectedReason:   "Volume elevado de vazios",
		},
		{
			name:             "CPF Sujo (99.5% puro)",
			blankRatio:       0.0,
			consistencyRatio: 0.995,
			dtype:            TypeCPF,
			expectedScore:    SlaWarning,
			expectedReason:   "Indícios de sujeira",
		},

		{
			name:             "Inteiro Padrão (1% vazio)",
			blankRatio:       0.01,
			consistencyRatio: 1.0,
			dtype:            TypeInteger,
			expectedScore:    SlaGood,
			expectedReason:   "Dados saudáveis",
		},
		{
			name:             "Inteiro Crítico (11% vazio)",
			blankRatio:       0.11,
			consistencyRatio: 1.0,
			dtype:            TypeInteger,
			expectedScore:    SlaCritical,
			expectedReason:   "Volume crítico",
		},
		{
			name:             "Inteiro Muito Sujo (85% puro)",
			blankRatio:       0.0,
			consistencyRatio: 0.85,
			dtype:            TypeInteger,
			expectedScore:    SlaCritical,
			expectedReason:   "Alta poluição",
		},

		{
			name:             "String Flexível (15% vazio)",
			blankRatio:       0.15,
			consistencyRatio: 0.0,
			dtype:            TypeString,
			expectedScore:    SlaGood,
			expectedReason:   "Dados saudáveis",
		},
		{
			name:             "String Warning (30% vazio)",
			blankRatio:       0.30,
			consistencyRatio: 1.0,
			dtype:            TypeString,
			expectedScore:    SlaWarning,
			expectedReason:   "Volume elevado",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotScore, gotReason := CalculateSLA(tt.blankRatio, tt.consistencyRatio, tt.dtype)

			if gotScore != tt.expectedScore {
				t.Errorf("Score incorreto para %s: esperado %s, recebeu %s", tt.name, tt.expectedScore, gotScore)
			}

			if !strings.Contains(gotReason, tt.expectedReason) {
				t.Errorf("Motivo incorreto para %s: esperava conter '%s', recebeu '%s'", tt.name, tt.expectedReason, gotReason)
			}
		})
	}
}
