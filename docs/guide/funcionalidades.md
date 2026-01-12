# Regras de Negócio e Funcionalidades

O diferencial do DataProfiler é sua capacidade de "entender" o dado, não apenas lê-lo. Abaixo detalhamos os algoritmos de inferência.

## 1. Cálculo de SLA (Nível de Qualidade)

O sistema atribui um selo de qualidade para cada coluna processada. Isso permite que engenheiros de dados decidam rapidamente se aquela coluna pode ser usada em um modelo de Machine Learning ou Dashboard.

### A Lógica Matemática

O cálculo é baseado na **Densidade de Informação**. O sistema contabiliza, em tempo real, quantos valores são considerados "sujos" (Nulos, Vazios, `NA`, `NULL`).

$$\text{Score} = \frac{\text{Total Linhas} - \text{Linhas Sujas}}{\text{Total Linhas}} \times 100$$

### Classificação

| Selo                                                           | Critério              | Interpretação                                                                                                |
| :------------------------------------------------------------- | :-------------------- | :----------------------------------------------------------------------------------------------------------- |
| <span style="color:#D4AF37; font-weight:bold">🥇 Ouro</span>   | **Score ≥ 99%**       | **Alta Confiabilidade.** Dados praticamente íntegros. Seguros para chaves primárias ou métricas financeiras. |
| <span style="color:#C0C0C0; font-weight:bold">🥈 Prata</span>  | **95% ≤ Score < 99%** | **Confiabilidade Média.** Dados úteis para análises de tendência, mas requerem atenção em casos de borda.    |
| <span style="color:#CD7F32; font-weight:bold">🥉 Bronze</span> | **Score < 95%**       | **Baixa Qualidade.** Requer tratamento (imputação de dados) antes do uso. Alto risco de viés.                |

---

## 2. Detecção de Sensibilidade (LGPD/GDPR)

Para garantir conformidade com leis de proteção de dados, o DataProfiler escaneia o conteúdo em busca de **PII (Personally Identifiable Information)**.

O algoritmo funciona em duas camadas:

1. **Análise de Cabeçalho:** Verifica se o nome da coluna sugere dados sensíveis (ex: "cpf_cliente", "email_contato").
2. **Análise de Conteúdo (Regex):** Verifica se os valores batem com padrões conhecidos.

### Padrões Detectados

!!! warning "Atenção"
Se uma coluna for marcada como Sensível, o ícone 🛡️ aparecerá no relatório. Recomenda-se aplicar hashing ou mascaramento nesses dados.

- **CPF (Brasil):** Validação de formato `111.222.333-44` ou `11122233344`.
- **E-mail:** Padrão RFC 5322 (`usuario@dominio.com`).
- **Cartão de Crédito:** Detecção de sequências numéricas compatíveis com PANs (Luhn Algorithm check básico).
- **Telefone:** Padrões globais E.164 e nacionais.

---

## 3. Inferência de Tipos (Polimorfismo)

Como o CSV é um formato sem tipo (tudo é texto), o DataProfiler realiza uma inferência estatística. Ele lê uma amostragem dos dados e tenta "promover" o tipo para o mais específico possível.

**Ordem de Tentativa:**

1. **Integer:** É um número inteiro? (ex: `42`)
2. **Float:** É decimal? (ex: `42.5` ou `42,5`) -> _Suporta ponto e vírgula como decimal._
3. **Boolean:** É lógico? (ex: `true`, `1`, `sim`, `yes`)
4. **Date:** É data? (ex: `2023-01-01`, `01/01/2023`) -> _Suporta ISO8601 e BR._
5. **String:** Se falhar em tudo, é texto.
