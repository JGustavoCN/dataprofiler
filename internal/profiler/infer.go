package profiler

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	// Essenciais
	RegexFiscalKey = regexp.MustCompile(`^\d{44}$`)                          // NFe, CTe, MDFe
	RegexPlaca     = regexp.MustCompile(`^[A-Z]{3}-?[0-9][0-9A-Z][0-9]{2}$`) // ABC1234 ou ABC1C34

	// Documentos Brasileiros
	RegexCPF  = regexp.MustCompile(`^\d{3}\.\d{3}\.\d{3}-\d{2}$`)       // 000.000.000-00
	RegexCNPJ = regexp.MustCompile(`^\d{2}\.\d{3}\.\d{3}/\d{4}-\d{2}$`) // 00.000.000/0000-00

	// Datas (Formatos comuns BR e ISO)
	RegexDateBr  = regexp.MustCompile(`^\d{2}/\d{2}/\d{4}$`) // DD/MM/YYYY
	RegexDateIso = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`) // YYYY-MM-DD

	// Logística Avançada
	RegexContainer = regexp.MustCompile(`^[A-Z]{4}\d{7}$`)                  // Padrão ISO
	Regex8Digits   = regexp.MustCompile(`^\d{8}$`)                          // NCM, RNTRC, CEP sem traço, Data compacta
	Regex11Digits  = regexp.MustCompile(`^\d{11}$`)                         // CPF sem formatação
	RegexCEP       = regexp.MustCompile(`^\d{5}-\d{3}$`)                    // CEP com traço
	RegexMobile    = regexp.MustCompile(`^\(?\d{2}\)?\s?9\d{4}-?\d{4}$`)    // Celular com 9
	RegexEmail     = regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`) // Email
	RegexEAN       = regexp.MustCompile(`^\d{13,14}$`)                      // GTIN/EAN (Produtos)
)

func InferType(value string, headerName string) DataType {
	if value == "" {
		return TypeEmpty
	}
	headerLower := strings.ToLower(headerName)

	if Regex8Digits.MatchString(value) {
		if containsAny(headerLower, "ncm", "fiscal", "classificacao", "sh") {
			return TypeNCM
		}
		if containsAny(headerLower, "rntrc", "antt", "transportador") {
			return TypeRNTRC
		}
		if containsAny(headerLower, "cep", "zip", "postal") {
			return TypeCEP
		}
		if isCompactDate(value) {
			return TypeDateCompact
		}
	}

	if RegexFiscalKey.MatchString(value) {
		return TypeFiscalKey44
	}
	if RegexEmail.MatchString(value) {
		return TypeEmail
	}
	if RegexPlaca.MatchString(value) {
		return TypePlaca
	}
	if RegexContainer.MatchString(value) {
		return TypeContainer
	}
	if RegexCEP.MatchString(value) {
		return TypeCEP
	}
	if RegexMobile.MatchString(value) {
		return TypeMobile
	}

	if RegexDateBr.MatchString(value) || RegexDateIso.MatchString(value) {
		return TypeDate
	}

	if RegexEAN.MatchString(value) {
		if containsAny(headerLower, "ean", "gtin", "barras", "item", "produto", "sku") {
			return TypeEAN
		}

		if containsAny(headerLower, "cnpj", "fornecedor", "empresa", "transportadora") {
			return TypeCNPJ
		}
	}

	if RegexCNPJ.MatchString(value) {
		return TypeCNPJ
	}

	if Regex11Digits.MatchString(value) {
		if containsAny(headerLower, "cpf", "cliente", "consumidor", "pessoa", "colaborador", "funcionario", "funcionário", "usuario", "usuário", "rg", "identidade", "documento") {
			return TypeCPF
		}
	}
	if RegexCPF.MatchString(value) {
		return TypeCPF
	}

	if isInt(value) {
		return TypeInteger
	}
	if isFloat(value) {
		return TypeFloat
	}
	if isBool(value) {
		return TypeBoolean
	}

	return TypeString
}

func containsAny(text string, keywords ...string) bool {
	for _, k := range keywords {
		if strings.Contains(text, k) {
			return true
		}
	}
	return false
}

func isCompactDate(value string) bool {
	_, err := time.Parse("20060102", value)
	if err == nil {
		return true
	}
	_, err = time.Parse("02012006", value)
	return err == nil
}

func isInt(value string) bool {
	_, err := strconv.Atoi(value)

	return err == nil
}

func isBool(value string) bool {
	lower := strings.ToLower(value)
	return lower == "true" || lower == "false" || lower == "s" || lower == "n"
}

func isFloat(value string) bool {
	value = strings.Replace(value, ",", ".", 1)
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}
