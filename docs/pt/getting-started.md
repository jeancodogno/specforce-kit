[ [English](../en/getting-started.md) | [Português](getting-started.md) | [Español](../es/getting-started.md) ]

# Primeiros Passos com o Specforce Kit

Bem-vindo ao **Specforce Kit**! Este guia vai ajudar você a configurar seu ambiente e dar os primeiros passos com o Desenvolvimento Orientado por Especificação (SDD - Spec-Driven Development).

## Pré-requisitos

- Node.js e NPM (para a instalação recomendada)
- Go 1.26+ e Make (se você quiser compilar a partir do código-fonte)
- Um [agente de IA suportado](supported-tools.md) (Claude Code, Qwen, Kimi Code, OpenCode, KiloCode, Codex, Antigravity ou Cursor)

## Instalação

A forma mais fácil de instalar o Specforce Kit é globalmente via NPM:

```bash
npm i -g @jeancodogno/specforce-kit
```

Se preferir instalar a partir do código-fonte:

```bash
git clone https://github.com/jeancodogno/specforce-kit.git
cd specforce-kit
make build
sudo make install   # installs to /usr/local/bin; use "make install PREFIX=$HOME/.local" to avoid sudo
```

### Solução de problemas

Se você encontrar o erro `command not found: specforce` após a instalação, normalmente significa que o `PATH` do seu sistema não inclui o diretório onde os binários ficam armazenados.

#### 1. Correção para macOS/Linux (NPM)
Se você instalou via `npm i -g`, adicione o seguinte ao seu perfil de shell (por exemplo, `~/.zshrc` ou `~/.bashrc`):

```bash
# Add npm global bin to PATH
export PATH="$(npm config get prefix)/bin:$PATH"
```

Depois recarregue o perfil: `source ~/.zshrc`.

#### 2. Correção para Windows (NPM)
Certifique-se de que o diretório global do npm (normalmente `%AppData%\npm`) esteja nas Variáveis de Ambiente do sistema.

#### 3. Correção para instalação a partir do código-fonte/Git
Se você instalou a partir do código-fonte ou do Git e o comando não é encontrado, confirme que compilou o binário nativo:

```bash
cd specforce-kit
make build
# This creates a 'specforce' binary in the root
```

Se a compilação automática falhou durante a instalação, verifique se você tem `go` e `make` instalados no seu sistema.

#### 4. Diagnóstico automático
Se o comando `specforce` funciona mas não encontra o seu "cérebro" interno (o binário nativo), ele executa automaticamente uma verificação de diagnóstico e imprime o comando exato que você precisa rodar para corrigir o seu `PATH`.

## Atualizações

O Specforce se mantém atualizado em segundo plano, seguindo o canal pelo qual você o instalou:

- **Quando ele verifica:** no máximo uma vez a cada 6 horas, e apenas quando você executa um comando interativo em um terminal. A verificação roda desacoplada e nunca deixa o comando mais lento; estar offline não causa verificações repetidas.
- **Onde ele procura:** instalações via npm seguem a tag npm do canal em execução (por exemplo `beta` para uma versão `1.0.0-beta.x`, `latest` para versões estáveis); instalações nativas seguem a release mais recente do GitHub. O Specforce nunca faz downgrade.
- **Quando ele aplica:** instalações via npm instalam a nova versão exata; binários nativos são baixados, verificados contra o checksum da release e substituídos na próxima execução. O próximo comando interativo mostra `(Specforce updated to <version>)` uma vez.
- **Quando ele não consegue:** se a pasta de instalação não tiver permissão de escrita, ou se uma instalação ou verificação falhar, nada é repetido em loop e um aviso **MANUAL UPDATE REQUIRED** é mostrado uma vez para aquela versão, com o comando exato:

```bash
# npm installs
npm install -g @jeancodogno/specforce-kit@<version>

# native installs: download the release from GitHub, or reinstall from source
sudo make install
```

- **Os agentes nunca são incomodados:** comandos executados por agentes (`constitution`, `spec`, `implementation`) nunca iniciam uma verificação nem mostram avisos, e os avisos são escritos apenas na saída de erro padrão, então a saída dos agentes (incluindo `--json`) nunca é poluída.
- **Diagnóstico:** execute qualquer comando com `DEBUG=1` (por exemplo `DEBUG=1 specforce console`) para imprimir cada decisão do atualizador na saída de erro padrão. O estado do atualizador por usuário fica em `~/.specforce/state.json`.

## Primeiros Passos

