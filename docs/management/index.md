# Gestão do Projeto e Backlog

Este documento detalha a evolução cronológica do **Data Profiler**, cobrindo desde o core matemático até a entrega da versão Enterprise e visões de futuro.

## 📅 Backlog Cronológico (Milestone 1)

### :material-check-circle: Fase 1: Fundação & Infraestrutura

Fase focada em fazer o sistema funcionar, ser resiliente e processar dados em streaming para evitar OOM (_Out of Memory_).

#### Sprint 0: O Core Lógico (A Matemática)

- [x] **Task 0.1 (InferType):** Detetive de Tipos com Regex (Int, Float, String).
- [x] **Task 0.2 (StatsCalc):** Calculadora estatística (Média, Min, Max).
- [x] **Task 0.3 (AnalyzeColumn):** Analista síncrono para contagem de nulos.

#### Sprint 1: O MVP "Happy Path"

- [x] **Task 1.1:** Leitura de CSV com `ReadAll` (Carregamento total na RAM).
- [x] **Task 1.2:** Servidor HTTP Básico síncrono.
- [x] **Task 1.3:** Frontend MVP com CSS Artesanal.

#### Sprint 2: Resiliência HTTP

- [x] **Task 2.1 (Timeouts):** Configuração de `ReadTimeout` e `WriteTimeout` no servidor.
- [x] **Task 2.2 (CORS):** Middleware para comunicação Frontend (5173) <-> Backend (8080).
- [x] **Task 2.3 (Context):** Timeout lógico para cancelamento de requests longos.

#### Sprint 3: Arquitetura Go Way (Streaming)

- [x] **Task 3.1 (Pipeline):** Substituição por leitura linha-a-linha via `Channels`.
- [x] **Task 3.2 (Accumulator):** Cálculo estatístico sem histórico (Stream).
- [x] **Task 3.3 (Async):** Processamento em background (`goroutines`).
- [x] **Task 3.4 (Sniffer):** Detecção automática de Encoding e Separadores.

---

### :material-check-circle: Fase 2: Robustez & Engenharia

Esta fase transforma o script funcional em um software de engenharia robusta, focado em observabilidade e controle de recursos.

#### Sprint 4: Robustez & Backend Engineering :material-server-network

!!! abstract "Foco: Estabilidade e Observabilidade"

    Preparar o motor para suportar carga pesada e ser auditável.

    - [x] **Task 4.1 (Observabilidade):** Migração para `log/slog` e Baseline de memória com `pprof`.
    - [x] **Task 4.2 (Gestão de Memória):** Implementação de `sync.Pool` para redução de GC.
    - [x] **Task 4.3 (Lifecycle):** Graceful Shutdown e cancelamento via `Context`.
    - [x] **Task 4.4 (Fuzzing):** Testes de estresse com dados aleatórios na inferência.

#### Sprint 5: Lógica de Negócio & Dados :material-database-search

!!! abstract "Foco: Regras de Negócio Logísticas"

    Implementar as regras que geram valor para o cliente (SLA, Sujeira, Idempotência).

    - [x] **Task 5.1 (Dirty Data):** Tratamento de linhas irregulares e suporte a `.jsonl`.
    - [x] **Task 5.2 (SLA Logístico):** Regex para CEP, CNPJ, Placas e Score de Qualidade.
    - [x] **Task 5.3 (Estatística Stream):** Reservoir Sampling (Preview) e Histogramas em streaming.

---

### :material-check-circle: Fase 3: Experiência & Entrega

Foco na usabilidade profissional e empacotamento para distribuição.

#### Sprint 6: Frontend Enterprise (Material UI) :material-monitor-shimmer

!!! abstract "Foco: UX Profissional"

    Substituição do CSS artesanal por componentes de dados robustos.

    - [x] **Task 6.1 (Validação):** Bloqueio de extensões/tamanhos e Drag Zone ativa.
    - [x] **Task 6.2 (Real-Time):** Feedback de progresso via Server Sent Events (SSE).
    - [x] **Task 6.3 (MUI Migration):** Refatoração para Material UI e uso de `DataGrid`.
    - [x] **Task 6.4 (Dashboards):** Gráficos estatísticos e Cards de SLA visual.

#### Sprint 7: Packaging & DevOps :material-package-variant

!!! abstract "Foco: Deploy em Arquivo Único"

    Gerar um artefato final fácil de executar em qualquer ambiente.

    - [x] **Task 7.1 (Single Binary):** Embed do React dentro do binário Go.
    - [x] **Task 7.2 (CLI & Benchmark):** Modo terminal e validação do desafio 10GB/512MB.
    - [x] **Task 7.3 (Docker):** Dockerfile Multi-stage e Makefile de automação.
    - [x] **Task 7.4 (Deploy):** Fazer deploy do docker multi-stage do single-binary no Render.
    - [ ] **Task 7.5 (Documentação):** Fazer a organização e documentação do que foi feito.
    - [ ] **Task 7.6 (Ghost Mode):** Flag `windowsgui` para rodar como serviço oculto.

---

## 🐛 Histórico de Bugs Resolvidos

Registro de problemas técnicos identificados e solucionados durante o desenvolvimento da Milestone.

| ID         | Problema                                                                     | Solução Aplicada                                       | Status                  |
| ---------- | ---------------------------------------------------------------------------- | ------------------------------------------------------ | ----------------------- |
| **Bug-01** | **O "Fantasma" do JSON:** Nome do arquivo antigo persistia após novo upload. | Limpeza de estado `fileName` no reset do componente.   | :material-check-circle: |
| **Bug-02** | **Erro no "Catálogo de Cursos":** Falha ao processar CSV específico.         | Ajuste no `SmartReader` para encoding e delimitadores. | :material-check-circle: |

---

## 🔮 Futuro (Milestones 2 e 3)

Backlog de itens para análise pós-release da versão 1.0.

### 🛠️ Dívida Técnica & Refatoração

- [ ] **Task 8.1 (Config Central):** Centralizar "Magic Numbers" e ENVs em pacote `internal/config`.
- [ ] **Refactor (Multipart):** Migrar de buffer em disco para streaming de rede (apenas se latência for crítica).
- [ ] **DevEx:** Configurar `go fmt` e `go race` no pipeline de CI/CD.

### ✨ Novas Features (Roadmap)

- [ ] **Persistência:** Banco de dados (SQLite/Postgres) para histórico de análises.
- [ ] **Cardinalidade:** Algoritmo _HyperLogLog_ para contagem de únicos em Big Data.
- [ ] **Webhooks:** Notificação passiva para sistemas externos.
- [ ] **Exportação:** Gerar relatórios em PDF/HTML estático.

### 🧪 Pesquisa & Inovação

!!! tip "Ideias de Expansão"

    - **Módulo Educacional:** Exibir fórmulas estatísticas (Desvio Padrão, Variância) ao clicar nos dados, servindo como ferramenta de ensino.
    - **Engenharia Reversa:** Analisar padrões de PII (Dados Sensíveis) da lib _Capital One DataProfiler_.
    - **Contexto Logístico:** Mapear validações específicas de transporte (Tempo Certo).
