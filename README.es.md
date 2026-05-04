
<p align="center">
  <a href="https://github.com/jeancodogno/specforce-kit/">
    <picture>
      <img src="assets/logo.png" alt="Logotipo de OpenSpec" height="128">
    </picture>
  </a>
</p>
<p align="center">Ecosistema para el desarrollo asistido por IA.</p>

# Specforce Kit

🌎 **Idiomas / Languages:** [English](README.md) | [Português](README.pt.md) | [Español](README.es.md)

> **La Capa de Orquestación para la Ingeniería Nativa de IA.**

Specforce Kit es un marco de trabajo de nivel profesional diseñado para coordinar agentes de IA (Claude Code, Cursor, KiloCode) a través del **Desarrollo Orientado por Especificaciones (SDD - Spec-Driven Development)**. Proporciona los **Kits**, la **Constitución** y las **Herramientas de Orquestación** estandarizadas necesarias para transformar la IA de un asistente "basado en vibras" en una fuerza de ingeniería determinista.

---

## 🧠 Filosofía

En la era de la IA, la implementación es un commodity, pero el diseño es una responsabilidad. **La implementación sin una especificación está estrictamente prohibida en Specforce.**

El Desarrollo Orientado por Especificaciones (SDD) desacopla el diseño de la implementación. Una especificación actúa como un contrato estricto entre los diseñadores humanos y los ejecutores de máquinas. Al hacerlo, evitamos el "vibe coding" (la acumulación de decisiones ad-hoc y no verificadas), asegurando que su sistema permanezca coherente, seguro y mantenible, independientemente de qué agente de IA escriba el código.

---

## 📊 ¿Por qué usar Specforce?

La codificación con IA sin guías estrictas conduce inevitablemente al "vibe coding": prompts vagos, resultados impredecibles y una montaña de deuda técnica que se acumula rápidamente. Aunque otras herramientas intentan resolver esto, a menudo te obligan a entrar en nuevos ecosistemas de IDE, requieren configuraciones pesadas o saturan la ventana de contexto de tu IA.

Specforce está diseñado para ser la **capa de orquestación invisible**. Aporta previsibilidad, control arquitectónico y verificación estricta a tus flujos de trabajo de IA, actuando como una máquina de estados activa en lugar de una simple colección de plantillas estáticas.

### Fortalezas Clave:

- **Constitución Segmentada**: En lugar de un único prompt masivo, las reglas se dividen en dominios especializados (Arquitectura, UI/UX, Seguridad). El agente solo carga el contexto exacto que necesita, manteniendo bajos los costos de tokens y reduciendo las alucinaciones.
- **Hooks de Verificación Dinámica**: Ejecuta pruebas, linters o scripts personalizados automáticamente antes de que un agente pueda marcar una tarea como finalizada.
- **Soporte de Git Worktree**: Descubre y visualiza de forma nativa las especificaciones en múltiples ramas de git simultáneamente sin cambiar de contexto.
- **Instrucciones Sensibles al Contexto**: Inyecta reglas y restricciones específicas de forma dinámica en función del artefacto o la fase en la que el agente esté trabajando.
- **Agnóstico a Herramientas e IDE**: Lleva el Desarrollo Orientado por Especificaciones directamente a tu terminal. Funciona a la perfección con tu IDE favorito y agentes de IA basados en CLI (Gemini, Claude, Cursor, KiloCode, etc.).
- **Ultrarrápido y Sin Dependencias**: Distribuido como un único binario, no requiere pesadas configuraciones de Python ni ecosistemas forzados para arrancar tu proyecto.

---

## 🤖 Uso del Agente (Comandos Slash)

Specforce interactúa a la perfección con tu agente de IA preferido a través de comandos slash sencillos. Este es el flujo de trabajo estándar:

### 1. Descubrimiento (`/spf:discovery`)
**Comando:** `/spf:discovery {idea o informe de error}`

La fase de exploración temprana. Úsalo para intercambiar ideas sobre nuevas funciones o investigar problemas técnicos sin modificar ningún archivo.
- **Exploración de Solo Lectura:** El agente tiene estrictamente prohibido escribir código o artefactos. Se centra en la investigación, el análisis de la causa raíz y la estrategia técnica.
- **Personas Expertas:** Cambia automáticamente entre "Arquitecto Senior de Producto" (Lluvia de ideas) e "Ingeniero Senior de Sistemas" (Detective) según tu entrada.
- **Alineación con el Proyecto:** Lee la Constitución de tu proyecto (`.specforce/docs/`) para asegurar que todas las ideas se alineen con tus principios y arquitectura.

### 2. La Constitución (`/spf:constitution`)
**Comando:** `/spf:constitution {descripción del proyecto}`

