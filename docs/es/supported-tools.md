[ [English](../en/supported-tools.md) | [Português](../pt/supported-tools.md) | [Español](supported-tools.md) ]

# Agentes y herramientas de IA soportados

Specforce es agnóstico respecto a la herramienta: la Constitución, las especificaciones y la CLI son las mismas sea cual sea el asistente de programación con IA que uses. Además, instala sus *Skills* de forma nativa para estas herramientas:

- **Claude Code** (Anthropic)
- **Qwen** (Alibaba)
- **Kimi Code** (Moonshot AI)
- **OpenCode**
- **KiloCode**
- **Codex**
- **Antigravity**
- **Cursor**

## Instalar las skills

`specforce init` abre un selector interactivo donde eliges las herramientas a configurar. Para saltarte el selector, pasa directamente las claves de las herramientas:

```bash
specforce init claude cursor
```

Claves de herramienta aceptadas: `claude`, `qwen`, `kimi-code`, `open-code`, `kilo-code`, `codex`, `antigravity`, `cursor`.

Todas las herramientas reciben las mismas skills: `spf.discovery`, `spf.constitution`, `spf.spec`, `spf.implement` y `spf.archive`, más `consultative-grill`, un auxiliar interno que usa `spf.spec` durante la planificación. Cada skill se instala como una carpeta llamada `spf-<name>` (por ejemplo `spf-spec/SKILL.md`) cuyo nombre visible es `spf.<name>`.

## Invocar las skills

Todas las herramientas pueden elegir una skill por su cuenta cuando tu petición coincide con ella, y todas entienden una petición en lenguaje natural que nombre la skill, por ejemplo "Usa spf.spec para planificar: añadir exportación CSV a los informes". La tabla muestra dónde se instalan las skills y la sintaxis de invocación explícita documentada por cada herramienta, usando `spf.spec` como ejemplo.

| Herramienta | Carpeta de instalación | Invocación explícita | Notas |
|---|---|---|---|
| Claude Code | `.claude/skills/` | `/spf.spec <request>` (el nombre de la carpeta, `/spf-spec`, también funciona) | Clave `claude`. También recibe `.claude/rules/AGENTS.md` (consulta más abajo). |
| Qwen | `.qwen/skills/` | `/spf.spec <request>` | Clave `qwen`. |
| Kimi Code | `.kimi/skills/` | `/skill:spf.spec <request>` | Clave `kimi-code`. |
| OpenCode | `.opencode/skills/` | Menciona `@spf-spec` en tu prompt | Clave `open-code`. |
| KiloCode | `.kilocode/skills/` | Pídela por su nombre en el chat (p. ej. 'usa spf.spec ...') | Clave `kilo-code`. |
| Codex | `.codex/skills/` | `/prompts:spf-spec <request>` | Clave `codex`. También escribe los archivos de prompt `spf-*.md` **fuera del proyecto**, en `${CODEX_PROMPTS_DIR:-${CODEX_HOME:-~/.codex}/prompts}`: la carpeta definida por `CODEX_PROMPTS_DIR`, si no `$CODEX_HOME/prompts`, y si no `~/.codex/prompts`. |
| Antigravity | `.agents/skills/` | Pídela por su nombre en el chat (p. ej. 'usa spf.spec ...') | Clave `antigravity`. |
| Cursor | `.cursor/skills/` | Pídela por su nombre en el chat (p. ej. 'usa spf.spec ...') | Clave `cursor`. Lee el `AGENTS.md` raíz de forma nativa. |

La tabla usa `spf.spec` como ejemplo. Para las demás skills, cambia solo el nombre: en Claude Code, `/spf.discovery` o `/spf.archive`; en Codex, `/prompts:spf-implement`. En las herramientas marcadas con "Pídela por su nombre en el chat" no hay comando: escribe la petición como siempre nombrando la skill, por ejemplo "usa spf.archive para archivar csv-export".

## Configuración automatizada

Cuando ejecutas `specforce init` o actualizas las herramientas, Specforce se asegura de que cada agente pueda descubrir las reglas del proyecto y sus skills:

- **`AGENTS.md` raíz:** Specforce escribe sus reglas entre los marcadores `<!-- SPECFORCE_AGENTS_START -->` y `<!-- SPECFORCE_AGENTS_END -->`. Todo lo que escribas fuera de ese bloque se conserva.
- **Claude Code:** solo Claude Code recibe un enlace: `.claude/rules/AGENTS.md` se crea como enlace simbólico al `AGENTS.md` raíz.
- **Cursor y las demás herramientas:** no se crea ningún enlace. Cursor lee el `AGENTS.md` raíz de forma nativa; las demás herramientas leen directamente la guía raíz o su carpeta de skills.
- **Enlaces antiguos:** se eliminan los enlaces simbólicos `AGENTS.md` antiguos dentro de las carpetas `.agents/**/rules/`.

*¿No encuentras tu agente favorito? ¡Envía un PR creando un Kit para él!*
