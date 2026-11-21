[ [English](../en/artifacts.md) | [Português](artifacts.md) | [Español](../es/artifacts.md) ]

# Artefatos do Specforce

Para orquestrar agentes de IA sem acumular contexto desnecessário, o Specforce divide o conhecimento do projeto em documentos focados chamados **Artefatos**. Esta página explica cada um deles e as regras de ciclo de vida que decidem o que o agente produz e por que uma etapa pode ser bloqueada. As skills são citadas pelo nome (por exemplo [`spf.spec`](supported-tools.md#invocando-skills)).

## 1. Constituição (`.specforce/docs/`)

A **Constituição** fica na pasta `.specforce/docs/` do seu projeto. Ela contém as regras duradouras do sistema e é criada e mantida com o `spf.constitution`:

- **`principles.md`**: valores centrais, filosofia e axiomas técnicos.
- **`architecture.md`**: fronteiras do sistema, direção das dependências e topologia de persistência.
- **`ui-ux.md`**: direção visual, padrões de interação e regras estéticas (quando aplicável).
- **`security.md`**: autorização, papéis, permissões e proteção de dados.
- **`engineering.md`**: padrões de código, estratégia de testes e diretrizes de refatoração.
- **`governance.md`**: regras de ciclo de vida, responsabilidades e limites da IA.
- **`modules/<module>.md`**: especificações vivas de cada domínio (regras de negócio, casos de uso, superfícies públicas), atualizadas no momento do arquivamento a partir dos deltas das especificações.

Os documentos não são carregados todos de uma vez: o agente lê os que são relevantes para a tarefa em questão.

## 2. O pacote da especificação (`.specforce/specs/<timestamp>-<slug>/`)

Toda mudança começa com uma especificação criada pelo `spf.spec`. A pasta dela pode conter:

- **`spec.yaml`**: metadados gerenciados pela CLI: tipo, tamanho, declaração e selo do delta, auditoria de refinamento, registros de tempo das tarefas e o relatório de QA.
- **`proposal.md`** *(opcional)*: o resultado formalizado de uma sessão de `spf.discovery` (problema, descobertas verificadas, decisões, perguntas em aberto, tipo e tamanho sugeridos). O planejamento o reaproveita em vez de perguntar de novo.
- **`requirements.md`**: o "quê" e o "porquê": histórias de usuário, critérios de aceitação BDD (GIVEN/WHEN/THEN) e restrições de negócio. Sem detalhes técnicos.
- **`design.md`**: o "como": arquitetura, contratos e o inventário de arquivos e componentes. Sem código final.
- **`tasks.md`**: o roadmap atômico de execução (veja [Formato das tarefas](#5-formato-das-tarefas)).
- **`deltas/<module>.md`**: as mudanças propostas para cada especificação viva de módulo afetada (veja [Deltas de módulo](#6-deltas-de-módulo)).

Quais desses arquivos são obrigatórios depende do tamanho da especificação.

## 3. Tamanhos

Toda especificação tem um tamanho, escolhido na criação (`specforce spec init <slug> --size <size>`) e alterável depois com `specforce spec resize`. O tamanho padrão é `medium`.

| Tamanho | Artefatos obrigatórios | Delta de módulo |
|---|---|---|
| `small` | `tasks.md` | Isento |
| `medium` (padrão) | `requirements.md`, `tasks.md` | Obrigatório (ou `delta.impact: none`) |
| `large` | `requirements.md`, `design.md`, `tasks.md` | Obrigatório (ou `delta.impact: none`) |
| `complex` | `requirements.md`, `design.md`, `tasks.md` | Obrigatório (ou `delta.impact: none`) |

- A implementação fica bloqueada até que todos os artefatos exigidos pelo tamanho existam.
- Promover uma especificação (por exemplo de `small` para `large`) marca os novos artefatos como ausentes; rebaixá-la relaxa a validação e mantém os arquivos existentes como referência.
- Em especificações `small`, as tarefas podem omitir o vínculo `**Context:** [US-X]`; a partir de `medium`, toda tarefa deve se vincular a um requisito.
- Roadmaps `large` e `complex` devem incluir uma seção de Black-box Scenarios (veja abaixo).

## 4. Especificações de bug

Uma especificação criada com `--type bug` usa as variantes de bug dos artefatos de planejamento, a partir dos templates `bug-requirements` e `bug-design`:

- **Requisitos (`bug-requirements`)**: descrição do problema, evidências (logs, erros), passos concretos de reprodução, uma análise de causa raiz baseada em evidências e critérios de aceitação de regressão (`[FIX-x]`). Quando a causa raiz não está estabelecida, o agente para e pergunta ou recomenda o `spf.discovery` em vez de adivinhar.
- **Design (`bug-design`)**: os caminhos de código a alterar, a estratégia de regressão e os efeitos colaterais.
- **Roadmap**: a primeira tarefa é um teste de regressão que falha antes da correção (uma tarefa `[RED]`). Nenhuma tarefa de correção pode vir antes dela.

Os arquivos mantêm os nomes de sempre (`requirements.md`, `design.md`, `tasks.md`).

## 5. Formato das tarefas

Cada tarefa em `tasks.md` é um item de checklist como `- [ ] T1.1: [CODE] Title`, seguido de:

- **`**Target:**`**: o caminho exato do arquivo ou o comando da CLI que a tarefa altera.
- **`**Context:**`**: o requisito que ela implementa (`[US-X]`), obrigatório a partir de `medium`.
- **`**Constraints:**`** *(opcional)*: apenas o que o implementador não consegue deduzir do design e do código (uma armadilha, um desvio da convenção, uma decisão de design referenciada). Sem assinaturas, constantes ou trechos de código.
- **`**Acceptance Check:**`**: um único comando não interativo ou cenário de teste com um resultado esperado observável. Critérios que exigem julgamento humano são reescritos ou sinalizados como não verificáveis.

Regras:
- Toda tarefa altera pelo menos um arquivo do produto ou artefato. Tarefas apenas de verificação (rodar a suíte completa, uma execução de smoke, uma checagem de build) são proibidas: isso pertence à etapa de QA independente e aos hooks `on_all_tasks_finished`.
- As verificações black-box ficam em uma seção não executável `## 3. Black-box Scenarios` após as fases de tarefas, nunca como tarefas. Cada cenário (`S1`, `S2`...) se vincula a um requisito e informa um comando e seu resultado esperado. A seção é obrigatória para `large` e `complex`, opcional para especificações `medium` que produzem um artefato executável, e omitida em `small`.
- `specforce spec status` mostra avisos não bloqueantes para tarefas que contêm blocos de código ou citam mais de três caminhos de arquivo.

## 6. Deltas de módulo

A partir de `medium`, uma especificação deve declarar, antes da implementação, como ela altera as especificações vivas dos módulos em `.specforce/docs/modules/`:

- **Documentos de delta**: um `deltas/<module>.md` por módulo afetado, com seções `## ADDED`, `## MODIFIED` e `## REMOVED` contendo blocos completos. Novas regras e casos de uso usam chaves provisórias `NEW:<key>`; os identificadores finais são alocados no momento do arquivamento.
- **Sem impacto**: quando nenhum módulo muda, a especificação declara `delta.impact: none` no `spec.yaml`.
- **Selo**: `specforce spec delta seal <slug>` registra impressões digitais dos blocos que o delta modifica ou remove. Se outra especificação alterar o mesmo bloco depois, a verificação do delta reporta um conflito e o delta precisa ser refeito com base no módulo atual e selado de novo.
- **Estados** reportados por `specforce spec delta check <slug>`:
  - `none`: `delta.impact: none` está declarado.
  - `exempt`: uma especificação `small` sem delta.
  - `missing`: um delta é obrigatório, mas não foi escrito.
  - `pending`: o delta ainda não está refletido nos módulos.
  - `applied`: os módulos correspondem ao delta.
  - `invalid`: o delta está malformado ou conflita com os módulos atuais.
- **Barreira de arquivamento**: o `spf.archive` aplica o delta aos módulos; uma especificação só é arquivada quando seu estado é `applied`, `none` ou `exempt`.

## 7. Relatório de QA

Após a implementação, o `spf.implement` executa uma etapa de QA independente: um testador que não implementou nenhuma tarefa verifica cada critério de aceitação (e cada item do delta e cenário black-box) e dá a cada um deles um veredito:

- **`PASSED`**: verificado com evidência.
- **`FAILED`**: o critério não é atendido. Critérios que falharam disparam no máximo dois ciclos de correção antes de o agente escalar para você.
- **`NOT-VERIFIABLE`**: o critério não pode ser executado por um testador; o motivo é informado.

O relatório é armazenado no `spec.yaml` (um relatório por especificação, substituído a cada nova execução). Seu status geral é `FAILED` quando qualquer critério falhou e, caso contrário, `PASSED`. `specforce implementation status` o mostra.

Barreira de arquivamento: um relatório `FAILED` bloqueia o arquivamento (mesmo com `--force`); a ausência de relatório só produz um aviso.

## 8. Tarefas test-first (RED)

Uma tarefa cujo título contém o token exato `[RED]` (diferencia maiúsculas de minúsculas) é uma tarefa test-first: ela escreve um teste que deve falhar até que uma tarefa `[GREEN]` posterior implemente o comportamento.

- Quando todas as tarefas finalizadas em um lote são `[RED]`, o Specforce pula apenas os hooks por tarefa `on_task_finished` e imprime um aviso citando-os, já que uma suíte de testes falhando é esperada nesse ponto.
- Os hooks `on_phase_finished` e `on_all_tasks_finished` continuam rodando, mesmo quando uma tarefa `[RED]` fecha a fase.
- Um lote que contenha qualquer tarefa não marcada como `[RED]` roda todos os hooks normalmente.

Para tornar os roadmaps test-first, adicione uma instrução em `tasks` no `config.yaml` (veja [Configuração](configuration.md)).

## 9. A pasta `.specforce/`

A pasta `.specforce/` da raiz guarda o estado da orquestração. Com exceção do arquivo de configuração, ela não foi feita para ser editada à mão:

- **`config.yaml`**: instruções do projeto, variáveis de contexto e hooks de validação (veja [Configuração](configuration.md)).
- **`docs/`**: a Constituição e as especificações vivas dos módulos.
- **`specs/`**: especificações ativas.
- **`archive/`**: especificações encerradas pelo `spf.archive`, o histórico das decisões técnicas do projeto.
- **`scratch/`**: evidências descartáveis criadas pelo `spf.discovery` (testes, scripts, logs, capturas de tela). Ela se ignora no controle de versão e é tratada como desatualizada em sessões posteriores.

---

Manter a qualidade desses artefatos é a melhor forma de garantir que a IA produza um código excelente.
