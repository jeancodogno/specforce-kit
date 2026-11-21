[ [English](../en/artifacts.md) | [Português](../pt/artifacts.md) | [Español](artifacts.md) ]

# Artefactos de Specforce

Para orquestar agentes de IA sin acumular contexto innecesario, Specforce divide el conocimiento del proyecto en documentos enfocados llamados **Artefactos**. Esta página explica cada uno de ellos y las reglas del ciclo de vida que deciden qué produce el agente y por qué un paso puede quedar bloqueado. Las skills se mencionan por su nombre (por ejemplo [`spf.spec`](supported-tools.md#invocar-las-skills)).

## 1. Constitución (`.specforce/docs/`)

La **Constitución** vive en la carpeta `.specforce/docs/` de tu proyecto. Contiene las reglas duraderas del sistema y se crea y mantiene con `spf.constitution`:

- **`principles.md`**: valores fundamentales, filosofía y axiomas técnicos.
- **`architecture.md`**: límites del sistema, dirección de las dependencias y topología de persistencia.
- **`ui-ux.md`**: dirección visual, patrones de interacción y reglas estéticas (cuando aplica).
- **`security.md`**: autorización, roles, permisos y protección de datos.
- **`engineering.md`**: estándares de código, estrategia de testing y pautas de refactorización.
- **`governance.md`**: reglas del ciclo de vida, responsabilidades y límites de la IA.
- **`modules/<module>.md`**: living specs de cada dominio (reglas de negocio, casos de uso, superficies públicas), actualizadas en el momento del archivado a partir de los deltas de las especificaciones.

Los documentos no se cargan todos a la vez: el agente lee los relevantes para la tarea en curso.

## 2. El paquete de especificación (`.specforce/specs/<timestamp>-<slug>/`)

Todo cambio empieza con una especificación creada por `spf.spec`. Su carpeta puede contener:

- **`spec.yaml`**: metadatos gestionados por la CLI: tipo, tamaño, declaración y sello del delta, auditoría de refinamiento, registros de tiempo de las tareas y el informe de QA.
- **`proposal.md`** *(opcional)*: el resultado formalizado de una sesión de `spf.discovery` (problema, hallazgos verificados, decisiones, preguntas abiertas, tipo y tamaño sugeridos). La planificación lo reutiliza en lugar de volver a preguntar.
- **`requirements.md`**: el "qué" y el "por qué": historias de usuario, criterios de aceptación BDD (GIVEN/WHEN/THEN) y restricciones de negocio. Sin detalles técnicos.
- **`design.md`**: el "cómo": arquitectura, contratos y el inventario de archivos y componentes. Sin código final.
- **`tasks.md`**: la hoja de ruta de ejecución atómica (consulta [Formato de las tareas](#5-formato-de-las-tareas)).
- **`deltas/<module>.md`**: los cambios propuestos a cada living spec de módulo afectada (consulta [Deltas de módulo](#6-deltas-de-módulo)).

Cuáles de estos archivos son obligatorios depende del tamaño de la especificación.

## 3. Tamaños

Toda especificación tiene un tamaño, elegido al crearla (`specforce spec init <slug> --size <size>`) y modificable después con `specforce spec resize`. El tamaño por defecto es `medium`.

| Tamaño | Artefactos obligatorios | Delta de módulo |
|---|---|---|
| `small` | `tasks.md` | Exento |
| `medium` (por defecto) | `requirements.md`, `tasks.md` | Obligatorio (o `delta.impact: none`) |
| `large` | `requirements.md`, `design.md`, `tasks.md` | Obligatorio (o `delta.impact: none`) |
| `complex` | `requirements.md`, `design.md`, `tasks.md` | Obligatorio (o `delta.impact: none`) |

- La implementación queda bloqueada hasta que existan todos los artefactos requeridos por el tamaño.
- Promover una especificación (por ejemplo de `small` a `large`) marca los nuevos artefactos como faltantes; degradarla relaja la validación y conserva los archivos existentes como referencia.
- En las especificaciones `small`, las tareas pueden omitir el vínculo `**Context:** [US-X]`; a partir de `medium`, toda tarea debe vincularse a un requisito.
- Las hojas de ruta `large` y `complex` deben incluir una sección de Black-box Scenarios (consulta más abajo).

## 4. Especificaciones de bug

Una especificación creada con `--type bug` usa las variantes de bug de los artefactos de planificación, a partir de las plantillas `bug-requirements` y `bug-design`:

- **Requisitos (`bug-requirements`)**: descripción del problema, evidencia (logs, errores), pasos de reproducción concretos, un análisis de causa raíz basado en evidencia y criterios de aceptación de regresión (`[FIX-x]`). Cuando la causa raíz no está establecida, el agente se detiene para preguntar o recomienda `spf.discovery` en lugar de adivinar.
- **Diseño (`bug-design`)**: las rutas de código a cambiar, la estrategia de regresión y los efectos secundarios.
- **Hoja de ruta**: la primera tarea es un test de regresión que falla antes de la corrección (una tarea `[RED]`). Ninguna tarea de corrección puede ir antes que ella.

Los archivos mantienen los nombres habituales (`requirements.md`, `design.md`, `tasks.md`).

## 5. Formato de las tareas

Cada tarea de `tasks.md` es un elemento de checklist como `- [ ] T1.1: [CODE] Title`, seguido de:

- **`**Target:**`**: la ruta exacta del archivo o el comando de la CLI que la tarea cambia.
- **`**Context:**`**: el requisito que implementa (`[US-X]`), obligatorio a partir de `medium`.
- **`**Constraints:**`** *(opcional)*: solo lo que el implementador no puede deducir del diseño y del código (una trampa, una desviación de la convención, una decisión de diseño referenciada). Sin firmas, constantes ni fragmentos de código.
- **`**Acceptance Check:**`**: un único comando no interactivo o escenario de test con un resultado esperado observable. Los criterios que necesitan juicio humano se reescriben o se marcan como no verificables.

Reglas:
- Toda tarea cambia al menos un archivo del producto o un artefacto. Las tareas solo de verificación (ejecutar la suite completa, una ejecución de humo, una comprobación de compilación) están prohibidas: eso corresponde al paso de QA independiente y a los hooks `on_all_tasks_finished`.
- Las comprobaciones de caja negra van en una sección no ejecutable `## 3. Black-box Scenarios` después de las fases de tareas, nunca como tareas. Cada escenario (`S1`, `S2`...) se vincula a un requisito e indica un comando y su resultado esperado. La sección es obligatoria para `large` y `complex`, opcional para las especificaciones `medium` que producen un artefacto ejecutable, y se omite en `small`.
- `specforce spec status` muestra avisos no bloqueantes para las tareas que contienen bloques de código o citan más de tres rutas de archivo.

## 6. Deltas de módulo

A partir de `medium`, una especificación debe declarar cómo cambia las living specs de módulo de `.specforce/docs/modules/` antes de la implementación:

- **Documentos de delta**: un `deltas/<module>.md` por cada módulo afectado, con secciones `## ADDED`, `## MODIFIED` y `## REMOVED` que llevan bloques completos. Las nuevas reglas y casos de uso usan claves provisionales `NEW:<key>`; los identificadores finales se asignan en el momento del archivado.
- **Sin impacto**: cuando ningún módulo cambia, la especificación declara en su lugar `delta.impact: none` en `spec.yaml`.
- **Sello**: `specforce spec delta seal <slug>` registra huellas (fingerprints) de los bloques que el delta modifica o elimina. Si otra especificación cambia el mismo bloque después, la comprobación del delta informa un conflicto y el delta debe rehacerse contra el módulo actual y volver a sellarse.
- **Estados** informados por `specforce spec delta check <slug>`:
  - `none`: se declaró `delta.impact: none`.
  - `exempt`: una especificación `small` sin delta.
  - `missing`: se requiere un delta pero no se ha escrito.
  - `pending`: el delta todavía no se refleja en los módulos.
  - `applied`: los módulos coinciden con el delta.
  - `invalid`: el delta está mal formado o entra en conflicto con los módulos actuales.
- **Barrera de archivado**: `spf.archive` aplica el delta a los módulos; una especificación solo se archiva cuando su estado es `applied`, `none` o `exempt`.

## 7. Informe de QA

Después de la implementación, `spf.implement` ejecuta un paso de QA independiente: un tester que no implementó ninguna tarea comprueba cada criterio de aceptación (y cada elemento del delta y cada escenario de caja negra) y da a cada uno un veredicto:

- **`PASSED`**: verificado con evidencia.
- **`FAILED`**: el criterio no se cumple. Los criterios fallidos desencadenan como máximo dos ciclos de corrección antes de que el agente te lo escale.
- **`NOT-VERIFIABLE`**: un tester no puede ejecutar el criterio; se indica el motivo.

El informe se guarda en `spec.yaml` (un informe por especificación, reemplazado en cada nueva ejecución). Su estado global es `FAILED` cuando algún criterio falló y `PASSED` en caso contrario. `specforce implementation status` lo muestra.

Barrera de archivado: un informe `FAILED` bloquea el archivado (incluso con `--force`); la ausencia de informe solo produce un aviso.

## 8. Tareas test-first (RED)

Una tarea cuyo título lleva el token exacto `[RED]` (distingue mayúsculas y minúsculas) es una tarea test-first: escribe un test que se espera que falle hasta que una tarea `[GREEN]` posterior implemente el comportamiento.

- Cuando todas las tareas finalizadas en un lote son `[RED]`, Specforce omite solo los hooks por tarea `on_task_finished` e imprime un aviso que los nombra, ya que en ese punto se espera que la suite de tests falle.
- Los hooks `on_phase_finished` y `on_all_tasks_finished` se siguen ejecutando, incluso cuando una tarea `[RED]` cierra la fase.
- Un lote que contenga alguna tarea no marcada como `[RED]` ejecuta todos los hooks como de costumbre.

Para que las hojas de ruta sean test-first, añade una instrucción a `tasks` en `config.yaml` (consulta [Configuración](configuration.md)).

## 9. La carpeta `.specforce/`

La carpeta raíz `.specforce/` contiene el estado de la orquestación. Salvo el archivo de configuración, no está pensada para editarse a mano:

- **`config.yaml`**: instrucciones del proyecto, variables de contexto y hooks de validación (consulta [Configuración](configuration.md)).
- **`docs/`**: la Constitución y las living specs de módulo.
- **`specs/`**: especificaciones activas.
- **`archive/`**: especificaciones cerradas por `spf.archive`, el historial de las decisiones técnicas del proyecto.
- **`scratch/`**: evidencia desechable creada por `spf.discovery` (tests, scripts, logs, capturas de pantalla). Se excluye a sí misma del control de versiones y se considera desactualizada en sesiones posteriores.

---

Mantener la calidad de estos artefactos es la mejor forma de asegurar que la IA produzca un código excelente.
