[ [English](../en/getting-started.md) | [Português](../pt/getting-started.md) | [Español](getting-started.md) ]

# Primeros pasos con Specforce Kit

¡Bienvenido a **Specforce Kit**! Esta guía te ayudará a configurar tu entorno y a dar tus primeros pasos con el Desarrollo Orientado por Especificaciones (SDD, Spec-Driven Development).

## Requisitos previos

- Node.js y NPM (para la instalación recomendada)
- Go 1.26+ y Make (si deseas compilar desde el código fuente)
- Un [agente de IA soportado](supported-tools.md) (Claude Code, Qwen, Kimi Code, OpenCode, KiloCode, Codex, Antigravity o Cursor)

## Instalación

La forma más sencilla de instalar Specforce Kit es de forma global mediante NPM:

```bash
npm i -g @jeancodogno/specforce-kit
```

Si prefieres instalarlo desde el código fuente:

```bash
git clone https://github.com/jeancodogno/specforce-kit.git
cd specforce-kit
make build
sudo make install   # installs to /usr/local/bin; use "make install PREFIX=$HOME/.local" to avoid sudo
```

### Solución de problemas

Si encuentras un error `command not found: specforce` después de la instalación, normalmente significa que el `PATH` de tu sistema no incluye el directorio donde se guardan los binarios.

#### 1. Solución para macOS/Linux (NPM)
Si instalaste con `npm i -g`, añade lo siguiente al perfil de tu shell (por ejemplo, `~/.zshrc` o `~/.bashrc`):

```bash
# Add npm global bin to PATH
export PATH="$(npm config get prefix)/bin:$PATH"
```

Luego recarga tu perfil: `source ~/.zshrc`.

#### 2. Solución para Windows (NPM)
Asegúrate de que tu directorio global de npm (normalmente `%AppData%\npm`) esté en las Variables de Entorno del sistema.

#### 3. Solución para la instalación desde código fuente/Git
Si instalaste desde el código fuente o desde Git y el comando no aparece, asegúrate de haber compilado el binario nativo:

```bash
cd specforce-kit
make build
# This creates a 'specforce' binary in the root
```

Si la compilación automática falló durante la instalación, comprueba que tienes `go` y `make` instalados en tu sistema.

#### 4. Diagnóstico automático
Si el comando `specforce` funciona pero no encuentra su "cerebro" interno (el binario nativo), ejecutará automáticamente una comprobación de diagnóstico e imprimirá el comando exacto que debes ejecutar para corregir tu `PATH`.

## Actualizaciones

Specforce se mantiene actualizado en segundo plano, siguiendo el canal desde el que lo instalaste:

- **Cuándo comprueba:** como máximo una vez cada 6 horas, y solo cuando ejecutas un comando interactivo en una terminal. La comprobación se ejecuta de forma desacoplada y nunca ralentiza el comando; estar sin conexión no provoca comprobaciones repetidas.
- **Dónde busca:** las instalaciones con npm siguen la etiqueta npm del canal en uso (por ejemplo `beta` para una versión `1.0.0-beta.x`, `latest` para las versiones estables); las instalaciones nativas siguen la última release de GitHub. Specforce nunca instala una versión anterior.
- **Cuándo aplica:** las instalaciones con npm instalan la nueva versión exacta; los binarios nativos se descargan, se verifican contra el checksum de la release y se sustituyen en la siguiente ejecución. El siguiente comando interactivo muestra `(Specforce updated to <version>)` una sola vez.
- **Cuándo no puede:** si la carpeta de instalación no tiene permisos de escritura, o si una instalación o verificación falla, no se reintenta nada en bucle y se muestra una sola vez para esa versión un aviso **MANUAL UPDATE REQUIRED**, con el comando exacto:

```bash
# npm installs
npm install -g @jeancodogno/specforce-kit@<version>

# native installs: download the release from GitHub, or reinstall from source
sudo make install
```

- **Los agentes nunca se ven afectados:** los comandos ejecutados por agentes (`constitution`, `spec`, `implementation`) nunca inician una comprobación ni muestran un aviso, y los avisos solo se escriben en la salida de error estándar, así que la salida de los agentes (incluido `--json`) nunca se contamina.
- **Diagnóstico:** ejecuta cualquier comando con `DEBUG=1` (por ejemplo `DEBUG=1 specforce console`) para imprimir cada decisión del actualizador en la salida de error estándar. El estado del actualizador por usuario se guarda en `~/.specforce/state.json`.

## Primeros pasos

