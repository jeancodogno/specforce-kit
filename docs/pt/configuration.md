[ [English](../en/configuration.md) | [Português](configuration.md) | [Español](../es/configuration.md) ]

# Configuração do Specforce

O comportamento do Specforce é configurado no arquivo `.specforce/config.yaml` (criado pelo `specforce init`). Com ele, você adiciona regras do projeto, variáveis compartilhadas e barreiras de qualidade que se aplicam automaticamente aos seus agentes de IA.

## 1. Instruções Sensíveis ao Contexto

Em vez de repetir regras ("Lembre-se de usar BDD") em todo prompt, você pode configurá-las para serem injetadas automaticamente no momento exato em que o agente trabalha em um artefato ou fase específicos.

No `config.yaml`:

```yaml
instructions:
  requirements:
    - "Sempre use a sintaxe BDD GIVEN/WHEN/THEN nos critérios de aceitação."
    - "Garanta que a acessibilidade seja sempre mencionada quando houver UI envolvida."
  bug-requirements:
    - "Sempre inclua a mensagem de erro exata na seção de evidências."
  design:
    - "Sempre avalie os impactos de segurança de novas rotas de API."
  tasks:
    - "Estruture cada fase como tarefas de teste [RED] seguidas de tarefas de implementação [GREEN]."
  implementation:
    - "Não contorne as tipagens do TypeScript. Use Type Guards."
  archive:
    - "Sempre consolide as especificações vivas dos módulos com as lições aprendidas."
    - "Resuma qualquer dívida técnica introduzida durante a implementação."
```

### Chaves de instrução aceitas

| Chave | Aplicada quando |
|---|---|
| `requirements` | o agente escreve o `requirements.md` (especificações de funcionalidade e de bug) |
| `bug-requirements` | o agente escreve o `requirements.md` de uma especificação de bug |
| `design` | o agente escreve o `design.md` (especificações de funcionalidade e de bug) |
| `bug-design` | o agente escreve o `design.md` de uma especificação de bug |
| `tasks` | o agente escreve o `tasks.md` |
| `implementation` | o agente carrega o roadmap para implementá-lo |
| `archive` | o agente arquiva uma especificação |

**Chaves genéricas e específicas são combinadas.** Quando o agente trabalha em um artefato variante como `bug-requirements`, ele recebe primeiro as instruções da chave genérica (`requirements`) e depois as da chave específica (`bug-requirements`); uma regra listada nas duas aparece apenas uma vez. O mesmo vale para `design` e `bug-design`.

## 2. Injeção de Variáveis nas Instruções

Para tornar as instruções ainda mais poderosas, você pode definir variáveis no nível do projeto no bloco `context`. Essas variáveis podem então ser injetadas em qualquer instrução ou template de artefato usando a sintaxe `{{variable_name}}`.

No `config.yaml`:

```yaml
# Variáveis de contexto globais
context:
  project_name: "Specforce Kit"
  tech_stack: "Go, React, PostgreSQL"
  primary_branch: "main"

instructions:
  tasks:
    - "Lembre-se: o código será mesclado na branch {{primary_branch}}."
    - "Toda nova funcionalidade do {{project_name}} deve respeitar a stack {{tech_stack}}."
```

### Por que usar variáveis?
- **Fonte Única da Verdade:** Atualize o nome do projeto ou a stack uma única vez, e todas as instruções voltadas ao agente ficam sincronizadas.
- **Templates Reutilizáveis:** Use o mesmo conjunto de instruções em projetos diferentes alterando apenas o bloco `context`.
- **Consistência:** Garante que os agentes sempre tenham dados corretos do projeto, sem alucinações.

Nunca coloque segredos (tokens, senhas, chaves) em `context` ou `instructions`: esses valores são mostrados ao agente. Mantenha os segredos em arquivos de ambiente.

## 3. Hooks de Validação (Quality Gates)

Hooks são comandos externos que o Specforce executa quando o agente marca tarefas como finalizadas. Seu principal uso é atuar como **Quality Gates** (barreiras de qualidade) que impedem o agente de finalizar tarefas enquanto o código estiver inadequado.

No `config.yaml`:

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

### Hooks suportados:
- **`on_task_finished`**: Disparado sempre que tarefas são marcadas como finalizadas.
- **`on_phase_finished`**: Disparado quando uma tarefa finalizada é a última da sua Fase (título H3 no `tasks.md`).
- **`on_all_tasks_finished`**: Disparado quando uma tarefa finalizada é a última do roadmap.

### Como funciona?
1. O agente marca uma ou mais tarefas como finalizadas (`specforce implementation update <slug> --task <ids> --status finished`).
2. Antes de alterar qualquer estado, o Specforce executa os hooks aplicáveis uma única vez para o lote inteiro, sem duplicatas (no exemplo acima, `npm run lint` e `npm run test`).
3. **Se todos os hooks terminarem com código 0:** as tarefas são marcadas como finalizadas.
4. **Se qualquer hook terminar com código diferente de zero:** a transição é bloqueada e a saída do hook é devolvida ao agente para que ele corrija o problema antes de tentar de novo.

Tarefas test-first: quando todas as tarefas do lote têm `[RED]` no título, apenas os hooks `on_task_finished` são pulados (com um aviso citando-os); os hooks de fase e de todas as tarefas continuam rodando. Veja [Artefatos](artifacts.md#8-tarefas-test-first-red).

### Dicas essenciais para hooks:
- **Um comando por entrada:** cada entrada é dividida nos espaços e executada diretamente, não por meio de um shell, então pipes (`|`), redirecionamentos e `&&` não são interpretados. Coloque pipelines em um script ou em um alvo do `make` e chame-o (por exemplo `make test-quiet`).
- **Os hooks de um lote rodam ao mesmo tempo:** não faça um hook depender do resultado de outro.
- **Seja rápido e silencioso:** Use comandos rápidos e silenciosos (`--silent`, `--quiet`).
- **Filtre a saída:** Não despeje milhares de linhas de log no contexto do agente. Dentro do seu script, envie apenas os erros relevantes, por exemplo `golangci-lint run --out-format=line-number` ou `npm run test --silent | grep -A 10 "FAIL"`.

---

Combinando as *Instruções Sensíveis ao Contexto* e os *Hooks de Validação*, você cria uma "proteção" inquebrável para a sua equipe e para os seus agentes de IA.
