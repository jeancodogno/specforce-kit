# Contribuindo com o Specforce Kit

Antes de tudo, obrigado por considerar contribuir com o Specforce Kit! São pessoas como você que fazem do Specforce uma ótima ferramenta para toda a comunidade de Engenharia com IA.

## 🚀 Nossa Filosofia: Especificação Primeiro
Nós usamos o que construímos. **Todas as novas funcionalidades, correções de bugs ou mudanças arquiteturais neste repositório DEVEM seguir o ciclo de vida do Desenvolvimento Orientado por Especificação (SDD).** As skills são citadas pelo nome a seguir; veja [Invocando skills](docs/pt/supported-tools.md#invocando-skills) para a sintaxe da sua ferramenta de IA.

1.  **Especificação:** Toda mudança começa com uma especificação feita pelo seu agente de IA (`spf.spec`).
2.  **Implementação:** Depois de aprovada, implementamos seguindo o roadmap do Specforce (`spf.implement`).
3.  **Verificação e Arquivamento:** Testes automatizados e verificação manual na TUI são obrigatórios antes de fechar o ciclo (`spf.archive`).

---

## 🐛 Reportando Problemas

### Relatos de Bug
- Use o Issue tracker do GitHub.
- Descreva o comportamento esperado vs. o comportamento atual.
- Forneça um script ou passos mínimos de reprodução.
- Inclua seu sistema operacional e a versão do Specforce (`specforce --version`).

### Pedidos de Funcionalidade
- Adoramos ideias novas! Abra uma issue para discutir o "Porquê" e o "Valor" antes de entrar nos detalhes técnicos.

---

## 🛠 Processo de Pull Request

1.  **Fork e Branch:** Crie uma branch de funcionalidade a partir da `main`.
2.  **Crie uma Especificação:** Use seu agente de IA preferido (Claude Code, Cursor etc.) e peça, por exemplo, "Use spf.spec para planejar: {descrição}". Ele gera os artefatos exigidos pelo tamanho da especificação (`requirements.md`, `design.md`, `tasks.md` e os deltas de módulo). Use o `spf.discovery` antes quando o problema ou a causa raiz ainda não estiverem claros.
3.  **Implemente:** Peça "Use spf.implement para implementar {especificação}" para deixar o agente executar as tarefas e rodar a etapa de QA independente.
4.  **Padrões de Código:**
    - Linguagem: Go 1.26+
    - Formatação: `gofmt` é obrigatório.
    - Linting: Deve passar no `golangci-lint`.
    - Padrões: Siga a Constituição do projeto (`.specforce/docs/`).
5.  **Testes:** Toda lógica nova deve ter testes unitários. Testes de regressão são obrigatórios para correções de bugs.
6.  **Arquivamento e Commit:** Peça "Use spf.archive para arquivar {especificação}" para atualizar as especificações vivas dos módulos e a Constituição. Use [Conventional Commits](https://www.conventionalcommits.org/) (ex.: `feat:`, `fix:`, `docs:`).
7.  **Envie:** Abra um PR contra a branch `main`. Garanta que o pipeline de CI passe.

---

## 🤖 Desenvolvimento de Kits (Suporte a Novos Agentes)

O Specforce Kit foi construído para ser agnóstico a modelos. Contribuições que adicionem suporte a novos agentes de IA ou melhorem os existentes são bem-vindas.

### O que é um Kit?
O Kit é o conjunto de skills que ensina um agente de IA a trabalhar com o framework Specforce, mais o mapeamento que as instala no diretório oculto de cada ferramenta para descoberta nativa (ex.: `.claude/skills/`, `.cursor/skills/`, `.kilocode/skills/`). Ele é embutido no binário:

- **Blueprints de skills:** `src/internal/agent/kit/skills/<id>/SKILL.yaml` (por exemplo `spf-spec/SKILL.yaml`). Cada blueprint é traduzido no `SKILL.md` nativo de cada ferramenta suportada durante o `specforce init`.
- **Manifesto de ferramentas:** `src/internal/agent/kit/kit.yaml` declara todas as ferramentas suportadas: seu `name`, `description`, diretório `target` e os `mappings` das categorias que ela recebe (ex.: `skills`).

### Adicionando uma Nova Ferramenta
1.  **Registre o diretório:** Adicione o diretório oculto da ferramenta (ex.: `.newagent/`) à slice `ToolPrefixes` em `src/internal/core/constants.go`.
2.  **Defina o mapeamento:** Adicione uma entrada de ferramenta em `src/internal/agent/kit/kit.yaml` com seu `name`, `description`, `target` e `mappings`. Use a sobrescrita de `target` no nível do mapeamento e a expansão de variáveis de ambiente (ex.: `${VAR:-DEFAULT}`) quando a ferramenta precisar de arquivos fora da sua pasta padrão.
3.  **Escritas globais (opcional):** Se a ferramenta exigir arquivos fora da raiz do projeto, declare `security: { global_write: true }` na sua entrada do `kit.yaml` (veja a entrada `codex`) e documente esse local.
4.  **Garanta a Segmentação de Contexto:** Certifique-se de que a nova integração respeite o nosso modelo de **Constituição Segmentada**, carregando apenas as regras necessárias de arquitetura, UI/UX ou segurança quando preciso, em vez de despejar tudo na janela de contexto.

Para mudar o que as skills dizem, edite os blueprints `SKILL.yaml`; nunca edite os arquivos gerados dentro das pastas das ferramentas. Durante o desenvolvimento, execute a CLI com `--dev` a partir da raiz do repositório para carregar o kit de `src/internal/agent/kit/` em vez da cópia embutida no binário (veja a [Referência da CLI](docs/pt/cli.md)).

---

## 🏗 Ambiente de Desenvolvimento

### Pré-requisitos
- Go 1.26+
- Make

### Configuração
```bash
git clone https://github.com/jeancodogno/specforce-kit.git
cd specforce-kit
make build
make install
```

### Executando os Testes
```bash
make test
```

---

## 📜 Código de Conduta
Ao participar deste projeto, você concorda em seguir o Código de Conduta do projeto e em manter um ambiente profissional e respeitoso para todas as pessoas que contribuem.

---
<p align="center">Vamos construir juntos o futuro da Engenharia Orquestrada por IA.</p>
