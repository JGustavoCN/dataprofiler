package profiler

type DataType string

const (
	// Tipos Primitivos
	TypeEmpty   DataType = "EMPTY"
	TypeInteger DataType = "INTEGER"
	TypeFloat   DataType = "FLOAT"
	TypeBoolean DataType = "BOOLEAN"
	TypeString  DataType = "STRING"

	// Documentos & Chaves
	TypeFiscalKey44 DataType = "FISCAL_KEY_44"
	TypeCNPJ        DataType = "CNPJ"
	TypeCPF         DataType = "CPF"
	TypePlaca       DataType = "LICENSE_PLATE"
	TypeNCM         DataType = "NCM"
	TypeRNTRC       DataType = "RNTRC"
	TypeEAN         DataType = "EAN_PRODUCT"

	// Logística & Contatos
	TypeContainer DataType = "CONTAINER_ID"
	TypeCEP       DataType = "CEP"
	TypeMobile    DataType = "MOBILE_PHONE"
	TypeEmail     DataType = "EMAIL"

	// Datas
	TypeDate        DataType = "DATE"
	TypeDateCompact DataType = "DATE_COMPACT"
)

func (d DataType) String() string {
	return string(d)
}
