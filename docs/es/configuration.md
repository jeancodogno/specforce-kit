[ [English](../en/configuration.md) | [Português](../pt/configuration.md) | [Español](configuration.md) ]

# Configuración de Specforce

El comportamiento de Specforce se configura en el archivo `.specforce/config.yaml` (creado por `specforce init`). Con él añades reglas del proyecto, variables compartidas y barreras de calidad que se aplican automáticamente a tus agentes de IA.

## 1. Instrucciones según el contexto

En lugar de repetir reglas ("Recuerda usar BDD") en cada prompt, puedes configurarlas para que se inyecten automáticamente en el momento exacto en que el agente trabaja en un artefacto o en una fase concreta.

En `config.yaml`:

```yaml
instructions:
  requirements:
    - "Always name the end-user persona in every user story."
    - "Ensure accessibility is always mentioned if UI is involved."
  bug-requirements:
    - "Always include the exact error message in the evidence section."
  design:
    - "Always evaluate security impacts for new API routes."
  tasks:
    - "Structure every phase as [RED] test tasks followed by [GREEN] implementation tasks."
  implementation:
    - "Do not bypass TypeScript typings. Use Type Guards."
  archive:
    - "Always consolidate module living specs with lessons learned."
    - "Summarize any technical debt introduced during implementation."
```

### Claves de instrucción aceptadas

| Clave | Se aplica cuando |
|---|---|
| `requirements` | el agente escribe `requirements.md` (especificaciones de funcionalidad y de bug) |
| `bug-requirements` | el agente escribe el `requirements.md` de una especificación de bug |
| `design` | el agente escribe `design.md` (especificaciones de funcionalidad y de bug) |
| `bug-design` | el agente escribe el `design.md` de una especificación de bug |
| `tasks` | el agente escribe `tasks.md` |
| `implementation` | el agente carga la hoja de ruta para implementarla |
| `archive` | el agente archiva una especificación |

**Las claves genéricas y las específicas se combinan.** Cuando el agente trabaja en un artefacto variante como `bug-requirements`, recibe primero las instrucciones de la clave genérica (`requirements`) y después las de la clave específica (`bug-requirements`); una regla listada en ambas aparece una sola vez. Lo mismo se aplica a `design` y `bug-design`.

## 2. Inyección de variables en las instrucciones

Para que las instrucciones sean todavía más potentes, puedes definir variables a nivel de proyecto en el bloque `context`. Después, estas variables pueden inyectarse en cualquier instrucción o plantilla de artefacto usando la sintaxis `{{variable_name}}`.

En `config.yaml`:

```yaml
# Global context variables
context:
  project_name: "Specforce Kit"
  tech_stack: "Go, React, PostgreSQL"
  primary_branch: "main"

instructions:
  tasks:
    - "Remember: the code will be merged into the {{primary_branch}} branch."
    - "Every new feature in {{project_name}} must respect the {{tech_stack}} stack."
```

### ¿Por qué usar variables?
- **Fuente única de verdad:** actualiza el nombre del proyecto o el stack una sola vez, y todas las instrucciones dirigidas al agente se mantienen sincronizadas.
- **Plantillas reutilizables:** usa el mismo conjunto de instrucciones en distintos proyectos cambiando solo el bloque `context`.
- **Consistencia:** garantiza que los agentes siempre tengan datos del proyecto correctos y no alucinados.

Nunca pongas secretos (tokens, contraseñas, claves) en `context` ni en `instructions`: estos valores se muestran al agente. Guarda los secretos en archivos de entorno.

## 3. Hooks de validación (Quality Gates)

Los hooks son comandos externos que Specforce ejecuta cuando el agente marca tareas como finalizadas. Su uso principal es actuar como **Quality Gates** (barreras de calidad) que impiden que el agente finalice tareas mientras el código no sea adecuado.

En `config.yaml`:

```yaml
hooks:
  on_task_finished:
    - "npm run lint"
    - "npm run test"
  on_phase_finished:
    - "npm run integration-tests"
  on_all_tasks_finished:
    - "npm run e2e-tests"
```

### Hooks soportados:
- **`on_task_finished`**: se dispara cada vez que se marcan tareas como finalizadas.
- **`on_phase_finished`**: se dispara cuando una tarea finalizada es la última de su Fase (encabezado H3 en `tasks.md`).
- **`on_all_tasks_finished`**: se dispara cuando una tarea finalizada es la última de la hoja de ruta.

### ¿Cómo funciona?
1. El agente marca una o más tareas como finalizadas (`specforce implementation update <slug> --task <ids> --status finished`).
2. Antes de cambiar cualquier estado, Specforce ejecuta los hooks aplicables una sola vez para todo el lote, sin duplicados (en el ejemplo anterior, `npm run lint` y `npm run test`).
3. **Si todos los hooks terminan con código 0:** las tareas se marcan como finalizadas.
4. **Si algún hook termina con un código distinto de cero:** la transición se bloquea y la salida del hook se devuelve al agente para que pueda corregir el problema antes de volver a intentarlo.

Tareas test-first: cuando todas las tareas del lote llevan `[RED]` en su título, solo se omiten los hooks `on_task_finished` (con un aviso que los nombra); los hooks de fase y de todas las tareas se siguen ejecutando. Consulta [Artefactos](artifacts.md#8-tareas-test-first-red).

### Consejos esenciales para los hooks:
- **Un comando por entrada:** cada entrada se divide por espacios y se ejecuta directamente, no a través de un shell, así que las tuberías (`|`), las redirecciones y `&&` no se interpretan. Pon las tuberías en un script o en un target de `make` y llama a eso en su lugar (por ejemplo `make test-quiet`).
- **Los hooks de un mismo lote se ejecutan a la vez:** no hagas que un hook dependa del resultado de otro.
- **Sé rápido y silencioso:** usa comandos rápidos y silenciosos (`--silent`, `--quiet`).
- **Filtra la salida:** no vuelques miles de líneas de log en el contexto del agente. Dentro de tu script, envía solo los errores relevantes, p. ej. `golangci-lint run --out-format=line-number` o `npm run test --silent | grep -A 10 "FAIL"`.

---

Combinando las *Instrucciones según el contexto* y los *Hooks de validación*, creas una "barandilla" irrompible para tu equipo y tus agentes de IA.
