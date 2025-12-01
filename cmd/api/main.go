package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/JGustavoCN/dataprofiler/internal/infra"
	"github.com/JGustavoCN/dataprofiler/internal/profiler"
)

func main() {
	filePath := "C:\\Users\\joseg\\Downloads\\catalago_cursos.csv"

	fmt.Println("🚀 Iniciando DataProfiler...")
	fmt.Printf("📂 Lendo arquivo: %s\n", filePath)

	columns, fileName, err := infra.LoadCSV(filePath)
	if err != nil {
		fmt.Printf("❌ Erro ao ler arquivo: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Sucesso! %d colunas carregadas.\n", len(columns))
	fmt.Println("🧠 Iniciando análise dos dados...")

	result := profiler.Profile(columns, fileName)

	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("❌ Erro ao gerar JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("📊 Relatório Final:")
	fmt.Println(string(jsonOutput))
}
