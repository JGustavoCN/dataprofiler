package profiler

import "fmt"

type QualityScore string

const (
	SlaGood     QualityScore = "GOOD"
	SlaWarning  QualityScore = "WARNING"
	SlaCritical QualityScore = "CRITICAL"
)

type SeverityLevel int

const (
	SeverityHigh   SeverityLevel = 3
	SeverityMedium SeverityLevel = 2
	SeverityLow    SeverityLevel = 1
)

func (q QualityScore) String() string {
	return string(q)
}

func CalculateSLA(blankRatio float64, consistencyRatio float64, dtype DataType) (QualityScore, string) {

	severity := getSeverity(dtype)

	scoreCompleteness := evaluateCompleteness(blankRatio, severity)

	if scoreCompleteness == SlaCritical {
		return SlaCritical, fmt.Sprintf("Volume crítico de dados ausentes (%.1f%% vazios)", blankRatio*100)
	}

	scoreConsistency := evaluateConsistency(consistencyRatio, severity)

	if scoreConsistency == SlaCritical {
		return SlaCritical, fmt.Sprintf("Alta poluição de dados: %.1f%% dos valores não são %s", (1-consistencyRatio)*100, dtype)
	}

	if scoreCompleteness == SlaWarning {
		return SlaWarning, fmt.Sprintf("Atenção: Volume elevado de vazios (%.1f%%)", blankRatio*100)
	}
	if scoreConsistency == SlaWarning {
		return SlaWarning, fmt.Sprintf("Atenção: Indícios de sujeira nos dados (%.1f%% inválidos)", (1-consistencyRatio)*100)
	}

	return SlaGood, "Dados saudáveis"
}

func evaluateCompleteness(ratio float64, severity SeverityLevel) QualityScore {
	var limitGood, limitWarning float64

	switch severity {
	case SeverityHigh:
		limitGood, limitWarning = 0.00, 0.01
	case SeverityMedium:
		limitGood, limitWarning = 0.01, 0.10
	case SeverityLow:
		limitGood, limitWarning = 0.20, 0.50
	}

	if ratio <= limitGood {
		return SlaGood
	}
	if ratio <= limitWarning {
		return SlaWarning
	}
	return SlaCritical
}

func evaluateConsistency(ratio float64, severity SeverityLevel) QualityScore {
	if ratio == 0 {
		return SlaGood
	}

	var limitCritical, limitWarning float64

	switch severity {
	case SeverityHigh:
		limitCritical, limitWarning = 0.99, 0.999
	case SeverityMedium:
		limitCritical, limitWarning = 0.90, 0.98
	case SeverityLow:
		return SlaGood
	}

	if ratio < limitCritical {
		return SlaCritical
	}
	if ratio < limitWarning {
		return SlaWarning
	}
	return SlaGood
}

func getSeverity(t DataType) SeverityLevel {
	switch t {
	case TypeCPF, TypeCNPJ, TypeFiscalKey44, TypePlaca, TypeRNTRC, TypeContainer:
		return SeverityHigh
	case TypeInteger, TypeFloat, TypeDate, TypeDateCompact, TypeEmail, TypeEAN, TypeMobile, TypeCEP:
		return SeverityMedium
	default:
		return SeverityLow
	}
}
