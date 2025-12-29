package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/JGustavoCN/dataprofiler/internal/infra"
	"github.com/JGustavoCN/dataprofiler/internal/profiler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/uploadDeprecated", uploadHandlerDeprecated)
	mux.HandleFunc("/api/upload", uploadHandlerStreaming)
	handlerComCORS := CORSMiddleware(mux)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handlerComCORS,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
		IdleTimeout:  600 * time.Second,
	}
	fmt.Println("🚀 Servidor Blindado rodando na porta :8080")

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func uploadHandlerStreaming(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
	defer cancel()

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erro ao recuperar arquivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Printf("📂 Recebido arquivo: %s (Streaming)\n", handler.Filename)
	fmt.Println("🚩 ============ Começo do parse async")
	headers, dataChan, err := infra.ParseDataAsync(file)
	
	if err != nil {
		fmt.Println("Erro ao iniciar leitura do CSV:", err)
		http.Error(w, "Erro ao ler CSV", http.StatusInternalServerError)
		return
	}
	type processingResult struct {
		data interface{}
	}
	done := make(chan processingResult)

	go func() {
		fmt.Println("🚩 ============ Começo do profile")
		fmt.Println("🔄 Iniciando processamento via Stream...")
		res := profiler.ProfileAsync(headers, dataChan, handler.Filename)
		done <- processingResult{data: res}
		close(done)
	}()

	select {
	case res := <-done:
		fmt.Println("✅ Processamento concluído com sucesso!")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res.data); err != nil {
			fmt.Println("Erro ao codificar JSON:", err)
		}

	case <-ctx.Done():
		fmt.Println("⏱️ Timeout! O processamento demorou demais.")
		http.Error(w, "Timeout no processamento", http.StatusGatewayTimeout)
		return
	}
}

func uploadHandlerDeprecated(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
	defer cancel()

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erro ao recuperar arquivo", http.StatusBadRequest)
		return
	}
	defer file.Close()
	fmt.Printf("📂 Recebido arquivo: %s\n", handler.Filename)

	type processingResult struct {
		data interface{}
		err  error
	}

	done := make(chan processingResult)

	go func() {

		fmt.Println("============ Começo do parse")

		columns, err := infra.ParseData(file)
		if err != nil {
			done <- processingResult{err: err}
			return
		}
		fmt.Println("============ Terminou o parse")

		fmt.Println("============ Começo do profile")
		result := profiler.Profile(columns, handler.Filename)
		done <- processingResult{data: result}
		fmt.Println("============ Terminou o profile")
	}()

	select {
	case res := <-done:
		if res.err != nil {
			fmt.Println("❌ Erro interno no processamento")
			http.Error(w, "Erro ao processar CSV: "+res.err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Preparação e envio do json")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res.data); err != nil {
			http.Error(w, "Erro ao gerar JSON", http.StatusInternalServerError)
			return
		}
		fmt.Println("Terminou o envio do json")

	case <-ctx.Done():
		fmt.Println("⏱️ Timeout Lógico atingido! Cancelando resposta.")
		http.Error(w, "O processamento demorou demais e foi cancelado.", http.StatusGatewayTimeout)
		return
	}

}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, DELETE, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)

	})

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
