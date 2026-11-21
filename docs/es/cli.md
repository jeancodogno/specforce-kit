[ [English](../en/cli.md) | [Português](../pt/cli.md) | [Español](cli.md) ]

# Interfaz de línea de comandos (CLI)

Esta página lista todos los comandos y flags de la CLI `specforce`, agrupados por área. Muchos de estos comandos los ejecuta el agente de IA mientras sigue una skill (consulta [Invocar las skills](supported-tools.md#invocar-las-skills)); se documentan aquí para que sepas qué existe y qué hace cada uno.

Los comandos que reciben un `<slug>` aceptan el nombre de la carpeta de la especificación (por ejemplo `20261121-0035-csv-export`) o solo su parte legible (`csv-export`). Los flags marcados como `--json` imprimen JSON legible por máquinas, pensado para automatizaciones y agentes.

## Global

### `specforce init [agents...]`
Inicializa Specforce en el directorio actual: crea la carpeta `.specforce/` (con un `config.yaml` por defecto), el bloque gestionado del `AGENTS.md` raíz y las skills de las herramientas seleccionadas. Sin argumentos, un selector interactivo te permite elegir las herramientas. Con argumentos, las claves de herramienta indicadas (`claude`, `qwen`, `kimi-code`, `open-code`, `kilo-code`, `codex`, `antigravity`, `cursor`) se configuran inmediatamente. Vuelve a ejecutarlo en un proyecto ya inicializado para añadir o quitar herramientas y refrescar las skills instaladas. Con argumentos, cada herramienta instalada que dejes fuera se quita eliminando su carpeta entera (por ejemplo `.claude/`, incluidos tus propios ajustes), sin confirmación; el selector interactivo avisa antes de eliminar. Consulta [Herramientas soportadas](supported-tools.md).

### `specforce console`
Lanza la **Consola de Specforce**, una interfaz de terminal interactiva para explorar las especificaciones, la Constitución y las tareas de implementación, y para seguir el progreso en tiempo real (incluidas las especificaciones de [Git worktrees](git-worktrees.md) vinculados). Déjala abierta en un panel de terminal aparte mientras tu agente trabaja.

---

## Flags globales y entorno

- **`-v`, `--version`**: imprime la versión instalada (`specforce version <x.y.z>`).
- **`-h`, `--help`**: muestra la ayuda de cualquier comando o grupo.
- **`--dev`**: usa la carpeta local `kit/` en lugar del kit incrustado en el binario. Pensado solo para quienes contribuyen al kit.
- **`DEBUG=1`**: variable de entorno que escribe diagnósticos en la salida de error estándar, incluida cada decisión de actualización automática (consulta [Actualizaciones](getting-started.md#actualizaciones)).

### `specforce completion <shell>`
Genera el script de autocompletado de shell para `bash`, `zsh`, `fish` o `powershell`. Ejecuta `specforce completion <shell> --help` para ver los pasos de instalación de cada shell.

---

## Constitución

### `specforce constitution status [--json]`
Muestra el estado de completitud de los documentos de la Constitución del proyecto en `.specforce/docs/` (principles, architecture, ui-ux, security, engineering, governance y modules).

### `specforce constitution artifact <name> [--json]`
Muestra los detalles de un artefacto de la Constitución: su propósito, su plantilla y sus instrucciones. Los agentes lo usan para obtener las reglas correctas cuando escriben o actualizan ese documento.

---

## Especificación

### `specforce spec init <slug> [--type feature|bug] [--size small|medium|large|complex] [--json]`
Crea una nueva carpeta de especificación `.specforce/specs/<timestamp>-<slug>/` con sus metadatos `spec.yaml`.
- `--type`: `feature` (por defecto) o `bug`. Las especificaciones de bug usan las variantes de bug de los requisitos y del diseño.
- `--size`: `small`, `medium` (por defecto), `large` o `complex`. El tamaño decide qué artefactos son obligatorios (consulta [Artefactos](artifacts.md)).

### `specforce spec list [--json]`
Lista todas las especificaciones activas.

### `specforce spec status <slug> [--json]`
Muestra la completitud de una especificación: qué artefactos requiere su tamaño, cuáles existen, los errores de validación y los avisos no bloqueantes de las tareas.

### `specforce spec artifact <name> [--json]`
Muestra las instrucciones y la plantilla de un artefacto de especificación (por ejemplo `requirements`, `design`, `tasks`, `bug-requirements`, `bug-design`), incluidas las instrucciones del proyecto de [`config.yaml`](configuration.md). Los agentes lo usan antes de escribir cada artefacto.

### `specforce spec resize <slug> --size <size> [--json]`
Cambia el tamaño de una especificación existente. Promoverla marca los artefactos recién requeridos como faltantes y bloquea la implementación hasta que existan; degradarla relaja la validación y conserva en disco los archivos existentes.

### `specforce spec audit <slug> [flags]`
Actualiza el estado de la auditoría de refinamiento guardado en `spec.yaml`, que usa la skill de planificación mientras comprueba la coherencia entre artefactos.
- `--valid` / `--invalid`: marca la especificación como válida o inválida.
- `--error <msg>`: añade un mensaje de error de coherencia.
- `--clear`: borra todos los errores de coherencia.
- `--iteration <n>`: establece el número de la iteración de refinamiento actual.
- `--json`: salida legible por máquinas.

### `specforce spec archive <slug> [--force]`
Finaliza una especificación completada y mueve su carpeta a `.specforce/archive/`. Se niega cuando:
- todavía hay tareas pendientes;
- una especificación `medium`, `large` o `complex` tiene un delta que no está `applied` (o declarado `none`); las especificaciones `small` informan `exempt`;
- el informe de QA guardado es `FAILED` (se listan los criterios fallidos).

Una especificación sin informe de QA se archiva con un aviso visible. `--force` anula solo la comprobación de tareas pendientes; nunca se salta la comprobación del delta ni un veredicto de QA `FAILED`. Un rechazo deja el sistema de archivos intacto.

---

## Delta

El delta describe los cambios que una especificación hace en las living specs de módulo de `.specforce/docs/modules/`. El archivado exige que el delta esté aplicado, sea none o esté exento.

### `specforce spec delta check <slug> [--json]`
Informa si el delta de la especificación (los cambios propuestos a las living specs de módulo) se ha aplicado. Solo lectura. El estado del informe es uno de:
- `none`: la especificación declara `delta.impact: none` (sin impacto en módulos).
- `exempt`: una especificación `small` sin delta; no lo necesita.
- `missing`: se requiere un delta pero no se ha escrito.
- `pending`: existe un delta pero sus cambios todavía no se reflejan en los módulos.
- `applied`: los módulos ya coinciden con el delta.
- `invalid`: el delta está mal formado o entra en conflicto con los módulos actuales.

### `specforce spec delta seal <slug> [--json]`
Sella el delta de la especificación registrando las huellas (fingerprints) de los bloques de módulo en los que se basa, de modo que las ediciones manuales posteriores de esos bloques se detecten como conflictos de base.

---

## Implementación y QA

### `specforce implementation status <slug> [--json]`
Muestra la hoja de ruta de implementación de una especificación (a partir de `tasks.md`): si está lista, los artefactos que faltan para su tamaño, cada tarea con su estado y sus restricciones, las mitigaciones de la hoja de ruta, las instrucciones `implementation` del proyecto y, cuando existe, el estado de QA guardado.

### `specforce implementation update <slug> --task <id>[,<id>...] --status <status>`
Actualiza el estado de una o más tareas. `--task` acepta IDs separados por comas o puede repetirse (`--task T1.1 --task T1.2`). `--status` es `pending`, `in-progress` o `finished`.

Cuando el estado es `finished`, Specforce ejecuta primero los [hooks de validación](configuration.md#3-hooks-de-validación-quality-gates) configurados, una vez por lote y sin duplicados: `on_task_finished`, más `on_phase_finished` cuando una tarea es la última de su fase y `on_all_tasks_finished` cuando es la última de la hoja de ruta. Un hook que falla bloquea la actualización. Cuando todas las tareas del lote llevan el token `[RED]` en su título, solo se omiten los hooks `on_task_finished`, con un aviso que los nombra; los hooks de fase y de todas las tareas se siguen ejecutando.

### `specforce implementation qa record <slug>`
Registra el informe de QA independiente de una especificación. Lee el informe como JSON desde la entrada estándar (criterios con escenario, evidencia y veredicto `PASSED`, `FAILED` o `NOT-VERIFIABLE`, más los hallazgos fuera de alcance), lo valida, deriva el estado global (`FAILED` cuando algún criterio falló y `PASSED` en caso contrario) y reemplaza el informe anterior en `spec.yaml`. Una entrada inválida se rechaza sin cambiar el informe guardado. La skill de implementación lo usa después del paso de QA.

---

## Archivado

### `specforce archive instructions`
Muestra el protocolo de archivado completo: la descripción de cada documento de la Constitución, las reglas centrales del archivado (aplicar el delta a las living specs de módulo, actualizar la Constitución) y las instrucciones `archive` del proyecto de `config.yaml`. La skill de archivado lo sigue antes de ejecutar `specforce spec archive`.
