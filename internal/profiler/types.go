package profiler

type DataType string

const (
	TypeEmpty   DataType = "EMPTY"
	TypeInteger DataType = "INTEGER"
	TypeFloat   DataType = "FLOAT"
	TypeBoolean DataType = "BOOLEAN"
	TypeString  DataType = "STRING"

	TypeFiscalKey44 DataType = "FISCAL_KEY_44"
	TypeCNPJ        DataType = "CNPJ"
	TypeCPF         DataType = "CPF"
	TypePlaca       DataType = "LICENSE_PLATE"
	TypeNCM         DataType = "NCM"
	TypeRNTRC       DataType = "RNTRC"
	TypeEAN         DataType = "EAN_PRODUCT"

	TypeContainer DataType = "CONTAINER_ID"
	TypeCEP       DataType = "CEP"
	TypeMobile    DataType = "MOBILE_PHONE"
	TypeEmail     DataType = "EMAIL"

	TypeDate        DataType = "DATE"
	TypeDateCompact DataType = "DATE_COMPACT"
)

func (d DataType) String() string {
	return string(d)
}

type DataSensitivity string

const (
	SensitivityPublic DataSensitivity = "PUBLIC"

	SensitivityInternal DataSensitivity = "INTERNAL"

	SensitivityConfidential DataSensitivity = "CONFIDENTIAL"
)

func (ds DataSensitivity) String() string {
	return string(ds)
}

func ClassifySensitivity(t DataType) (DataSensitivity, string) {
	switch t {
	case TypeCPF, TypeCNPJ:
		return SensitivityConfidential, "Identificação Pessoal/Empresarial (PII)"
	case TypeEmail, TypeMobile:
		return SensitivityConfidential, "Dado de Contato Pessoal (PII)"
	case TypeFiscalKey44:
		return SensitivityConfidential, "Sigilo Fiscal (NFe/CTe)"

	case TypePlaca, TypeRNTRC, TypeContainer:
		return SensitivityInternal, "Rastreabilidade Logística (Segurança Operacional)"
	case TypeEAN, TypeNCM:
		return SensitivityInternal, "Inteligência Comercial/Produto"

	default:
		return SensitivityPublic, "Dado Geral / Não Classificado"
	}
}
