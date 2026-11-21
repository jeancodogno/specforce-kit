[ [English](../en/supported-tools.md) | [Português](supported-tools.md) | [Español](../es/supported-tools.md) ]

# Agentes de IA e Ferramentas Suportadas

O Specforce é agnóstico a ferramentas: a Constituição, as especificações e a CLI são as mesmas qualquer que seja o assistente de programação com IA que você use. Além disso, ele instala suas *Skills* nativamente para estas ferramentas:

- **Claude Code** (Anthropic)
- **Qwen** (Alibaba)
- **Kimi Code** (Moonshot AI)
- **OpenCode**
- **KiloCode**
- **Codex**
- **Antigravity**
- **Cursor**

## Instalando as skills

O `specforce init` abre um seletor interativo em que você escolhe as ferramentas a configurar. Para pular o seletor, passe as chaves das ferramentas diretamente:

```bash
specforce init claude cursor
```

Chaves de ferramenta aceitas: `claude`, `qwen`, `kimi-code`, `open-code`, `kilo-code`, `codex`, `antigravity`, `cursor`.

Todas as ferramentas recebem as mesmas skills: `spf.discovery`, `spf.constitution`, `spf.spec`, `spf.implement` e `spf.archive`, mais `consultative-grill`, um auxiliar interno usado pelo `spf.spec` durante o planejamento. Cada skill é instalada como uma pasta chamada `spf-<name>` (por exemplo `spf-spec/SKILL.md`) cujo nome de exibição é `spf.<name>`.

## Invocando skills

Toda ferramenta pode escolher uma skill por conta própria quando o seu pedido corresponde a ela, e toda ferramenta entende um pedido em linguagem natural que cite o nome da skill, por exemplo "Use spf.spec para planejar: adicionar exportação CSV aos relatórios". A tabela mostra onde as skills são instaladas e a sintaxe de invocação explícita documentada por cada ferramenta, usando `spf.spec` como exemplo.

| Ferramenta | Pasta de instalação | Invocação explícita | Observações |
|---|---|---|---|
| Claude Code | `.claude/skills/` | `/spf.spec <request>` (o nome da pasta, `/spf-spec`, também funciona) | Chave `claude`. Também recebe `.claude/rules/AGENTS.md` (veja abaixo). |
| Qwen | `.qwen/skills/` | `/spf.spec <request>` | Chave `qwen`. |
| Kimi Code | `.kimi/skills/` | `/skill:spf.spec <request>` | Chave `kimi-code`. |
| OpenCode | `.opencode/skills/` | Mencione `@spf-spec` no seu prompt | Chave `open-code`. |
| KiloCode | `.kilocode/skills/` | Peça pelo nome no chat (ex.: 'use spf.spec ...') | Chave `kilo-code`. |
| Codex | `.codex/skills/` | `/prompts:spf-spec <request>` | Chave `codex`. Também escreve os arquivos de prompt `spf-*.md` **fora do projeto**, em `${CODEX_PROMPTS_DIR:-${CODEX_HOME:-~/.codex}/prompts}`: a pasta definida por `CODEX_PROMPTS_DIR`, senão `$CODEX_HOME/prompts`, senão `~/.codex/prompts`. |
| Antigravity | `.agents/skills/` | Peça pelo nome no chat (ex.: 'use spf.spec ...') | Chave `antigravity`. |
| Cursor | `.cursor/skills/` | Peça pelo nome no chat (ex.: 'use spf.spec ...') | Chave `cursor`. Lê o `AGENTS.md` da raiz nativamente. |

A tabela usa `spf.spec` como exemplo. Para as outras skills, troque só o nome: no Claude Code, `/spf.discovery` ou `/spf.archive`; no Codex, `/prompts:spf-implement`. Nas ferramentas marcadas com "Peça pelo nome no chat", não existe comando: escreva o pedido normalmente citando a skill, por exemplo "use spf.archive para arquivar csv-export".

## Configuração Automatizada

Quando você executa o `specforce init` ou atualiza as ferramentas, o Specforce garante que cada agente consiga descobrir as regras do projeto e as suas skills:

- **`AGENTS.md` da raiz:** O Specforce escreve suas regras entre os marcadores `<!-- SPECFORCE_AGENTS_START -->` e `<!-- SPECFORCE_AGENTS_END -->`. Tudo o que você escrever fora desse bloco é preservado.
- **Claude Code:** Apenas o Claude Code recebe um link: `.claude/rules/AGENTS.md` é criado como um link simbólico para o `AGENTS.md` da raiz.
- **Cursor e as demais ferramentas:** Nenhum link é criado. O Cursor lê o `AGENTS.md` da raiz nativamente; as demais ferramentas leem o guia da raiz ou a sua pasta de skills diretamente.
- **Links legados:** Links simbólicos `AGENTS.md` antigos dentro de pastas `.agents/**/rules/` são removidos.

*Não encontrou seu agente favorito? Envie um PR criando um Kit para ele!*
