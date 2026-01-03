package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/JGustavoCN/dataprofiler/frontend"
	"github.com/JGustavoCN/dataprofiler/internal/infra"
	"github.com/JGustavoCN/dataprofiler/internal/infra/web"
	"github.com/JGustavoCN/dataprofiler/internal/profiler"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)
	sseBroker := web.NewBroker()

	go func() {
		slog.Info("🔧 Servidor Debug/Pprof iniciado", "addr", "localhost:6060")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			slog.Error("Falha no servidor Pprof", "error", err)
		}
	}()

	slog.Info(
		"Iniciando servidor DataProfiler",
		"port", 8080,
		"env", "production",
		"version", "v1.0.0",
	)

	mux := http.NewServeMux()
	mux.Handle("/events", sseBroker)
	mux.HandleFunc("/api/upload", func(w http.ResponseWriter, r *http.Request) {
		uploadHandlerStreaming(w, r, sseBroker)
	})
	mux.HandleFunc("/api/uploadDeprecated", uploadHandlerDeprecated)

	handlerComCORS := CORSMiddleware(mux)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handlerComCORS,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
		IdleTimeout:  600 * time.Second,
	}
	go func() {

		time.Sleep(1 * time.Second)

		slog.Info("Abrindo navegador automaticamente...")
		openBrowser("http://localhost:8080")
	}()
	assets, err := frontend.GetFileSystem()
	if err != nil {
		log.Fatalf("Falha ao carregar frontend embutido: %v", err)
	}

	fileServer := http.FileServer(assets)

	mux.Handle("/", spaHandler(assets, fileServer))
	go func() {
		slog.Info("Servidor pronto e escutando", "addr", ":8080")

		if err := srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				slog.Error("Servidor caiu", "error", err)
				os.Exit(1)
			}
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	<-stopChan
	slog.Warn("Sinal de desligamento recebido! Iniciando Graceful Shutdown...")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Erro forçando desligamento do servidor", "error", err)
	} else {
		slog.Info("Servidor desligado com sucesso (Gracefully)")
	}

}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":

		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("plataforma não suportada")
	}

	if err != nil {
		slog.Error("Não foi possível abrir o navegador automaticamente", "error", err)
	}
}

func spaHandler(fsys http.FileSystem, fileServer http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Path

		f, err := fsys.Open(path)

		if os.IsNotExist(err) {

			r.URL.Path = "/"
			fileServer.ServeHTTP(w, r)
			return
		} else if err != nil {

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f.Close()
		fileServer.ServeHTTP(w, r)
	})
}

func uploadHandlerStreaming(w http.ResponseWriter, r *http.Request, broker *web.Broker) {
	start := time.Now()
	requestID := start.UnixNano()

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

	onProgress := func(percentage float64, bytesRead int64) {
		msg := fmt.Sprintf(`{"status": "streaming", "progress": %.1f, "bytes": %d}`, percentage, bytesRead)
		broker.Broadcast(msg)
	}

	progressFile := infra.NewProgressReader(file, handler.Size, onProgress)

	log.Info("Iniciando processamento com rastreamento real",
		"filename", handler.Filename,
		"size_bytes", handler.Size,
	)

	headers, dataChan, err := infra.ParseDataAsync(ctx, log, progressFile)

	if err != nil {
		log.Error("Erro crítico no parser", "error", err)
		http.Error(w, "Erro ao ler", http.StatusInternalServerError)
		return
	}

	result := profiler.ProfileAsync(log, headers, dataChan, handler.Filename)
	broker.Broadcast(`{"status": "finishing", "progress": 100}`)

	duration := time.Since(start)
	log.Info("Sucesso",
		"filename", handler.Filename,
		"duration_ms", duration.Milliseconds(),
		"duration_human", duration.String(),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Error("Erro ao codificar JSON de resposta", "error", err)
	}
	broker.Broadcast(`{"status": "done", "progress": 100}`)
}

func uploadHandlerDeprecated(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := start.UnixNano()

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

		duration := time.Since(start)

		log.Info("Sucesso (Deprecated)",
			"filename", handler.Filename,
			"duration_ms", duration.Milliseconds(),
			"duration_human", duration.String(),
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
			"duration_elapsed", time.Since(start).String(),
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
