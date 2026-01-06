# ==========================================
# Configurações do Projeto
# ==========================================
BINARY_NAME=dataprofiler.exe
DOCKER_IMAGE=dataprofiler:latest
FRONTEND_DIR=frontend
CMD_DIR=cmd/api

# .PHONY diz ao Make que esses não são arquivos reais
.PHONY: setup install-tools run run-front build build-all build-windows test test-race test-fuzz fmt \
        frontend-install frontend-build docker-build docker-run benchmark clean profile-heap \
		docs-install docs-serve docs-build release

# ==========================================
# 🚀 Workflow Diário (Daily Driver)
# ==========================================

# Prepara a máquina (Instala deps do Go, do React e ferramentas extras como rsrc)
setup: install-tools frontend-install docs-install
	go mod tidy

# Roda o Backend
run:
	go run $(CMD_DIR)/main.go

# Roda o Frontend
run-front:
	cd $(FRONTEND_DIR) && npm run dev

# ==========================================
# 📚 Documentação (MkDocs)
# ==========================================

# Instala o MkDocs e o tema Material via Python
docs-install:
	@echo "📚 Instalando dependencias de documentacao..."
	pip install mkdocs mkdocs-material

# Roda o servidor local de documentação (Hot Reload)
# Usa 'python -m' para evitar problemas de PATH no Windows
docs-serve:
	@echo "📖 Iniciando servidor de documentacao em http://127.0.0.1:8000"
	python -m mkdocs serve

# Gera o site estático na pasta /site (para deploy)
docs-build:
	@echo "🔨 Compilando site estatico..."
	python -m mkdocs build
	@echo "✅ Documentacao gerada na pasta 'site/'"

# Publica no GitHub Pages
docs-deploy:
	@echo "🚀 Publicando documentacao no GitHub Pages..."
	python -m mkdocs gh-deploy --force
	@echo "✅ Documentacao publicada! Acesse em: https://jgustavocn.github.io/dataprofiler/"

# ==========================================
# 🏗️ Build & Distribuição
# ==========================================

# Build simples (Linux/Mac ou dev rápido)
build:
	go build -o $(BINARY_NAME) ./$(CMD_DIR)

# Build Completo (Front + Back)
build-all: frontend-build build

# Build Profissional Windows (Com Ícone e Otimizado)
build-windows: frontend-install frontend-build
	@echo "🎨 Gerando ícone (rsrc)..."
	rsrc -ico app.ico -o $(CMD_DIR)/rsrc.syso
	@echo "🔨 Compilando binário Windows..."
	go build -ldflags "-s -w -H=windowsgui" -o $(BINARY_NAME) ./$(CMD_DIR)
	@echo "✅ Build concluído: ./$(BINARY_NAME)"

# ==========================================
# 🧪 Qualidade & Testes
# ==========================================
test:
	go test ./...

# Detecta Race Conditions (Essencial para Go Routines)
# Nota: 'CGO_ENABLED=1' funciona no Git Bash e Linux. No Powershell puro falharia.
test-race:
	CGO_ENABLED=1 go test -race ./...

# Teste de Estresse (Fuzzing)
test-fuzz:
	go test ./internal/profiler -fuzz=FuzzInferType -fuzztime=10s

fmt:
	go fmt ./...
	go mod tidy

# ==========================================
# 🛠️ Ferramentas & Infra
# ==========================================

install-tools:
	@echo "🔧 Instalando ferramenta de ícone (rsrc)..."
	go install github.com/akavel/rsrc@latest

frontend-install:
	cd $(FRONTEND_DIR) && npm install

frontend-build:
	cd $(FRONTEND_DIR) && npm run build

# ==========================================
# 🐳 Docker
# ==========================================
docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-run:
	docker compose up app

benchmark:
	@echo "🔥 Iniciando Benchmark (Necessário arquivo large_dataset.csv)..."
	docker compose --profile test up benchmark

# ==========================================
# 🔍 Profiling (Baseado no seu histórico)
# ==========================================
profile-heap:
	@echo "📸 Capturando Heap Profile..."
	curl -o heap.out http://localhost:6060/debug/pprof/heap
	go tool pprof -http=:8081 heap.out

# ==========================================
# 📦 Release (Gera binários para GitHub)
# ==========================================
release: frontend-install frontend-build
	@echo "🚀 Preparando release..."
	-mkdir bin
	
	@echo "🎨 Gerando icone (rsrc)..."
	rsrc -ico app.ico -o $(CMD_DIR)/rsrc.syso
	
	@echo "📦 Compilando para Windows (amd64)..."
	go build -ldflags="-s -w" -o bin/dataprofiler.exe ./$(CMD_DIR)
	
	@echo "🐧 Compilando para Linux (amd64)..."
	set GOOS=linux& set GOARCH=amd64& go build -ldflags="-s -w" -o bin/dataprofiler-linux ./$(CMD_DIR)
	
	@echo "✅ Binarios criados na pasta bin/!"

# ==========================================
# 🧹 Limpeza
# ==========================================
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(CMD_DIR)/rsrc.syso
	rm -f *.out
	docker compose down --remove-orphans