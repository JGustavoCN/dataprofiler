# Guia de Instalação e Execução

O **DataProfiler** foi projetado para ser agnóstico de plataforma. Oferecemos três métodos de execução, variando do "clique e use" (para usuários finais) até a compilação total (para desenvolvedores).

## Método 1: Executável (Recomendado)

A forma mais simples de utilizar o DataProfiler é através do conceito de **Single Binary**. O Backend (Go) e o Frontend (React) foram compilados em um único arquivo. Não é necessário instalar Java, Python ou Node.js.

### Windows

1. Acesse a [Página de Releases](https://github.com/JGustavoCN/dataprofiler/releases) do projeto.
2. Baixe o arquivo `dataprofiler-windows-amd64.exe`.
3. Dê um **duplo clique** no arquivo baixado.
4. O terminal se abrirá e, em seguida, seu navegador padrão abrirá automaticamente em `http://localhost:8080`.

!!! warning "Alerta do Windows Defender"

    Como este é um software open-source não assinado digitalmente (o que custa caro), o Windows pode exibir a tela _"O Windows protegeu o computador"_.

    Isso é um **falso positivo**. Para prosseguir:

    1. Clique em **"Mais informações"**.
    2. Clique no botão **"Executar assim mesmo"**.

### Linux / macOS

1. Acesse a [Página de Releases](https://github.com/JGustavoCN/dataprofiler/releases) do projeto.
2. Baixe o arquivo `dataprofiler-linux-amd64` (ou `darwin` para Mac).
3. Abra o terminal na pasta do download e dê permissão de execução:

   ```bash
   chmod +x dataprofiler-linux-amd64
   ```

4. Execute o programa:

   ```bash
   ./dataprofiler-linux-amd64
   ```

---

## Método 2: Docker (Ambiente Isolado)

Se você prefere não rodar binários diretamente no seu sistema operacional, disponibilizamos uma imagem Docker oficial. Este método garante que o ambiente seja idêntico ao de produção.

**Pré-requisitos:** Docker e Docker Compose instalados.

Crie um arquivo `docker-compose.yml`:

```yaml
title="docker-compose.yml"
version: "3.8"
services:
  app:
    image: jgustavocn/dataprofiler:latest
    ports:
      - "8080:8080"
    volumes:
      - ./uploads:/app/uploads
```

No terminal, execute:

```bash
docker compose up -d
```

> O sistema estará disponível em: [http://localhost:8080](http://localhost:8080)

---

## Método 3: Compilando do Código (Para Desenvolvedores)

Se você deseja contribuir com o código ou testar funcionalidades experimentais, siga os passos de compilação manual.

### Pré-requisitos

- **Go:** Versão 1.22 ou superior.
- **Node.js:** Versão 20 (LTS) ou superior.
- **Make:** (Opcional, mas recomendado).

### Passo a Passo

1. **Clone o repositório:**

```bash
git clone https://github.com/JGustavoCN/dataprofiler.git
cd dataprofiler
```

1. **Instale as dependências (Backend e Frontend):**
   Utilizamos o Makefile para automatizar a instalação das libs do Go e os pacotes npm do React.

```bash
make setup
```

1. **Execute em modo de desenvolvimento:**
   Este comando roda o Backend com Hot-Reload (Air) e o Frontend em modo dev.

```bash
make run
```

---

## 🔧 Troubleshooting (Resolução de Problemas)

### Erro: "Address already in use" (Porta 8080 ocupada)

Se você já tiver outro serviço rodando na porta `8080` (comum em desenvolvedores Java/Tomcat), o DataProfiler não iniciará.

**Solução:** Defina a variável de ambiente `PORT` antes de executar.

### Windows (PowerShell)

```powershell
$env:PORT="9090"; .\dataprofiler.exe
```

### Linux / Mac

```bash
PORT=9090 ./dataprofiler
```

### Erro: Tela Branca ou "Connection Refused"

Certifique-se de que o backend Go está rodando. O Frontend React depende da API para funcionar. Se você estiver rodando via código fonte, garanta que ambos os terminais (Go e Node) estejam ativos.
