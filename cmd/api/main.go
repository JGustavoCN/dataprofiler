package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JGustavoCN/dataprofiler/internal/infra"
	"github.com/JGustavoCN/dataprofiler/internal/profiler"
)

func main() {
	http.HandleFunc("/api/upload", uploadHandler)

	fmt.Println("🚀 Servidor rodando na porta :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erro ao recuperar arquivo", http.StatusBadRequest)
		return
	}
	defer file.Close()
	fmt.Printf("📂 Recebido arquivo: %s\n", handler.Filename)

	fmt.Println("============ Começo do parse")
	columns, err := infra.ParseData(file)
	if err != nil {
		fmt.Println("Erro ao processar CSV", err.Error())
		http.Error(w, "Erro ao processar CSV", http.StatusInternalServerError)
		return
	}
	fmt.Println("============ Terminou o parse")

	fmt.Println("============ Começo do profile")
	result := profiler.Profile(columns, handler.Filename)
	fmt.Println("============ Terminou o profile")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Erro ao gerar JSON", http.StatusInternalServerError)
	}
	fmt.Println("Terminou o envio do json")
}

/**
import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/JGustavoCN/dataprofiler/internal/infra"
	"github.com/JGustavoCN/dataprofiler/internal/profiler"
)

func main() {
	filePath := "produtos_teste.csv"

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
}*/
