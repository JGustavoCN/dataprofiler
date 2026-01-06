---
title: Home
hide:
  - navigation
  - toc
---

<div class="hero-section">
  <img src="assets/logo.svg" alt="DataProfiler Logo" class="hero-logo">
  
  <h1>DataProfiler Enterprise</h1>
  <p>
    Análise de Qualidade e Perfilamento de Dados processando <strong>Gigabytes</strong> com consumo mínimo de RAM.
  </p>
  
  <div class="hero-buttons">
    <a href="guia-usuario/instalacao/" class="md-button md-button--primary">
      🚀 Começar Agora
    </a>
    <a href="engenharia/arquitetura-streaming/" class="md-button md-button--secondary">
      Entender a Engenharia
    </a>
  </div>
</div>

<h3 align="center"> O Problema: Big Data vs Hardware Limitado </h3>

<div class="features-grid">

  <div class="feature-card">
    <h3><span>🚀</span> Alta Performance</h3>
    <p>
      Esqueça o erro <code>Out of Memory</code>. Nossa arquitetura lê arquivos maiores que a RAM disponível, utilizando buffers inteligentes e <strong>I/O não bloqueante</strong>. Processa 10GB com apenas 512MB de RAM.
    </p>
  </div>

  <div class="feature-card">
    <h3><span>🛡️</span> SLA Automático</h3>
    <p>
      O sistema classifica automaticamente a qualidade das colunas (<strong>Ouro, Prata, Bronze</strong>) calculando a densidade de informação e consistência em tempo real para tomada de decisão.
    </p>
  </div>

  <div class="feature-card">
    <h3><span>👁️</span> Segurança & LGPD</h3>
    <p>
      Detector de <strong>PII (Dados Pessoais)</strong> integrado. O sistema varre e alerta sobre CPF, E-mails e Cartões de Crédito expostos para garantir conformidade.
    </p>
  </div>

  <div class="feature-card">
    <h3><span>📦</span> Single Binary</h3>
    <p>
      Zero dependências. O Backend (Go) e o Frontend (React) são compilados em um único arquivo executável <code>.exe</code>. Baixou, rodou, usou.
    </p>
  </div>

  <div class="feature-card">
    <h3><span>🧠</span> Inferência Inteligente</h3>
    <p>
      Esqueça o mapeamento manual (`schema`). O algoritmo de <strong>Type Inference</strong> analisa amostras dos dados para detectar automaticamente se a coluna é Inteiro, Decimal, Data ou Texto.
    </p>
  </div>

  <div class="feature-card">
    <h3><span>📊</span> Interface & Estatísticas</h3>
    <p>
      Frontend em <strong>React + Material UI</strong>. Oferece DataGrid com paginação nativa, filtros avançados e cálculo automático de estatísticas (Média, Mediana, Desvio Padrão) em tempo real.
    </p>
  </div>

</div>

<h3 align="center"> A Engenharia por trás do Streaming</h3>

O diferencial do DataProfiler é a arquitetura <strong>Producer-Consumer</strong>. O dado flui através de canais concorrentes sem nunca ser carregado totalmente na memória.

```mermaid
graph LR
    A[Arquivo CSV Massivo] -->|Stream Leitura| B(Go Reader / Buffer);
    B -->|Chunks de Dados| C{Canal de Distribuição};
    C -->|Worker 1| D[Validação de Tipos];
    C -->|Worker 2| E[Regex PII];
    C -->|Worker 3| F[Estatística];
    D & E & F -->|Agregação| G[Relatório JSON];
    G --> H[Dashboard React];

    style B fill:#3f51b5,stroke:#fff,stroke-width:2px,color:#fff
    style H fill:#2196f3,stroke:#fff,stroke-width:2px,color:#fff
```

<div align="center" class="hero-buttons" style="margin-top: 4rem; margin-bottom: 2rem;" markdown>

<h3 align="center">
Pronto para usar?
Não requer Python, Java ou Docker obrigatório.
</h3>

<a href="guia-usuario/instalacao/" class="md-button"> Baixar para Windows (.exe) </a>

<a href="decisoes/001-escolha-documentacao/" class="md-button"> Ver Decisões de Arquitetura (ADR) </a>

</div>
