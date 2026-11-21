# Contribuir a Specforce Kit

Antes que nada, ¡gracias por considerar contribuir a Specforce Kit! Son personas como tú las que hacen de Specforce una gran herramienta para toda la comunidad de ingeniería con IA.

## 🚀 Nuestra filosofía: primero la especificación
Usamos nuestra propia herramienta. **Todas las nuevas funcionalidades, correcciones de bugs o cambios arquitectónicos de este repositorio DEBEN seguir el ciclo de vida del Desarrollo Orientado por Especificaciones (SDD).** A continuación, las skills se mencionan por su nombre; consulta [Invocar las skills](docs/es/supported-tools.md#invocar-las-skills) para ver la sintaxis de tu herramienta de IA.

1.  **Especificación:** todo cambio empieza con una especificación creada a través de tu agente de IA (`spf.spec`).
2.  **Implementación:** una vez aprobada, implementamos usando la hoja de ruta de Specforce (`spf.implement`).
3.  **Verificación y archivado:** los tests automatizados y la verificación manual en la TUI son obligatorios antes de cerrar el ciclo (`spf.archive`).

---

## 🐛 Reportar problemas

### Reportes de bugs
- Usa el gestor de Issues de GitHub.
- Describe el comportamiento esperado frente al comportamiento real.
- Proporciona un script o unos pasos mínimos de reproducción.
- Incluye tu sistema operativo y la versión de Specforce (`specforce --version`).

### Solicitudes de funcionalidades
- ¡Nos encantan las ideas nuevas! Abre un issue para discutir el "Por qué" y el "Valor" antes de entrar en detalles técnicos.

---

## 🛠 Proceso de Pull Request

1.  **Fork y rama:** crea una rama de funcionalidad a partir de `main`.
2.  **Crea una especificación:** usa tu agente de IA preferido (Claude Code, Cursor, etc.) y pídele, por ejemplo, "Usa spf.spec para planificar: {descripción}". Genera los artefactos requeridos por el tamaño de la especificación (`requirements.md`, `design.md`, `tasks.md` y los deltas de módulo). Usa primero `spf.discovery` cuando el problema o la causa raíz todavía no estén claros.
3.  **Implementa:** pide "Usa spf.implement para implementar {spec}" para que el agente ejecute las tareas y el paso de QA independiente.
4.  **Estándares de código:**
    - Lenguaje: Go 1.26+
    - Formato: `gofmt` es obligatorio.
    - Linting: debe pasar `golangci-lint`.
    - Estándares: respeta la Constitución de nuestro proyecto (`.specforce/docs/`).
5.  **Tests:** toda lógica nueva debe tener tests unitarios. Las correcciones de bugs requieren tests de regresión.
6.  **Archivado y commit:** pide "Usa spf.archive para archivar {spec}" para actualizar las living specs de módulo y la Constitución. Usa [Conventional Commits](https://www.conventionalcommits.org/) (p. ej. `feat:`, `fix:`, `docs:`).
7.  **Envía:** abre un PR contra la rama `main`. Asegúrate de que el pipeline de CI pase.

---

## 🤖 Desarrollo de Kits (soporte para nuevos agentes)

Specforce Kit está construido para ser agnóstico respecto al modelo. Damos la bienvenida a las contribuciones que añaden soporte para nuevos agentes de IA o mejoran los existentes.

### ¿Qué es un Kit?
El Kit es el conjunto de skills que enseñan a un agente de IA a trabajar con el framework Specforce, más el mapeo que las instala en el directorio oculto de cada herramienta para su descubrimiento nativo (p. ej. `.claude/skills/`, `.cursor/skills/`, `.kilocode/skills/`). Está incrustado en el binario:

- **Plantillas de skills:** `src/internal/agent/kit/skills/<id>/SKILL.yaml` (por ejemplo `spf-spec/SKILL.yaml`). Cada plantilla se traduce al `SKILL.md` nativo de cada herramienta soportada al ejecutar `specforce init`.
- **Manifiesto de herramientas:** `src/internal/agent/kit/kit.yaml` declara cada herramienta soportada: su `name`, su `description`, el directorio `target` y los `mappings` de las categorías que recibe (p. ej. `skills`).

### Añadir una nueva herramienta
1.  **Registra el directorio:** añade el directorio oculto de la herramienta (p. ej. `.newagent/`) al slice `ToolPrefixes` de `src/internal/core/constants.go`.
2.  **Define el mapeo:** añade una entrada para la herramienta en `src/internal/agent/kit/kit.yaml` con su `name`, `description`, `target` y `mappings`. Usa la sobrescritura de `target` a nivel de mapeo y la expansión de variables de entorno (p. ej. `${VAR:-DEFAULT}`) cuando la herramienta necesite archivos fuera de su carpeta por defecto.
3.  **Escrituras globales (opcional):** si la herramienta requiere archivos fuera de la raíz del proyecto, declara `security: { global_write: true }` en su entrada de `kit.yaml` (consulta la entrada `codex`) y documenta esa ubicación.
4.  **Asegura la segmentación del contexto:** comprueba que la nueva integración respete nuestro modelo de **Constitución segmentada**, cargando solo las reglas de arquitectura, UI/UX o seguridad necesarias cuando hagan falta, en lugar de volcarlo todo en la ventana de contexto.

Para cambiar lo que dicen las skills, edita las plantillas `SKILL.yaml`; nunca edites los archivos generados dentro de las carpetas de las herramientas. Durante el desarrollo, ejecuta la CLI con `--dev` desde la raíz del repositorio para cargar el kit desde `src/internal/agent/kit/` en lugar de la copia incrustada en el binario (consulta la [Referencia de la CLI](docs/es/cli.md)).

---

## 🏗 Entorno de desarrollo

### Requisitos previos
- Go 1.26+
- Make

### Configuración
```bash
git clone https://github.com/jeancodogno/specforce-kit.git
cd specforce-kit
make build
make install
```

### Ejecutar los tests
```bash
make test
```

---

## 📜 Código de conducta
Al participar en este proyecto, aceptas cumplir el Código de Conducta del proyecto y mantener un entorno profesional y respetuoso para todas las personas que contribuyen.

---
<p align="center">Construyamos juntos el futuro de la ingeniería orquestada por IA.</p>