As skills são citadas pelo nome a seguir (por exemplo `spf.constitution`). Como chamá-las explicitamente depende da sua ferramenta: veja [Invocando skills](supported-tools.md#invocando-skills). Um pedido em linguagem natural que cite o nome da skill funciona em qualquer ferramenta.

### 1. Inicialize seu projeto

Navegue até o diretório do seu projeto (novo ou existente) e inicialize o Specforce:

```bash
specforce init
```

Escolha suas ferramentas de IA na lista interativa, ou passe as chaves delas diretamente (por exemplo `specforce init claude cursor`). Isso cria a pasta `.specforce/`, as regras no `AGENTS.md` da raiz e as skills de cada ferramenta selecionada.

### 2. Defina a Constituição

A **Constituição** é o conjunto fundamental de regras do seu projeto (princípios, arquitetura, UI/UX, segurança, engenharia, governança). Em vez de poluir o contexto do seu agente com dezenas de regras soltas, o Specforce as centraliza de forma estruturada.

Inicie seu agente de IA e peça, no chat dele:

```text
Use spf.constitution para configurar a constituição: <descrição do seu projeto e stack tecnológica>
```

O agente fará algumas perguntas para alinhar o projeto e criará os documentos fundamentais na pasta `.specforce/docs/`.

### 3. Crie sua primeira especificação

Toda mudança de código no Specforce exige uma especificação. Para adicionar uma nova funcionalidade, peça ao seu agente:

```text
Use spf.spec para planejar: <descreva a funcionalidade que você quer construir>
```

O agente cria os artefatos exigidos pelo tamanho da especificação; uma especificação `medium`, o tamanho padrão, precisa de `requirements.md` e `tasks.md` (veja [Artefatos](artifacts.md)).

---

## Fluxo de Trabalho: O Ciclo de Vida

O fluxo de trabalho do Specforce Kit é baseado no **Desenvolvimento Orientado por Especificação (SDD)**. Na era da IA, a implementação virou uma commodity, mas um design arquitetural sem validação é um risco enorme.

O SDD separa o **design** da **implementação**, garantindo que nenhum código seja gerado no improviso ("vibe coding") sem um contrato formal. O ciclo de desenvolvimento segue um caminho determinístico. Para fluxos prontos, dicas e o que fazer quando uma etapa fica bloqueada, veja [Fluxos de Trabalho e Dicas](workflows.md).

### 1. Fundação: A Constituição (`spf.constitution`)
Tudo começa pela **Constituição**. Antes de escrever qualquer código, o projeto precisa de regras claras de arquitetura, UI/UX, segurança e engenharia.
- O agente interage com você para entender o escopo do projeto, ou lê o código existente de um projeto que você está adotando.
- Ele documenta essas decisões em `.specforce/docs/`, que servem como as "leis" de todas as implementações futuras.
- **Vantagem:** O agente carrega apenas os documentos de que precisa, economizando tokens e evitando alucinações.

### 2. Exploração: A Descoberta (`spf.discovery`)
Use a descoberta para fazer um brainstorming de uma ideia, investigar um bug ou encontrar uma causa raiz antes de se comprometer com um plano.
- O código do produto permanece somente leitura; as evidências descartáveis (testes, scripts, capturas de tela) ficam em `.specforce/scratch/`.
- Quando estiver pronto, peça ao agente para formalizar o resultado em um `proposal.md` que o planejamento reaproveita.

### 3. O Contrato: A Especificação (`spf.spec`)
Toda nova funcionalidade, correção ou refatoração começa aqui.
- Você explica o que precisa. O agente escreve os artefatos exigidos pelo tamanho da especificação:
  - **Requisitos** (o "quê")
  - **Design** (o "como", do ponto de vista arquitetural)
  - **Tarefas** (o roadmap atômico de execução)
  - **Deltas de módulo** (como a mudança afeta as especificações vivas dos módulos)
- *Ação humana:* Você deve revisar e aprovar esta especificação. Este é o seu contrato com a máquina.

### 4. Execução: A Implementação (`spf.implement`)
Com a especificação aprovada, o agente entra em modo de execução.
- O agente segue rigorosamente a lista de tarefas gerada na etapa anterior e depois executa uma etapa de QA independente que registra um veredito para cada critério de aceitação.
- Os [Hooks de Validação](configuration.md) garantem que o agente não pule testes, linters ou etapas de segurança antes de marcar uma tarefa como "finalizada".
- *Ação humana:* Use o comando `specforce console` em outro terminal para acompanhar o progresso do agente em tempo real.

### 5. Encerramento: O Arquivamento (`spf.archive`)
Depois que a funcionalidade está totalmente codificada e verificada, a especificação precisa ser encerrada.
- O agente aplica os deltas de módulo às especificações vivas em `.specforce/docs/modules/`, atualiza a Constituição se surgiram novas decisões sistêmicas e move a especificação para `.specforce/archive/`.
- O arquivamento é recusado enquanto houver tarefas pendentes, o delta não estiver aplicado ou o veredito de QA for FAILED.
- Você pode personalizar essas regras no `config.yaml` para incluir etapas como atualizar ferramentas de acompanhamento externas ou tarefas específicas de limpeza.

---

## Papel Humano vs. Papel da Máquina

No Specforce, os papéis são claros:

- **O Arquiteto Humano (Governança):** Define a direção do projeto, aprova as especificações e garante que o sistema como um todo faça sentido para o negócio.
- **O Agente de IA (Executor):** Executa o roadmap aprovado, gerando código, testes e documentação interna de alta qualidade com eficiência e rapidez, sempre respeitando a Constituição.

---

Para explorar os comandos de terminal, veja a [Referência da CLI](cli.md).
Saiba mais sobre os documentos gerados neste processo na página de [Artefatos](artifacts.md).
Siga fluxos passo a passo e dicas em [Fluxos de Trabalho e Dicas](workflows.md).
