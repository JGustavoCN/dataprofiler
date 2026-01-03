package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"
)

type Broker struct {
	mu       sync.Mutex
	clients  map[chan string]bool
	Notifier chan string
}

func NewBroker() *Broker {
	return &Broker{
		clients:  make(map[chan string]bool),
		Notifier: make(chan string),
	}
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	slog.Info("Nova conexão SSE iniciada", "remote_addr", r.RemoteAddr)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		slog.Error("Streaming não suportado pelo cliente", "remote_addr", r.RemoteAddr)
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	messageChan := make(chan string)
	b.mu.Lock()
	b.clients[messageChan] = true
	b.mu.Unlock()

	defer func() {
		b.mu.Lock()
		delete(b.clients, messageChan)
		b.mu.Unlock()
		close(messageChan)
		slog.Info("Cliente SSE desconectado", "remote_addr", r.RemoteAddr)
	}()

	for {
		select {
		case msg := <-messageChan:

			_, err := fmt.Fprintf(w, "data: %s\n\n", msg)
			if err != nil {
				slog.Error("Erro ao escrever dados SSE", "error", err)
				return
			}
			flusher.Flush()
		case <-r.Context().Done():

			return
		}
	}
}

func (b *Broker) Broadcast(msg string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	slog.Debug("Broadcast SSE", "active_clients", len(b.clients), "msg_size", len(msg))

	for clientChan := range b.clients {
		select {
		case clientChan <- msg:
		default:

			slog.Warn("Cliente SSE lento/travado, pulando mensagem")
		}
	}
}