El primer paso en cualquier proyecto. Esto genera la Constitución de tu proyecto, que contiene todas las reglas, principios, pautas de UI/UX, arquitectura, seguridad y memoria del agente.
- **Agnóstico a Herramientas y Segmentado:** En lugar de volcar las reglas en archivos específicos de la herramienta (como `.clauderc` o `.gemini/GEMINI.md`) o mantener todo en el contexto, Specforce mantiene sus propios archivos de memoria especializados y segmentados. El agente solo carga lo que necesita. Puedes cambiar de herramienta (por ejemplo, de Gemini CLI a Claude Code) y la memoria del proyecto permanece completamente intacta.
- **Configuración Interactiva:** Describe tu idea de proyecto y stack. El agente te hará preguntas aclaratorias para ayudarte a tomar decisiones fundamentales.
- **Flexible:** Ejecútalo una vez para configurar, o vuelve a ejecutarlo en cualquier momento para actualizar la constitución. Puedes usarlo en ideas nuevas o apuntarlo a un proyecto existente para analizar y documentar su arquitectura actual.

### 2. Creación de una Spec (`/spf:spec`)
**Comando:** `/spf:spec {descripción de qué añadir o modificar}`

¿Listo para construir? El agente aclarará tus requisitos y generará tres documentos principales: `requirements.md`, `design.md` y `tasks.md`.

### 3. Implementación (`/spf:implement`)
**Comando:** `/spf:implement`

Después de validar los documentos de especificación generados, ejecuta este comando. El agente seguirá estrictamente las tareas planificadas en la especificación para implementar la función, garantizando una desviación cero del diseño acordado.

### 4. Archivado (`/spf:archive`)
**Comando:** `/spf:archive`

Una vez que la implementación está completada y verificada, este comando analiza qué debe actualizarse en la Constitución global del proyecto basándose en la nueva implementación, y archiva la especificación completada. Sigue las **Instrucciones de Archivado** estandarizadas que aseguran que las lecciones aprendidas se capturen en el memorial del proyecto y que se realice cualquier limpieza específica del proyecto.

---

## ⚡ Toque Humano

Mientras Specforce potencia a tus agentes de IA, los humanos interactúan principalmente a través de dos comandos de terminal:

### 1. Inicialización del Proyecto
Configura la gobernanza y los kits de agentes para tu proyecto.
```bash
specforce init
```

### 2. La Consola (TUI)
Inicia el centro de mando para supervisar la salud del proyecto, el progreso de las especificaciones y las implementaciones de los agentes en tiempo real.
```bash
specforce console
```

---

## 🛠️ Agentes de IA y Herramientas Soportadas

Specforce Kit está diseñado para ser agnóstico a las herramientas, pero viene con kits preconstruidos y soporte nativo de comandos slash para los asistentes de codificación de IA más populares:

- **Gemini CLI** (Google)
- **Claude Code** (Anthropic)
- **Qwen** (Alibaba)
- **Kimi Code** (Moonshot AI)
- **OpenCode**
- **KiloCode**
- **Codex**
- **Antigravity**

*¿No ves tu herramienta favorita? ¡Aceptamos PRs para añadir nuevos kits!*

---

## 🚀 Primeros Pasos

### Instalación

**Usando NPM (Recomendado)**
```bash
npm i -g @jeancodogno/specforce-kit
```

**Desde el Código Fuente (Requiere Go 1.26+ y Make)**
```bash
git clone https://github.com/jeancodogno/specforce-kit.git
cd specforce-kit
make build
make install
```

---

## 📚 Documentación

Para profundizar en Specforce Kit, consulta nuestra documentación oficial:

- [Primeros Pasos y Flujo de Trabajo](docs/es/getting-started.md): Configuración del entorno, primeros pasos y el ciclo de vida recomendado del Desarrollo Orientado por Especificaciones.
- [Soporte de Git Worktree](docs/es/git-worktrees.md): Guía para el descubrimiento de especificaciones entre ramas, consola unificada y restricciones de solo lectura.
- [Artefactos](docs/es/artifacts.md): Detalles sobre los archivos generados (Constitución, Requisitos, Diseño, Tareas).
- [Configuración](docs/es/configuration.md): Guía para personalizar hooks e instrucciones.
- [Referencia de la CLI](docs/es/cli.md): Uso de la terminal y comandos slash.
- [Herramientas Soportadas](docs/es/supported-tools.md): Lista de agentes de IA y asistentes de codificación soportados de forma nativa.

---

## 🔍 Resolución de Problemas

### "command not found: specforce"
Si te encuentras con este error después de la instalación, suele significar que el directorio bin global de npm no está en el `PATH` de tu sistema.

**Solución para macOS/Linux:**
Añade la siguiente línea a tu perfil de shell (por ejemplo, `~/.zshrc` o `~/.bashrc`):
```bash
export PATH="$(npm config get prefix)/bin:$PATH"
```

**Solución para Windows:**
Asegúrate de que `%AppData%\npm` esté en tus variables de entorno.

Para una resolución de problemas más detallada, consulta [Primeros Pasos: Resolución de Problemas](docs/es/getting-started.md#troubleshooting).

---

## 🤝 Contribuir
¡Agradecemos las contribuciones! Consulta nuestro [CONTRIBUTING.md](CONTRIBUTING.md) para aprender cómo añadir soporte para nuevos agentes de IA o mejorar los kits existentes.

---

## 📜 Licencia
Specforce Kit se publica bajo la **Licencia MIT**. Consulta [LICENSE](LICENSE) para más detalles.
