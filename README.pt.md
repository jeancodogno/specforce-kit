
<p align="center">
  <a href="https://github.com/jeancodogno/specforce-kit/">
    <picture>
      <img src="assets/logo.png" alt="Logo OpenSpec" height="128">
    </picture>
  </a>
</p>
<p align="center">Ecossistema para desenvolvimento assistido por IA.</p>

# Specforce Kit

🌎 **Idiomas / Languages:** [English](README.md) | [Português](README.pt.md) | [Español](README.es.md)

> **A Camada de Orquestração para Engenharia Nativa de IA.**

O Specforce Kit é um framework de nível profissional projetado para coordenar agentes de IA (Claude Code, Cursor, KiloCode) através do **Desenvolvimento Orientado por Especificação (SDD - Spec-Driven Development)**. Ele fornece os **Kits**, a **Constituição** e as **Ferramentas de Orquestração** padronizadas necessárias para transformar a IA de um assistente "baseado em vibes" em uma força de engenharia determinística.

---

## 🧠 Filosofia

Na era da IA, a implementação é uma commodity, mas o design é uma responsabilidade. **A implementação sem uma especificação é estritamente proibida no Specforce.**

O Desenvolvimento Orientado por Especificação (SDD) desacopla o design da implementação. Uma especificação atua como um contrato rigoroso entre designers humanos e executores de máquina. Ao fazer isso, evitamos o "vibe coding" - o acúmulo de decisões ad-hoc e não verificadas - garantindo que seu sistema permaneça coerente, seguro e sustentável, independentemente de qual agente de IA escreva o código.

---

## 📊 Por que usar o Specforce?

A codificação por IA sem diretrizes rígidas leva inevitavelmente ao "vibe coding": prompts vagos, resultados imprevisíveis e uma montanha de dívida técnica que se acumula rapidamente. Embora outras ferramentas tentem resolver isso, elas muitas vezes forçam você a novos ecossistemas de IDE, exigem configurações pesadas ou incham a janela de contexto da sua IA.

O Specforce foi projetado para ser a **camada de orquestração invisível**. Ele traz previsibilidade, controle arquitetônico e verificação rigorosa aos seus fluxos de trabalho de IA, atuando como uma máquina de estados ativa em vez de apenas uma coleção de templates estáticos.

### Principais Diferenciais:

- **Constituição Segmentada**: Em vez de um único prompt massivo, as regras são divididas em domínios especializados (Arquitetura, UI/UX, Segurança). O agente carrega apenas o contexto exato de que precisa, mantendo os custos de tokens baixos e reduzindo alucinações.
- **Hooks de Verificação Dinâmica**: Execute testes, linters ou scripts personalizados automaticamente antes que um agente seja autorizado a marcar uma tarefa como finalizada.
- **Suporte a Git Worktree**: Descubra e visualize nativamente especificações em múltiplas branches git simultaneamente, sem trocar de contexto.
- **Instruções Sensíveis ao Contexto**: Injete regras e restrições específicas dinamicamente com base no artefato ou fase em que o agente está trabalhando.
- **Agnóstico a Ferramentas e IDEs**: Traga o Desenvolvimento Orientado por Especificação diretamente para o seu terminal. Funciona perfeitamente com sua IDE favorita e agentes de IA baseados em CLI (Gemini, Claude, Cursor, KiloCode, etc).
- **Rápido e Sem Dependências**: Distribuído como um único binário, sem necessidade de configurações pesadas de Python ou ecossistemas forçados para iniciar seu projeto.

---

## 🤖 Uso pelo Agente (Comandos Slash)

O Specforce interage perfeitamente com seu agente de IA preferido através de comandos slash simples. Aqui está o fluxo de trabalho padrão:

### 1. Descoberta (`/spf:discovery`)
**Comando:** `/spf:discovery {ideia ou relatório de erro}`

A fase inicial de exploração. Use isso para debater novas funcionalidades ou investigar problemas técnicos sem modificar nenhum arquivo.
- **Exploração Somente Leitura:** O agente é estritamente proibido de escrever código ou artefatos. Ele se concentra em pesquisa, análise de causa raiz e estratégia técnica.
- **Personas Especialistas:** Alterna automaticamente entre "Arquiteto de Produto Sênior" (Brainstorming) e "Engenheiro de Sistemas Sênior" (Detetive) com base na sua entrada.
- **Alinhamento com o Projeto:** Lê a Constituição do seu projeto (`.specforce/docs/`) para garantir que todas as ideias estejam alinhadas com seus princípios e arquitetura.

### 2. A Constituição (`/spf:constitution`)
**Comando:** `/spf:constitution {descrição do projeto}`

