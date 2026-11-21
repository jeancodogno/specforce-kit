[ [English](../en/git-worktrees.md) | [Português](../pt/git-worktrees.md) | [Español](git-worktrees.md) ]

# Soporte para Git Worktrees

## Introducción

El Desarrollo Orientado por Especificaciones suele implicar trabajar en varias funcionalidades o correcciones de bugs en paralelo, en distintas ramas de git. Specforce simplifica este flujo de trabajo ofreciendo soporte nativo para **Git Worktrees**. Esta funcionalidad permite que la consola TUI de Specforce descubra y muestre las especificaciones y el progreso de los worktrees vinculados, ofreciendo una vista unificada del estado de todo el proyecto sin necesidad de cambiar de rama constantemente.

## Motor de descubrimiento

Specforce incluye un Motor de Descubrimiento Unificado que detecta automáticamente los git worktrees vinculados al repositorio actual. Este motor examina el directorio `.git/worktrees` para identificar otros worktrees activos y sus rutas respectivas.

La principal ventaja de esta integración es poder ver simultáneamente, dentro de la consola TUI de Specforce, especificaciones de varias ramas (p. ej. `main`, `feature/X`, `bugfix/Y`). Esto ofrece a los desarrolladores una vista global de todo el trabajo en curso del proyecto, facilitando una mejor coordinación y evitando esfuerzos duplicados.

## Restricciones

Para garantizar la integridad de los datos y evitar modificaciones accidentales entre ramas, Specforce impone las siguientes restricciones:

- **Acceso de solo lectura:** todas las especificaciones descubiertas en git worktrees externos (los que no son tu directorio de trabajo actual) son estrictamente de **solo lectura**.
- **Acciones bloqueadas:** para estas especificaciones externas se bloquean todas las acciones que modifican el estado. Esto incluye reclamar o finalizar tareas, cambiar el estado de las especificaciones y cualquier edición automatizada de archivos.
- **Aislamiento de la Constitución:** Specforce ignora el directorio `.specforce/docs/` (la Constitución del Proyecto) de los worktrees externos. Así se asegura que tu implementación actual siempre respete los principios de arquitectura e ingeniería definidos en tu rama local.

## Alcance del escaneo

El Motor de Descubrimiento Unificado realiza un escaneo profundo de los worktrees vinculados para reunir todos los artefactos de SDD (Spec-Driven Development) relevantes. En el escaneo entre worktrees se incluyen los siguientes elementos:

- **Especificaciones activas:** todos los archivos bajo `.specforce/specs/**` de todos los worktrees descubiertos.
- **Implementaciones:** cualquier archivo de implementación u hoja de ruta referenciado actualmente por especificaciones activas.
- **Especificaciones archivadas:** especificaciones completadas y archivadas previamente, ubicadas en `.specforce/archive/**`.

Al incluir tanto el trabajo activo como el archivado, Specforce ofrece un contexto histórico y actual de la evolución del proyecto en todas las ramas.

## Especificaciones en paralelo y deltas de módulo

Cuando especificaciones de distintos worktrees cambian la misma living spec de módulo en `.specforce/docs/modules/`, el sellado del delta evita que se sobrescriban entre sí sin avisar. Durante la planificación, `specforce spec delta seal <slug>` registra una huella (fingerprint) de cada bloque de módulo que el delta modifica o elimina. Si otra especificación cambia primero uno de esos bloques (por ejemplo, se archiva y se fusiona en tu rama), `specforce spec delta check <slug>` informa un conflicto para ese elemento: el delta debe rehacerse contra el módulo actual y volver a sellarse antes de que la especificación pueda archivarse. Consulta [Deltas de módulo](artifacts.md#6-deltas-de-módulo).
