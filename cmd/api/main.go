package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/JGustavoCN/dataprofiler/internal/infra"
	"github.com/JGustavoCN/dataprofiler/internal/profiler"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	slog.Info(
		"Iniciando servidor DataProfiler", 
        "port", 8080, 
        "env", "production",
        "version", "v1.0.0",
	)

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
	slog.Info("Servidor pronto e escutando", "addr", ":8080")

	if err := srv.ListenAndServe(); err != nil {
		slog.Error("Servidor caiu", "error", err)
		os.Exit(1)
	}
}

func uploadHandlerStreaming(w http.ResponseWriter, r *http.Request) {

	requestID := time.Now().UnixNano()

	log := slog.With(
		"req_id", requestID,
		"method", r.Method,
		"path", r.URL.Path,
	)

	log.Info("Nova requisição de upload recebida")
	defer log.Info("Finalizando requisição de upload")

	if r.Method != http.MethodPost {
		log.Warn("Tentativa de método inválido")
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
	defer cancel()

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Error("Erro parse form", "error", err)
		http.Error(w, "Erro", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Error("Falha ao recuperar arquivo do form", "error", err)
		http.Error(w, "Erro ao recuperar arquivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Info("Iniciando processamento de arquivo",
        "filename", handler.Filename,
        "size_bytes", handler.Size,
    )
	headers, dataChan, err := infra.ParseDataAsync(log, file)
	
	if err != nil {
		log.Error("Erro crítico no parser", "error", err)
		http.Error(w, "Erro ao ler", http.StatusInternalServerError)
		return
	}
	type processingResult struct {
		data interface{}
	}
	done := make(chan processingResult)

	go func() {
		log.Info("Iniciando profile async em background...")
		res := profiler.ProfileAsync(log, headers, dataChan, handler.Filename)
		done <- processingResult{data: res}
		close(done)
	}()

	select {
	case res := <-done:
		log.Info("Sucesso", 
            "filename", handler.Filename,
            "duration_ms", "TODO: Medir tempo",
        )
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res.data); err != nil {
			log.Error("Erro ao codificar JSON de resposta", "error", err)
		}

	case <-ctx.Done():
		log.Warn("Timeout no processamento", 
            "filename", handler.Filename, 
            "timeout_limit", "10m",
        )
		http.Error(w, "Timeout no processamento", http.StatusGatewayTimeout)
		return
	}
}

func uploadHandlerDeprecated(w http.ResponseWriter, r *http.Request) {

	requestID := time.Now().UnixNano()

	log := slog.With(
		"req_id", requestID,
		"method", r.Method,
		"path", r.URL.Path,
		"handler", "deprecated",
	)

	log.Info("Nova requisição recebida (Deprecated)")
	defer log.Info("Finalizando requisição (Deprecated)")

	if r.Method != http.MethodPost {
		log.Warn("Tentativa de método inválido")
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
	defer cancel()

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Error("Erro parse form", "error", err)
		http.Error(w, "Erro", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Error("Falha ao recuperar arquivo do form", "error", err)
		http.Error(w, "Erro ao recuperar arquivo", http.StatusBadRequest)
		return
	}
	defer file.Close()
	log.Info("Iniciando processamento de arquivo",
        "filename", handler.Filename,
        "size_bytes", handler.Size,
    )

	type processingResult struct {
		data interface{}
		err  error
	}

	done := make(chan processingResult)

	go func() {

		log.Info("Iniciando Parse Síncrono")

		columns, err := infra.ParseData(log, file)
		if err != nil {
			done <- processingResult{err: err}
			return
		}
		log.Info("Parse finalizado. Iniciando Profile Síncrono")
		result := profiler.Profile(log, columns, handler.Filename)
		done <- processingResult{data: result}
		log.Info("Profile finalizado")
	}()

	select {
	case res := <-done:
		if res.err != nil {
			log.Error("Erro ao processar", "error", res.err)
			http.Error(w, "Erro ao processar: "+res.err.Error(), http.StatusInternalServerError)
			return
		}
		
		log.Info("Sucesso (Deprecated)", 
            "filename", handler.Filename,
            "duration_ms", "TODO: Medir tempo", 
        )
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res.data); err != nil {
			log.Error("Erro ao gerar JSON", "error", err)
			http.Error(w, "Erro ao gerar JSON", http.StatusInternalServerError)
			return
		}

	case <-ctx.Done():
		log.Warn("Timeout no processamento", 
            "filename", handler.Filename, 
            "timeout_limit", "10m",
        )
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