A continuación, las skills se mencionan por su nombre (por ejemplo `spf.constitution`). Cómo llamarlas explícitamente depende de tu herramienta: consulta [Invocar las skills](supported-tools.md#invocar-las-skills). Una petición en lenguaje natural que nombre la skill funciona en todas partes.

### 1. Inicializa tu proyecto

Ve al directorio de tu proyecto (nuevo o existente) e inicializa Specforce:

```bash
specforce init
```

Elige tus herramientas de IA en la lista interactiva, o pasa sus claves directamente (por ejemplo `specforce init claude cursor`). Esto crea la carpeta `.specforce/`, las reglas del `AGENTS.md` raíz y las skills de cada herramienta seleccionada.

### 2. Define la Constitución

La **Constitución** es el conjunto fundacional de reglas de tu proyecto (principios, arquitectura, UI/UX, seguridad, ingeniería, gobernanza). En lugar de contaminar el contexto de tu agente con decenas de reglas sueltas, Specforce las centraliza de forma estructurada.

Inicia tu agente de IA y pide, en su chat:

```text
Usa spf.constitution para configurar la constitución: <la descripción de tu proyecto y su stack tecnológico>
```

El agente hará algunas preguntas para alinear el proyecto y creará los documentos fundacionales en la carpeta `.specforce/docs/`.

### 3. Crea tu primera especificación

Todo cambio de código en Specforce requiere una especificación. Para añadir una nueva funcionalidad, pide a tu agente:

```text
Usa spf.spec para planificar: <describe la funcionalidad que quieres construir>
```

El agente crea los artefactos requeridos por el tamaño de la especificación; una especificación `medium`, el tamaño por defecto, necesita `requirements.md` y `tasks.md` (consulta [Artefactos](artifacts.md)).

---

## Flujo de trabajo: el ciclo de vida

El flujo de trabajo de Specforce Kit se basa en el **Desarrollo Orientado por Especificaciones (SDD)**. En la era de la IA, la implementación se ha convertido en un commodity, pero el diseño arquitectónico sin validación es un gran riesgo.

El SDD separa el **diseño** de la **implementación**, garantizando que no se genere código sobre la marcha ("vibe coding") sin un contrato formal. El ciclo de desarrollo sigue un camino determinista. Para flujos listos para usar, consejos y qué hacer cuando un paso queda bloqueado, consulta [Flujos de trabajo y consejos](workflows.md).

### 1. Fundamento: la Constitución (`spf.constitution`)
Todo empieza con la **Constitución**. Antes de escribir código, el proyecto necesita reglas claras de arquitectura, UI/UX, seguridad e ingeniería.
- El agente interactúa contigo para entender el alcance del proyecto, o lee el código existente de un proyecto que estás adoptando.
- Documenta estas decisiones en `.specforce/docs/`, que sirven como las "leyes" para todas las implementaciones futuras.
- **Ventaja:** el agente solo carga los documentos que necesita, ahorrando tokens y evitando alucinaciones.

### 2. Exploración: Discovery (`spf.discovery`)
Usa discovery para hacer una lluvia de ideas, investigar un bug o encontrar una causa raíz antes de comprometerte con un plan.
- El código del producto queda en solo lectura; la evidencia desechable (tests, scripts, capturas de pantalla) vive en `.specforce/scratch/`.
- Cuando estés listo, pide al agente que formalice el resultado en un `proposal.md` que la planificación reutiliza.

### 3. El contrato: la Especificación (`spf.spec`)
Toda nueva funcionalidad, corrección o refactorización empieza aquí.
- Explicas lo que necesitas. El agente escribe los artefactos requeridos por el tamaño de la especificación:
  - **Requisitos** (el "qué")
  - **Diseño** (el "cómo", a nivel arquitectónico)
  - **Tareas** (la hoja de ruta de ejecución atómica)
  - **Deltas de módulo** (cómo afecta el cambio a las living specs de módulo)
- *Acción humana:* debes revisar y aprobar esta especificación. Es tu contrato con la máquina.

### 4. Ejecución: la Implementación (`spf.implement`)
Con la especificación aprobada, el agente entra en modo de ejecución.
- El agente sigue estrictamente la lista de tareas generada en el paso anterior y luego ejecuta un paso de QA independiente que registra un veredicto para cada criterio de aceptación.
- Los [Hooks de validación](configuration.md) garantizan que el agente no se salte tests, linters ni pasos de seguridad antes de marcar una tarea como "finalizada".
- *Acción humana:* usa el comando `specforce console` en otra terminal para seguir el progreso del agente en tiempo real.

### 5. Cierre: el Archivado (`spf.archive`)
Una vez que la funcionalidad está completamente programada y verificada, hay que cerrar la especificación.
- El agente aplica los deltas de módulo a las living specs de `.specforce/docs/modules/`, actualiza la Constitución si surgieron nuevas decisiones sistémicas y mueve la especificación a `.specforce/archive/`.
- El archivado se rechaza mientras haya tareas pendientes, el delta no esté aplicado o el veredicto de QA sea FAILED.
- Puedes personalizar estas reglas en `config.yaml` para incluir pasos como actualizar herramientas de seguimiento externas o tareas de limpieza específicas.

---

## Rol humano vs. rol de la máquina

En Specforce, los roles están claros:

- **El Arquitecto Humano (Gobernanza):** marca la dirección del proyecto, aprueba las especificaciones y se asegura de que el sistema en su conjunto tenga sentido para el negocio.
- **El Agente de IA (Ejecutor):** ejecuta la hoja de ruta aprobada, generando código de alta calidad, tests y documentación interna de forma eficiente y rápida, respetando siempre la Constitución.

---

Para explorar los comandos de la terminal, consulta la [Referencia de la CLI](cli.md).
Conoce más sobre los documentos generados en este proceso en la página de [Artefactos](artifacts.md).
Sigue flujos paso a paso y consejos en [Flujos de trabajo y consejos](workflows.md).