O primeiro passo em qualquer projeto. Isso gera a Constituição do seu projeto, contendo todas as regras, princípios, diretrizes de UI/UX, arquitetura, segurança e memória do agente.
- **Agnóstico a Ferramentas e Segmentado:** Em vez de despejar regras em arquivos específicos de ferramentas (como `.clauderc` ou `.gemini/GEMINI.md`) ou manter tudo no contexto, o Specforce mantém seus próprios arquivos de memória especializados e segmentados. O agente carrega apenas o que precisa. Você pode trocar de ferramenta (ex: de Gemini CLI para Claude Code) e a memória do projeto permanece completamente intacta.
- **Configuração Interativa:** Descreva sua ideia de projeto e stack. O agente fará perguntas esclarecedoras para ajudá-lo a tomar decisões fundamentais.
- **Flexível:** Execute uma vez para configurar, ou execute novamente a qualquer momento para atualizar a constituição. Você pode usá-lo em ideias novas ou apontá-lo para um projeto existente para analisar e documentar sua arquitetura atual.

### 2. Criando uma Spec (`/spf:spec`)
**Comando:** `/spf:spec {descrição do que adicionar ou modificar}`

Pronto para construir? O agente esclarecerá seus requisitos e gerará três documentos principais: `requirements.md`, `design.md` e `tasks.md`.

### 3. Implementação (`/spf:implement`)
**Comando:** `/spf:implement`

Após validar os documentos de especificação gerados, execute este comando. O agente seguirá estritamente as tarefas planejadas na spec para implementar a funcionalidade, garantindo desvio zero do design acordado.

### 4. Arquivamento (`/spf:archive`)
**Comando:** `/spf:archive`

Assim que a implementação estiver totalmente concluída e verificada, este comando analisa o que precisa ser atualizado na Constituição global do projeto com base na nova implementação e arquiva a spec concluída. Ele segue **Instruções de Arquivamento** padronizadas que garantem que as lições aprendidas sejam capturadas no memorial do projeto e qualquer limpeza específica do projeto seja realizada.

---

## ⚡ Toque Humano

Enquanto o Specforce potencializa seus agentes de IA, os humanos interagem principalmente através de dois comandos de terminal:

### 1. Inicialização do Projeto
Configure a governança e os kits de agentes para o seu projeto.
```bash
specforce init
```

### 2. O Console (TUI)
Inicie o centro de comando para monitorar a saúde do projeto, o progresso das especificações e as implementações dos agentes em tempo real.
```bash
specforce console
```

---

## 🛠️ Agentes de IA e Ferramentas Suportadas

O Specforce Kit foi projetado para ser agnóstico a ferramentas, mas vem com kits pré-construídos e suporte nativo a comandos slash para os assistentes de codificação de IA mais populares:

- **Gemini CLI** (Google)
- **Claude Code** (Anthropic)
- **Qwen** (Alibaba)
- **Kimi Code** (Moonshot AI)
- **OpenCode**
- **KiloCode**
- **Codex**
- **Antigravity**

*Não vê sua ferramenta favorita? Aceitamos PRs para adicionar novos kits!*

---

## 🚀 Primeiros Passos

### Instalação

**Usando NPM (Recomendado)**
```bash
npm i -g @jeancodogno/specforce-kit
```

**A partir do Código Fonte (Requer Go 1.26+ e Make)**
```bash
git clone https://github.com/jeancodogno/specforce-kit.git
cd specforce-kit
make build
make install
```

---

## 📚 Documentação

Para um mergulho mais profundo no Specforce Kit, confira nossa documentação oficial:

- [Primeiros Passos e Fluxo de Trabalho](docs/pt/getting-started.md): Configuração do ambiente, primeiros passos e o ciclo de vida recomendado do Desenvolvimento Orientado por Especificação.
- [Suporte a Git Worktree](docs/pt/git-worktrees.md): Guia para descoberta de especificações entre branches, console unificado e restrições de apenas leitura.
- [Artefatos](docs/pt/artifacts.md): Detalhes sobre os arquivos gerados (Constituição, Requisitos, Design, Tarefas).
- [Configuração](docs/pt/configuration.md): Guia para personalizar hooks e instruções.
- [Referência da CLI](docs/pt/cli.md): Uso do terminal e comandos slash.
- [Ferramentas Suportadas](docs/pt/supported-tools.md): Lista de agentes de IA e assistentes de codificação suportados nativamente.

---

## 🔍 Resolução de Problemas

### "command not found: specforce"
Se você encontrar este erro após a instalação, geralmente significa que o diretório bin global do npm não está no `PATH` do seu sistema.

**Correção para macOS/Linux:**
Adicione a seguinte linha ao seu perfil de shell (ex: `~/.zshrc` ou `~/.bashrc`):
```bash
export PATH="$(npm config get prefix)/bin:$PATH"
```

**Correção para Windows:**
Certifique-se de que `%AppData%\npm` esteja em suas variáveis de ambiente.

Para uma resolução de problemas mais detalhada, veja [Primeiros Passos: Resolução de Problemas](docs/pt/getting-started.md#troubleshooting).

---

## 🤝 Contribuindo
Contribuições são bem-vindas! Veja nosso [CONTRIBUTING.md](CONTRIBUTING.md) para aprender como adicionar suporte para novos agentes de IA ou melhorar os kits existentes.

---

## 📜 Licença
O Specforce Kit é lançado sob a **Licença MIT**. Veja [LICENSE](LICENSE) para detalhes.
