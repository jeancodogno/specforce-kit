[ [English](../en/workflows.md) | [Português](../pt/workflows.md) | [Español](workflows.md) ]

# Flujos de trabajo y consejos

Esta página explica algunos flujos de trabajo de Specforce y consejos para usarlo.

Las skills se mencionan por su nombre (por ejemplo [`spf.spec`](supported-tools.md#invocar-las-skills)). Cómo llamarlas explícitamente en tu herramienta se muestra en [Invocar las skills](supported-tools.md#invocar-las-skills); una petición en lenguaje natural que nombre la skill funciona en todas partes.

**Cómo leer los flujos**

- 💻 **Terminal:** un comando que **tú** escribes en tu propia terminal. Aquí solo aparecen unos pocos comandos: `specforce init`, `specforce console`, `specforce spec list` y un par de comandos `git`. Todo lo demás (estado, comprobación del delta, resultados de QA) lo comprueba el agente y te lo informa en el chat.
- 💬 **Chat del agente:** una petición que **tú** escribes en el chat de tu agente de IA. Sustituye las partes entre `<corchetes angulares>`.
- **Detrás de escena** (contraído): lo que el agente ejecuta por su cuenta. Nunca necesitas ejecutar esos comandos; se listan para que puedas seguir lo que ocurre. Todos están descritos en la [Referencia de la CLI](cli.md).
- `<slug>` es el nombre de la especificación, por ejemplo `csv-export`. `specforce spec list` lo muestra.

## Flujos

### 1. Prototipar en discovery antes de implementar

**Cuándo usarla:** la idea todavía es vaga, o quieres una prueba (un test rápido, un script, una captura de pantalla) antes de comprometerte con un plan.

1. 💬 **Chat del agente:**
   ```text
   Usa spf.discovery para explorar: en nuestro SaaS de facturación, permitir que el cliente exporte todas sus facturas a CSV desde el panel. Las cuentas grandes tienen 200 mil facturas, así que demuestra con un test desechable que generar la exportación en streaming mantiene estable el uso de memoria antes de planificar nada.
   ```
2. 💬 **Chat del agente:** una vez que estén de acuerdo con el enfoque:
   ```text
   Formaliza esto en una propuesta.
   ```
3. 💬 **Chat del agente:**
   ```text
   Usa spf.spec para planificar csv-export a partir de la propuesta.
   ```

<details><summary>Detrás de escena</summary>

- Discovery nunca edita, renombra ni elimina tus archivos existentes. La evidencia (tests, scripts, logs, capturas de pantalla) va a `.specforce/scratch/`, que se excluye a sí misma del control de versiones. Si una hipótesis solo puede probarse dentro del árbol de código fuente, el agente crea un archivo **nuevo** cuyo nombre contiene `spf-probe`, lo registra y lo elimina antes de que termine la sesión, comprobando `git status`.
- Formalizar ejecuta `specforce spec init <slug>` y escribe `proposal.md` (problema, hallazgos, decisiones tomadas, alternativas descartadas, preguntas abiertas, tipo y tamaño sugeridos, rutas de la evidencia).
- `spf.spec` repite las decisiones que ya están en la propuesta para una única confirmación y solo hace las preguntas abiertas.

</details>

**Si algo sale mal:**
- Discovery no te empujará a una especificación por su cuenta; pide la propuesta explícitamente cuando estés listo.
- Un archivo `spf-probe` sobrante aparece en `git status`: pide al agente que lo elimine.
- La evidencia de scratch de una sesión anterior está desactualizada; pide al agente que la vuelva a ejecutar antes de confiar en ella.

### 2. Corregir un bug con un test que falla primero

**Cuándo usarla:** algo está roto y quieres que la corrección quede demostrada por un test de regresión.

1. 💬 **Chat del agente:**
   ```text
   Usa spf.discovery para investigar este bug: en el panel de facturas, el filtro "mes pasado" deja fuera las facturas emitidas el último día del mes para clientes en UTC-3. No aparece ningún error, la lista simplemente llega incompleta.
   ```
2. 💬 **Chat del agente:** una vez confirmada la causa raíz:
   ```text
   Formaliza esto en una especificación de bug.
   ```
3. 💬 **Chat del agente:** revisa y aprueba la especificación, y luego:
   ```text
   Usa spf.implement para implementar invoice-date-filter.
   ```

<details><summary>Detrás de escena</summary>

- La especificación de bug se crea con `specforce spec init <slug> --type bug`. Sus requisitos contienen los pasos de reproducción, la evidencia y la causa raíz.
- La hoja de ruta empieza con una tarea `[RED]` que añade un test de regresión que falla antes de la corrección; ninguna tarea de corrección va antes que ella.

</details>

**Si algo sale mal:**
- `spf.spec` se detiene cuando la causa raíz no está establecida: vuelve al paso 1 en lugar de dejar que adivine.
- La tarea `[RED]` se finaliza mientras los tests fallan: es lo esperado; solo se omiten los hooks por tarea (consulta [Un hook fallido bloquea una tarea](#un-hook-fallido-bloquea-una-tarea)).

### 3. Hacer un cambio pequeño

**Cuándo usarla:** un cambio acotado (una corrección en un solo lugar, un ajuste de texto, una opción de configuración) que no necesita requisitos ni diseño.

1. 💬 **Chat del agente:**
   ```text
   Usa spf.spec para planificar: añadir la columna "Vencimiento" a la tabla de facturas del panel.
   ```
2. 💬 **Chat del agente:** revisa la lista de tareas, y luego:
   ```text
   Usa spf.implement para implementar invoice-due-date-column.
   ```

<details><summary>Detrás de escena</summary>

- No necesitas decir que el cambio es pequeño: el agente define el tamaño durante la planificación y ejecuta `specforce spec init <slug> --size small`: solo se requiere `tasks.md` y el delta de módulo es `exempt` (consulta [Tamaños](artifacts.md#3-tamaños)).

</details>

**Si algo sale mal:**
- El cambio resulta más grande de lo esperado: sigue el [flujo 4](#4-el-alcance-creció-redimensionar-la-especificación).

### 4. El alcance creció: redimensionar la especificación

**Cuándo usarla:** durante la planificación o la implementación, el cambio termina necesitando más (o menos) de lo planificado.

1. 💬 **Chat del agente:** basta con describir el nuevo alcance:
   ```text
   Usa spf.spec para actualizar invoice-due-date-column: además de la columna, enviar al cliente un correo de recordatorio 3 días antes del vencimiento.
   ```
2. 💬 **Chat del agente:** después de aprobar los nuevos artefactos:
   ```text
   Usa spf.implement para continuar invoice-due-date-column.
   ```

<details><summary>Detrás de escena</summary>

- No necesitas pedir el redimensionamiento: el agente nota que el alcance cambió y ejecuta `specforce spec resize <slug> --size <nuevo tamaño>` por su cuenta. Promover marca los artefactos recién requeridos como faltantes y bloquea la implementación hasta que existan; degradar relaja la validación y conserva los archivos existentes como referencia.
- Luego el agente escribe los artefactos que exige el nuevo tamaño y, a partir de `medium`, escribe y sella el delta de módulo, o declara `delta.impact: none`.

</details>

**Si algo sale mal:**
- La implementación se niega a continuar: consulta [Redimensionar durante la implementación](#redimensionar-durante-la-implementación).

### 5. Cambio después de implementar, antes de archivar

**Cuándo usarla:** la implementación está terminada (o casi), y quieres algo distinto antes de cerrar la especificación.

1. 💬 **Chat del agente:**
   ```text
   Usa spf.spec para actualizar invoice-due-date-column: enviar el recordatorio 5 días antes del vencimiento en lugar de 3. Mantén coherentes los requisitos, el diseño, el delta y las tareas.
   ```
2. 💬 **Chat del agente:** revisa y aprueba la especificación actualizada, y luego:
   ```text
   Usa spf.implement para implementar solo las tareas nuevas o modificadas de invoice-due-date-column.
   ```

<details><summary>Detrás de escena</summary>

- Una especificación aprobada es un contrato: el agente actualiza juntos todos los artefactos afectados, nunca solo `tasks.md`, y añade tareas nuevas en lugar de reescribir las ya finalizadas.
- `spf.implement` ejecuta después un nuevo paso de QA independiente y registra un nuevo veredicto.

</details>

**Si algo sale mal:**
- No edites el código ni `tasks.md` a mano para "arreglarlo y ya": consulta [Editar una especificación aprobada](#editar-una-especificación-aprobada).

### 6. Cambio después de archivar: especificación de seguimiento

**Cuándo usarla:** la especificación ya está archivada y quieres cambiar o ampliar lo que entregó.

1. 💬 **Chat del agente:**
   ```text
   Usa spf.spec para planificar un seguimiento de csv-export: permitir que el cliente elija un rango de fechas antes de exportar.
   ```
2. 💬 **Chat del agente:** revisa y aprueba, y luego usa el nombre que el agente dio a la nueva especificación:
   ```text
   Usa spf.implement para implementar csv-export-date-range.
   ```

<details><summary>Detrás de escena</summary>

- Las especificaciones archivadas son historial y nunca se reabren. La nueva especificación tiene su propia carpeta.
- Cuando el cambio altera reglas que ya están en una living spec de módulo, su delta las lleva en `## MODIFIED` (o `## REMOVED`) con el bloque nuevo completo (consulta [Deltas de módulo](artifacts.md#6-deltas-de-módulo)).
- El resumen escrito por `spf.archive` lista los próximos pasos sugeridos y las especificaciones de seguimiento: un buen punto de partida para la petición.

</details>

**Si algo sale mal:**
- La comprobación del delta informa un conflicto: consulta [Conflicto de base del delta](#conflicto-de-base-del-delta).

### 7. Falló el QA

**Cuándo usarla:** el paso de QA independiente informó criterios `FAILED` y el agente te lo escaló.

1. 💬 **Chat del agente:** lee el informe del agente: ya enumera cada criterio fallido con su evidencia. Si el código está mal:
   ```text
   Usa spf.implement para corregir los criterios de QA fallidos de invoice-due-date-column y volver a ejecutar el QA.
   ```
   Si el propio criterio está mal, sigue en su lugar el [flujo 5](#5-cambio-después-de-implementar-antes-de-archivar).

<details><summary>Detrás de escena</summary>

- El agente ya ejecutó hasta dos ciclos de corrección (un implementador corrige, un tester nuevo vuelve a comprobar) antes de escalar, así que lo que te llega es solo lo que no pudo corregir por su cuenta.
- El veredicto se guarda en `spec.yaml` con `specforce implementation qa record <slug>` después de cada ejecución de QA.

</details>

**Si algo sale mal:**
- El archivado se rechaza mientras el veredicto guardado sea `FAILED`, incluso con `--force`: consulta [Archivado rechazado: QA FAILED](#archivado-rechazado-qa-failed).

### 8. Especificaciones en paralelo en varios worktrees

**Cuándo usarla:** quieres dos agentes trabajando en dos especificaciones al mismo tiempo, cada uno en su propia rama.

1. 💻 **Terminal:** crea un worktree para la segunda especificación:
   ```bash
   git worktree add ../myapp-csv-export -b feature/csv-export
   ```
2. 💬 **Chat del agente:** inicia una segunda sesión del agente dentro de `../myapp-csv-export`, y luego:
   ```text
   Usa spf.spec para planificar: permitir que el cliente exporte todas sus facturas a CSV desde el panel.
   ```
3. 💻 **Terminal:** sigue ambas desde tu checkout principal:
   ```bash
   specforce console
   ```

<details><summary>Detrás de escena</summary>

- Las especificaciones de otros worktrees son de solo lectura en tu checkout: sus tareas solo pueden cambiarse desde su propio worktree. La Constitución siempre proviene de tu rama actual (consulta [Soporte para Git Worktrees](git-worktrees.md)).
- Cada especificación sella su delta. Cuando ambas cambian el mismo bloque de módulo, la segunda en archivarse ve un conflicto y debe rehacer su delta.

</details>

**Si algo sale mal:**
- El agente en el nuevo worktree no ve las skills: si las carpetas de tus herramientas no están en el repositorio, ejecuta `specforce init` dentro de ese worktree.
- Dos especificaciones tocan la misma regla de módulo: consulta [Conflicto de base del delta](#conflicto-de-base-del-delta).

### 9. Adoptar Specforce en un proyecto existente

**Cuándo usarla:** el proyecto ya tiene código y quieres que Specforce lo documente antes de la primera especificación.

1. 💻 **Terminal:** en la raíz del proyecto:
   ```bash
   specforce init
   ```
2. 💬 **Chat del agente:**
   ```text
   Usa spf.constitution para documentar este proyecto existente: lee el código y redacta la constitución a partir de lo que ya existe.
   ```
3. 💬 **Chat del agente:** revisa los documentos de `.specforce/docs/`, corrige lo que esté mal, y luego planifica el primer cambio:
   ```text
   Usa spf.spec para planificar: añadir un buscador a la tabla de facturas que filtre por nombre del cliente.
   ```

<details><summary>Detrás de escena</summary>

- `specforce init` escribe sus reglas dentro de un bloque marcado del `AGENTS.md` raíz; tu contenido existente fuera del bloque se conserva.
- `spf.constitution` lee el código, pregunta por las decisiones que no puede deducir y escribe los seis documentos más `modules/`.

</details>

**Si algo sale mal:**
- La constitución describe el código tal como es, incluidos hábitos que quieres abandonar: dilo en el chat; registra la regla objetivo, no el accidente.

### 10. Reanudar una implementación interrumpida

**Cuándo usarla:** la sesión se cayó, cerraste la terminal o el agente se detuvo en medio de una hoja de ruta.

1. 💬 **Chat del agente:**
   ```text
   Usa spf.implement para verificar el estado de la implementación de invoice-due-date-column y reanudarla.
   ```

<details><summary>Detrás de escena</summary>

- El agente lee los estados de las tareas y el repositorio. Las tareas que quedaron en `in-progress` fueron interrumpidas: inspecciona el código afectado y las vuelve a ejecutar. Las tareas finalizadas no se rehacen, y los cambios sin commit se conservan.

</details>

**Si algo sale mal:**
- No descartes los cambios sin commit antes de reanudar: el agente los usa para continuar de forma segura.

### 11. Cambiar de herramienta de IA a mitad del proyecto

**Cuándo usarla:** te pasas a otra herramienta de IA (o añades una) en medio de un proyecto o de una especificación.

1. 💻 **Terminal:** lista todas las herramientas que quieres mantener, usando las claves de [Herramientas soportadas](supported-tools.md):
   ```bash
   specforce init claude cursor
   ```
   ⚠️ Cada herramienta instalada que dejes fuera se desinstala y **se elimina su carpeta entera** (por ejemplo `.claude/`, con tus ajustes y reglas), sin pedir confirmación. Haz commit o una copia de esa carpeta antes si contiene algo tuyo.
2. 💬 **Chat del agente:** en la nueva herramienta:
   ```text
   Usa spf.implement para verificar el estado de la implementación de invoice-due-date-column y reanudarla.
   ```

<details><summary>Detrás de escena</summary>

- La Constitución, las especificaciones y su progreso viven en `.specforce/` y no dependen de ninguna herramienta, así que no se pierde nada.
- `specforce init` instala las skills de las herramientas listadas y elimina la carpeta de cada herramienta instalada que dejaste fuera. Sin argumentos, el selector interactivo muestra un aviso antes de eliminar nada.

</details>

**Si algo sale mal:**
- La carpeta de una herramienta desapareció: dejaste fuera su clave. Restaura la carpeta desde el control de versiones si tenía archivos tuyos, y vuelve a ejecutar `specforce init` incluyendo la clave, o sin argumentos para usar el selector.

### 12. Imponer TDD mediante la configuración

**Cuándo usarla:** quieres que cada hoja de ruta sea test-first y que cada tarea pase tus comprobaciones antes de marcarse como finalizada.

1. 💬 **Chat del agente:** (o edita el archivo tú mismo)
   ```text
   En .specforce/config.yaml, añade una instrucción de tasks que exija tareas de test [RED] seguidas de tareas [GREEN] en la misma fase, y los hooks on_task_finished "make lint" y "make test".
   ```
   El resultado se ve así:
   ```yaml
   instructions:
     tasks:
       - "Each phase must contain [RED] tasks (tests for happy path and edge cases) followed by [GREEN] tasks (the implementation)."
   hooks:
     on_task_finished:
       - "make lint"
       - "make test"
   ```
2. 💬 **Chat del agente:**
   ```text
   Usa spf.spec para planificar: permitir que el cliente descargue cada factura en PDF desde el panel.
   ```

<details><summary>Detrás de escena</summary>

- Las tareas llevan el token exacto `[RED]` en su título. Cuando todas las tareas finalizadas en un lote son `[RED]`, solo se omiten los hooks `on_task_finished` (con un aviso); los hooks de fase y de todas las tareas se siguen ejecutando.

</details>

**Si algo sale mal:**
- Una entrada de hook con `|` o `&&` nunca funciona: consulta [Un hook fallido bloquea una tarea](#un-hook-fallido-bloquea-una-tarea).

## Situaciones comunes

### Archivado rechazado: hay tareas pendientes

**Síntoma:** el archivado se detiene y lista tareas que no están finalizadas.

**Por qué:** una especificación solo se archiva cuando todas sus tareas están finalizadas. La skill de archivado se niega incluso si le pides que archive de todos modos; `--force` solo omite esta comprobación y nunca las del delta ni del QA.

**Qué hacer:** el agente ya lista las tareas sin terminar. Pide `Usa spf.implement para terminar <slug>`. Si una tarea ya no es necesaria, pide en su lugar a `spf.spec` que la quite de la especificación.

### Archivado rechazado: delta faltante o pendiente

**Síntoma:** el archivado se rechaza con un estado de delta `missing`, `pending` o `invalid`.

**Por qué:** a partir de `medium`, las living specs de módulo deben reflejar la especificación antes de archivar: el delta debe estar `applied` (o declarado `none`). `missing` significa que no se escribió ningún delta; `pending` significa que todavía no se aplicó a los módulos.

**Qué hacer:** el agente informa el estado del delta. Para `missing`, pide `Usa spf.spec para escribir el delta de módulo de <slug>, o declarar que no hay impacto en módulos`. Para `pending`, pide `Usa spf.archive para aplicar el delta y archivar <slug>`. Para `invalid`, pide a `spf.spec` que corrija los problemas informados.

### Conflicto de base del delta

**Síntoma:** el agente informa un conflicto de delta en un bloque de módulo, y el archivado deja el módulo sin tocar.

**Por qué:** después de que se selló tu delta, otra especificación (a menudo de una rama paralela) cambió el mismo bloque de módulo. Aplicar el tuyo sobrescribiría ese cambio sin avisar.

**Qué hacer:** pide `Usa spf.spec para rehacer el delta de <slug> contra el módulo actual y volver a sellarlo`. El agente vuelve a comprobar el delta antes de archivar.

### Archivado rechazado: QA FAILED

**Síntoma:** el archivado se detiene y lista criterios de QA fallidos.

**Por qué:** el veredicto de QA guardado es `FAILED`. El archivado nunca avanza sobre un veredicto fallido, ni siquiera con `--force`.

**Qué hacer:** lee los criterios fallidos que lista el agente, y luego sigue el [flujo 7](#7-falló-el-qa).

### Aviso de archivado: no hay informe de QA

**Síntoma:** el archivado continúa pero avisa de que no existe un informe de QA.

**Por qué:** no se registró ningún veredicto para la especificación, normalmente porque la implementación se hizo o se terminó fuera de `spf.implement`, o porque la sesión terminó antes del QA.

**Qué hacer:** si quieres la evidencia, detente y pide `Usa spf.implement para ejecutar el paso de QA independiente de <slug>`, y luego vuelve a archivar. De lo contrario, el aviso es solo informativo.

### Un hook fallido bloquea una tarea

**Síntoma:** el agente no puede marcar una tarea como finalizada y muestra la salida de un hook.

**Por qué:** un hook de `.specforce/config.yaml` terminó con un código distinto de cero. Los hooks se ejecutan sin shell (cada entrada se divide por espacios), en paralelo, así que `|`, `&&` y las redirecciones nunca se interpretan. Solo un lote compuesto enteramente por tareas `[RED]` omite los hooks `on_task_finished`.

**Qué hacer:** pide `Corrige la causa del hook que falla; no lo omitas ni lo elimines`. Si la propia entrada del hook está mal, edita `.specforce/config.yaml` para que llame a un script o a un target de `make` (consulta [Configuración](configuration.md)).

### El agente quiere programar sin especificación

**Síntoma:** pides un cambio y el agente se niega a editar el código, o se detiene en medio de la implementación pidiendo una actualización de la especificación.

**Por qué:** es intencionado. El código solo se cambia a través de una especificación aprobada, y un cambio de alcance o de arquitectura debe pasar antes por la planificación.

**Qué hacer:** pide `Usa spf.spec para planificar: <el cambio>`, o `Usa spf.discovery para explorar: <la idea>` cuando todavía no esté claro.

### El worker devuelve una pregunta BLOCKED

**Síntoma:** durante la implementación, el agente se detiene y te hace una pregunta con opciones y una recomendación.

**Por qué:** el worker que ejecuta un lote no puede preguntarte directamente. Cuando necesita una decisión, se detiene con un informe BLOCKED y la sesión principal te lo transmite.

**Qué hacer:** responde en el chat, eligiendo una opción o la recomendación. Si tu respuesta cambia el alcance, espera que el agente actualice primero la especificación con `spf.spec`.

### Redimensionar durante la implementación

**Síntoma:** después de un redimensionamiento, la implementación se niega a continuar y lista artefactos faltantes.

**Por qué:** promover una especificación (por ejemplo a `large`) marca los artefactos recién requeridos como faltantes, y la implementación queda bloqueada hasta que existan. Degradar nunca elimina archivos.

**Qué hacer:** pide `Usa spf.spec para escribir los artefactos que faltan de <slug>`. Después de aprobarlos, pide `Usa spf.implement para continuar <slug>`.

### Editar una especificación aprobada

**Síntoma:** quieres cambiar los requisitos o las tareas de una especificación que ya aprobaste.

**Por qué:** una especificación aprobada es un contrato. Cambiarla significa actualizar juntos todos los artefactos relacionados (requisitos, diseño, delta, tareas) y volver a aprobarla; editar solo `tasks.md` deja la especificación incoherente.

**Qué hacer:** pide `Usa spf.spec para actualizar <slug>: <el cambio>. Mantén coherentes todos los artefactos`, revisa, aprueba, y luego sigue el [flujo 5](#5-cambio-después-de-implementar-antes-de-archivar).
