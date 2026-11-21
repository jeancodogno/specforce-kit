[ [English](../en/cli.md) | [Português](cli.md) | [Español](../es/cli.md) ]

# Interface de Linha de Comando (CLI)

Esta página lista todos os comandos e flags da CLI `specforce`, agrupados por área. Muitos desses comandos são executados pelo agente de IA enquanto ele segue uma skill (veja [Invocando skills](supported-tools.md#invocando-skills)); eles estão documentados aqui para que você saiba o que existe e o que cada um faz.

Os comandos que recebem um `<slug>` aceitam o nome da pasta da especificação (por exemplo `20261121-0035-csv-export`) ou apenas a sua parte legível (`csv-export`). As flags marcadas com `--json` imprimem JSON legível por máquina, destinado a automações e agentes.

## Global

### `specforce init [agents...]`
Inicializa o Specforce no diretório atual: cria a pasta `.specforce/` (com um `config.yaml` padrão), o bloco gerenciado no `AGENTS.md` da raiz e as skills das ferramentas selecionadas. Sem argumentos, um seletor interativo permite escolher as ferramentas. Com argumentos, as chaves de ferramenta informadas (`claude`, `qwen`, `kimi-code`, `open-code`, `kilo-code`, `codex`, `antigravity`, `cursor`) são configuradas imediatamente. Execute-o de novo em um projeto já inicializado para adicionar ou remover ferramentas e atualizar as skills instaladas. Com argumentos, toda ferramenta instalada que você deixar de fora é removida apagando a pasta inteira dela (por exemplo `.claude/`, incluindo suas próprias configurações), sem confirmação; o seletor interativo avisa antes de apagar. Veja [Ferramentas Suportadas](supported-tools.md).

### `specforce console`
Abre o **Specforce Console**, uma interface interativa de terminal para navegar pelas especificações, pela Constituição e pelas tarefas de implementação, e para acompanhar o progresso em tempo real (incluindo especificações de [Git worktrees](git-worktrees.md) vinculadas). Deixe-o rodando em um painel de terminal separado enquanto seu agente trabalha.

---

## Flags globais e ambiente

- **`-v`, `--version`**: imprime a versão instalada (`specforce version <x.y.z>`).
- **`-h`, `--help`**: mostra a ajuda de qualquer comando ou grupo.
- **`--dev`**: usa a pasta local `kit/` em vez do kit embutido no binário. Destinada apenas a quem contribui com o kit.
- **`DEBUG=1`**: variável de ambiente que escreve diagnósticos na saída de erro padrão, incluindo cada decisão da atualização automática (veja [Atualizações](getting-started.md#atualizações)).

### `specforce completion <shell>`
Gera o script de autocompletar do shell para `bash`, `zsh`, `fish` ou `powershell`. Execute `specforce completion <shell> --help` para ver os passos de instalação de cada shell.

---

## Constituição

### `specforce constitution status [--json]`
Mostra o status de completude dos documentos da Constituição do projeto em `.specforce/docs/` (principles, architecture, ui-ux, security, engineering, governance e modules).

### `specforce constitution artifact <name> [--json]`
Mostra os detalhes de um artefato da Constituição: seu propósito, template e instruções. Os agentes o usam para buscar as regras certas quando escrevem ou atualizam esse documento.

---

## Especificação

### `specforce spec init <slug> [--type feature|bug] [--size small|medium|large|complex] [--json]`
Cria uma nova pasta de especificação `.specforce/specs/<timestamp>-<slug>/` com seus metadados `spec.yaml`.
- `--type`: `feature` (padrão) ou `bug`. Especificações de bug usam as variantes de bug dos requisitos e do design.
- `--size`: `small`, `medium` (padrão), `large` ou `complex`. O tamanho decide quais artefatos são obrigatórios (veja [Artefatos](artifacts.md)).

### `specforce spec list [--json]`
Lista todas as especificações ativas.

### `specforce spec status <slug> [--json]`
Mostra a completude de uma especificação: quais artefatos o seu tamanho exige, quais existem, erros de validação e avisos não bloqueantes sobre as tarefas.

### `specforce spec artifact <name> [--json]`
Mostra as instruções e o template de um artefato de especificação (por exemplo `requirements`, `design`, `tasks`, `bug-requirements`, `bug-design`), incluindo as instruções do projeto vindas do [`config.yaml`](configuration.md). Os agentes o usam antes de escrever cada artefato.

### `specforce spec resize <slug> --size <size> [--json]`
Altera o tamanho de uma especificação existente. Promovê-la marca os artefatos que passaram a ser obrigatórios como ausentes e bloqueia a implementação até que eles existam; rebaixá-la relaxa a validação e mantém os arquivos existentes em disco.

### `specforce spec audit <slug> [flags]`
Atualiza o estado da auditoria de refinamento armazenado no `spec.yaml`, usado pela skill de planejamento enquanto ela verifica a coerência entre os artefatos.
- `--valid` / `--invalid`: marca a especificação como válida ou inválida.
- `--error <msg>`: acrescenta uma mensagem de erro de coerência.
- `--clear`: limpa todos os erros de coerência.
- `--iteration <n>`: define o número da iteração de refinamento atual.
- `--json`: saída legível por máquina.

### `specforce spec archive <slug> [--force]`
Finaliza uma especificação concluída e move sua pasta para `.specforce/archive/`. Ele recusa quando:
- ainda há tarefas pendentes;
- uma especificação `medium`, `large` ou `complex` tem um delta que não está `applied` (ou declarado `none`); especificações `small` reportam `exempt`;
- o relatório de QA armazenado é `FAILED` (os critérios que falharam são listados).

Uma especificação sem relatório de QA é arquivada com um aviso visível. `--force` ignora apenas a verificação de tarefas pendentes; ele nunca contorna a verificação do delta nem um veredito de QA `FAILED`. Uma recusa deixa o sistema de arquivos intocado.

---

## Delta

O delta descreve as mudanças que uma especificação faz nas especificações vivas dos módulos em `.specforce/docs/modules/`. O arquivamento exige que o delta esteja aplicado, `none` ou `exempt`.

### `specforce spec delta check <slug> [--json]`
Reporta se o delta da especificação (as mudanças propostas para as especificações vivas dos módulos) já foi aplicado. Somente leitura. O estado reportado é um destes:
- `none`: a especificação declara `delta.impact: none` (sem impacto em módulos).
- `exempt`: uma especificação `small` sem delta; ela não precisa de um.
- `missing`: um delta é obrigatório, mas ainda não foi escrito.
- `pending`: existe um delta, mas suas mudanças ainda não estão refletidas nos módulos.
- `applied`: os módulos já correspondem ao delta.
- `invalid`: o delta está malformado ou conflita com os módulos atuais.

### `specforce spec delta seal <slug> [--json]`
Sela o delta da especificação registrando as impressões digitais dos blocos de módulo em que ele se baseia, para que edições manuais posteriores nesses blocos sejam detectadas como conflitos de base.

---

## Implementação e QA

### `specforce implementation status <slug> [--json]`
Mostra o roadmap de implementação de uma especificação (a partir do `tasks.md`): prontidão, artefatos ausentes para o seu tamanho, cada tarefa com seu estado e suas restrições, as mitigações do roadmap, as instruções de `implementation` do projeto e, quando existir, o estado de QA armazenado.

### `specforce implementation update <slug> --task <id>[,<id>...] --status <status>`
Atualiza o estado de uma ou mais tarefas. `--task` aceita IDs separados por vírgula ou pode ser repetido (`--task T1.1 --task T1.2`). `--status` é `pending`, `in-progress` ou `finished`.

Quando o status é `finished`, o Specforce primeiro executa os [hooks de validação](configuration.md#3-hooks-de-validação-quality-gates) configurados, uma vez por lote e sem duplicatas: `on_task_finished`, mais `on_phase_finished` quando uma tarefa é a última da sua fase e `on_all_tasks_finished` quando é a última do roadmap. Um hook que falha bloqueia a atualização. Quando todas as tarefas do lote têm o token `[RED]` no título, apenas os hooks `on_task_finished` são pulados, com um aviso citando-os; os hooks de fase e de todas as tarefas continuam rodando.

### `specforce implementation qa record <slug>`
Registra o relatório de QA independente de uma especificação. Ele lê o relatório como JSON na entrada padrão (critérios com cenário, evidência e veredito `PASSED`, `FAILED` ou `NOT-VERIFIABLE`, mais as descobertas fora do escopo), valida-o, deriva o status geral (`FAILED` quando qualquer critério falhou e, caso contrário, `PASSED`) e substitui o relatório anterior no `spec.yaml`. Uma entrada inválida é recusada sem alterar o relatório armazenado. A skill de implementação o usa após a etapa de QA.

---

## Arquivamento

### `specforce archive instructions`
Mostra o protocolo completo de arquivamento: a descrição de cada documento da Constituição, as regras centrais de arquivamento (aplicar o delta às especificações vivas dos módulos, atualizar a Constituição) e as instruções de `archive` do projeto vindas do `config.yaml`. A skill de arquivamento o segue antes de executar `specforce spec archive`.
