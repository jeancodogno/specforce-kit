[ [English](../en/workflows.md) | [Português](workflows.md) | [Español](../es/workflows.md) ]

# Fluxos de Trabalho e Dicas

Esta página explica alguns fluxos de trabalho do Specforce e dicas para usá-lo.

As skills são citadas pelo nome (por exemplo [`spf.spec`](supported-tools.md#invocando-skills)). Como chamá-las explicitamente na sua ferramenta é mostrado em [Invocando skills](supported-tools.md#invocando-skills); um pedido em linguagem natural que cite o nome da skill funciona em qualquer ferramenta.

**Como ler os fluxos**

- 💻 **Terminal:** um comando que **você** digita no seu próprio terminal. Apenas alguns comandos aparecem aqui: `specforce init`, `specforce console`, `specforce spec list` e alguns comandos `git`. Todo o resto (status, verificação do delta, resultados de QA) o agente confere e informa a você no chat.
- 💬 **Chat do agente:** um pedido que **você** digita no chat do seu agente de IA. Substitua as partes entre `<sinais de menor e maior>`.
- **Por trás dos panos** (recolhido): o que o agente executa por conta própria. Você nunca precisa rodar esses comandos; eles estão listados para que você acompanhe o que acontece. Todos estão descritos na [Referência da CLI](cli.md).
- `<slug>` é o nome da especificação, por exemplo `csv-export`. `specforce spec list` mostra esse nome.

## Fluxos

### 1. Prototipar na descoberta antes de implementar

**Quando usar:** a ideia ainda está vaga, ou você quer uma prova (um teste rápido, um script, uma captura de tela) antes de se comprometer com um plano.

1. 💬 **Chat do agente:**
   ```text
   Use spf.discovery para explorar: no nosso SaaS de faturamento, permitir que o cliente exporte todas as faturas em CSV pelo painel. Contas grandes têm 200 mil faturas, então prove com um teste descartável que gerar a exportação em streaming mantém o uso de memória estável antes de planejarmos qualquer coisa.
   ```
2. 💬 **Chat do agente:** quando vocês concordarem com a abordagem:
   ```text
   Formalize isto em uma proposta.
   ```
3. 💬 **Chat do agente:**
   ```text
   Use spf.spec para planejar csv-export a partir da proposta.
   ```

<details><summary>Por trás dos panos</summary>

- A descoberta nunca edita, renomeia ou apaga seus arquivos existentes. As evidências (testes, scripts, logs, capturas de tela) vão para `.specforce/scratch/`, que se ignora no controle de versão. Se uma hipótese só puder ser provada dentro da árvore de código, o agente cria um arquivo **novo** cujo nome contém `spf-probe`, registra-o e o remove antes do fim da sessão, conferindo o `git status`.
- A formalização executa `specforce spec init <slug>` e escreve o `proposal.md` (problema, descobertas, decisões tomadas, alternativas rejeitadas, perguntas em aberto, tipo e tamanho sugeridos, caminhos das evidências).
- O `spf.spec` reapresenta as decisões que já estão na proposta para uma única confirmação e pergunta apenas as questões em aberto.

</details>

**Se der errado:**
- A descoberta não empurra você para uma especificação por conta própria; peça a proposta explicitamente quando estiver pronto.
- Um arquivo `spf-probe` esquecido aparece no `git status`: peça ao agente para removê-lo.
- As evidências em scratch de uma sessão anterior estão desatualizadas; peça ao agente para executá-las de novo antes de confiar nelas.

### 2. Corrigir um bug com um teste que falha primeiro

**Quando usar:** algo está quebrado e você quer que a correção seja comprovada por um teste de regressão.

1. 💬 **Chat do agente:**
   ```text
   Use spf.discovery para investigar este bug: no painel de faturas, o filtro "mês passado" deixa de fora as faturas emitidas no último dia do mês para clientes em UTC-3. Nenhum erro aparece, a lista só vem incompleta.
   ```
2. 💬 **Chat do agente:** quando a causa raiz estiver confirmada:
   ```text
   Formalize isto em uma especificação de bug.
   ```
3. 💬 **Chat do agente:** revise e aprove a especificação, depois:
   ```text
   Use spf.implement para implementar invoice-date-filter.
   ```

<details><summary>Por trás dos panos</summary>

- A especificação de bug é criada com `specforce spec init <slug> --type bug`. Seus requisitos contêm os passos de reprodução, as evidências e a causa raiz.
- O roadmap começa com uma tarefa `[RED]` que adiciona um teste de regressão que falha antes da correção; nenhuma tarefa de correção vem antes dela.

</details>

**Se der errado:**
- O `spf.spec` para quando a causa raiz não está estabelecida: volte ao passo 1 em vez de deixá-lo adivinhar.
- A tarefa `[RED]` é finalizada enquanto os testes falham: isso é esperado; apenas os hooks por tarefa são pulados (veja [Falha de hook bloqueia uma tarefa](#falha-de-hook-bloqueia-uma-tarefa)).

### 3. Fazer uma mudança pequena

**Quando usar:** uma mudança contida (uma correção em um só lugar, um ajuste de texto, uma opção de configuração) que não precisa de requisitos nem de design.

1. 💬 **Chat do agente:**
   ```text
   Use spf.spec para planejar: adicionar a coluna "Vencimento" na tabela de faturas do painel.
   ```
2. 💬 **Chat do agente:** revise a lista de tarefas, depois:
   ```text
   Use spf.implement para implementar invoice-due-date-column.
   ```

<details><summary>Por trás dos panos</summary>

- Você não precisa dizer que a mudança é pequena: o agente define o tamanho durante o planejamento e executa `specforce spec init <slug> --size small`: apenas o `tasks.md` é obrigatório e o delta de módulo é `exempt` (veja [Tamanhos](artifacts.md#3-tamanhos)).

</details>

**Se der errado:**
- A mudança acaba sendo maior do que o esperado: siga o [fluxo 4](#4-o-escopo-cresceu-redimensionar-a-especificação).

### 4. O escopo cresceu: redimensionar a especificação

**Quando usar:** durante o planejamento ou a implementação, a mudança acaba precisando de mais (ou de menos) do que o planejado.

1. 💬 **Chat do agente:** basta descrever o novo escopo:
   ```text
   Use spf.spec para atualizar invoice-due-date-column: além da coluna, enviar ao cliente um e-mail de lembrete 3 dias antes do vencimento.
   ```
2. 💬 **Chat do agente:** depois de aprovar os novos artefatos:
   ```text
   Use spf.implement para continuar invoice-due-date-column.
   ```

<details><summary>Por trás dos panos</summary>

- Você não precisa pedir o redimensionamento: o agente percebe que o escopo mudou e executa `specforce spec resize <slug> --size <novo tamanho>` por conta própria. Promover marca os artefatos que passaram a ser obrigatórios como ausentes e bloqueia a implementação até que eles existam; rebaixar relaxa a validação e mantém os arquivos existentes como referência.
- Em seguida o agente escreve os artefatos que o novo tamanho exige e, a partir de `medium`, escreve e sela o delta de módulo, ou declara `delta.impact: none`.

</details>

**Se der errado:**
- A implementação se recusa a continuar: veja [Redimensionamento durante a implementação](#redimensionamento-durante-a-implementação).

### 5. Mudar depois de implementar, antes de arquivar

**Quando usar:** a implementação está pronta (ou quase) e você quer algo diferente antes de fechar a especificação.

1. 💬 **Chat do agente:**
   ```text
   Use spf.spec para atualizar invoice-due-date-column: enviar o lembrete 5 dias antes do vencimento em vez de 3. Mantenha requisitos, design, delta e tarefas consistentes.
   ```
2. 💬 **Chat do agente:** revise e aprove a especificação atualizada, depois:
   ```text
   Use spf.implement para implementar apenas as tarefas novas ou alteradas de invoice-due-date-column.
   ```

<details><summary>Por trás dos panos</summary>

- Uma especificação aprovada é um contrato: o agente atualiza todos os artefatos afetados juntos, nunca apenas o `tasks.md`, e adiciona tarefas novas em vez de reescrever as finalizadas.
- O `spf.implement` executa depois uma nova etapa de QA independente e registra um novo veredito.

</details>

**Se der errado:**
- Não edite o código nem o `tasks.md` à mão para "só consertar": veja [Editar uma especificação aprovada](#editar-uma-especificação-aprovada).

### 6. Mudar depois de arquivar: especificação de continuação

**Quando usar:** a especificação já foi arquivada e você quer mudar ou estender o que ela entregou.

1. 💬 **Chat do agente:**
   ```text
   Use spf.spec para planejar uma continuação de csv-export: permitir que o cliente escolha um período antes de exportar.
   ```
2. 💬 **Chat do agente:** revise e aprove, depois use o nome que o agente deu à nova especificação:
   ```text
   Use spf.implement para implementar csv-export-date-range.
   ```

<details><summary>Por trás dos panos</summary>

- Especificações arquivadas são histórico e nunca são reabertas. A nova especificação tem sua própria pasta.
- Quando a mudança altera regras que já estão em uma especificação viva de módulo, o delta as traz em `## MODIFIED` (ou `## REMOVED`) com o bloco novo completo (veja [Deltas de módulo](artifacts.md#6-deltas-de-módulo)).
- O resumo escrito pelo `spf.archive` lista próximos passos sugeridos e especificações de continuação: um bom ponto de partida para o pedido.

</details>

**Se der errado:**
- A verificação do delta reporta um conflito: veja [Conflito de base do delta](#conflito-de-base-do-delta).

### 7. QA falhou

**Quando usar:** a etapa de QA independente reportou critérios `FAILED` e o agente escalou para você.

1. 💬 **Chat do agente:** leia o relatório do agente: ele já lista cada critério que falhou com a sua evidência. Se o código estiver errado:
   ```text
   Use spf.implement para corrigir os critérios de QA que falharam em invoice-due-date-column e rodar o QA de novo.
   ```
   Se o próprio critério estiver errado, siga o [fluxo 5](#5-mudar-depois-de-implementar-antes-de-arquivar).

<details><summary>Por trás dos panos</summary>

- O agente já executou até dois ciclos de correção (um implementador corrige, um testador novo verifica de novo) antes de escalar, então o que chega até você é só o que ele não conseguiu corrigir sozinho.
- O veredito é armazenado no `spec.yaml` com `specforce implementation qa record <slug>` após cada execução de QA.

</details>

**Se der errado:**
- O arquivamento é recusado enquanto o veredito armazenado for `FAILED`, mesmo com `--force`: veja [Arquivamento recusado: QA FAILED](#arquivamento-recusado-qa-failed).

### 8. Especificações em paralelo em worktrees

**Quando usar:** você quer dois agentes trabalhando em duas especificações ao mesmo tempo, cada um em sua própria branch.

1. 💻 **Terminal:** crie uma worktree para a segunda especificação:
   ```bash
   git worktree add ../myapp-csv-export -b feature/csv-export
   ```
2. 💬 **Chat do agente:** inicie uma segunda sessão do agente dentro de `../myapp-csv-export`, depois:
   ```text
   Use spf.spec para planejar: permitir que o cliente exporte todas as faturas em CSV pelo painel.
   ```
3. 💻 **Terminal:** acompanhe as duas a partir do seu checkout principal:
   ```bash
   specforce console
   ```

<details><summary>Por trás dos panos</summary>

- Especificações de outras worktrees são somente leitura no seu checkout: suas tarefas só podem ser alteradas a partir da própria worktree. A Constituição sempre vem da sua branch atual (veja [Suporte a Git Worktree](git-worktrees.md)).
- Cada especificação sela o seu delta. Quando as duas alteram o mesmo bloco de módulo, a segunda a ser arquivada vê um conflito e precisa refazer o seu delta.

</details>

**Se der errado:**
- O agente na nova worktree não enxerga as skills: se as pastas das suas ferramentas não estão commitadas, rode `specforce init` dentro dessa worktree.
- Duas especificações mexem na mesma regra de módulo: veja [Conflito de base do delta](#conflito-de-base-do-delta).

### 9. Adotar o Specforce em um projeto existente

**Quando usar:** o projeto já tem código e você quer que o Specforce o documente antes da primeira especificação.

1. 💻 **Terminal:** na raiz do projeto:
   ```bash
   specforce init
   ```
2. 💬 **Chat do agente:**
   ```text
   Use spf.constitution para documentar este projeto existente: leia o código e redija a constituição a partir do que já existe.
   ```
3. 💬 **Chat do agente:** revise os documentos em `.specforce/docs/`, corrija o que estiver errado e então planeje a primeira mudança:
   ```text
   Use spf.spec para planejar: adicionar uma busca na tabela de faturas que filtre pelo nome do cliente.
   ```

<details><summary>Por trás dos panos</summary>

- O `specforce init` escreve suas regras dentro de um bloco marcado do `AGENTS.md` da raiz; o seu conteúdo existente fora do bloco é preservado.
- O `spf.constitution` lê o código, pergunta sobre as decisões que não consegue inferir e escreve os seis documentos mais `modules/`.

</details>

**Se der errado:**
- A constituição descreve o código como ele é, incluindo hábitos que você quer abandonar: diga isso no chat; ela registra a regra desejada, não o acidente.

### 10. Retomar uma implementação interrompida

**Quando usar:** a sessão travou, você fechou o terminal ou o agente parou no meio de um roadmap.

1. 💬 **Chat do agente:**
   ```text
   Use spf.implement para verificar o status da implementação de invoice-due-date-column e retomá-la.
   ```

<details><summary>Por trás dos panos</summary>

- O agente lê os estados das tarefas e o repositório. Tarefas deixadas em `in-progress` foram interrompidas: ele inspeciona o código afetado e as executa de novo. Tarefas finalizadas não são refeitas, e as mudanças não commitadas são mantidas.

</details>

**Se der errado:**
- Não descarte as mudanças não commitadas antes de retomar: o agente as usa para continuar com segurança.

### 11. Trocar de ferramenta de IA no meio do projeto

**Quando usar:** você passa para outra ferramenta de IA (ou adiciona uma) no meio de um projeto ou de uma especificação.

1. 💻 **Terminal:** liste todas as ferramentas que quer manter, usando as chaves de [Ferramentas Suportadas](supported-tools.md):
   ```bash
   specforce init claude cursor
   ```
   ⚠️ Toda ferramenta instalada que você deixar de fora é desinstalada e **a pasta inteira dela é apagada** (por exemplo `.claude/`, com suas configurações e regras), sem pedir confirmação. Faça commit ou backup dessa pasta antes se ela tiver algo seu.
2. 💬 **Chat do agente:** na nova ferramenta:
   ```text
   Use spf.implement para verificar o status da implementação de invoice-due-date-column e retomá-la.
   ```

<details><summary>Por trás dos panos</summary>

- A Constituição, as especificações e o progresso delas ficam em `.specforce/` e não dependem de nenhuma ferramenta, então nada se perde.
- O `specforce init` instala as skills das ferramentas listadas e apaga a pasta de cada ferramenta instalada que você deixou de fora. Sem argumentos, o seletor interativo mostra um aviso antes de apagar qualquer coisa.

</details>

**Se der errado:**
- A pasta de uma ferramenta sumiu: você deixou a chave dela de fora. Restaure a pasta pelo controle de versão se ela tinha arquivos seus e rode `specforce init` de novo incluindo a chave, ou sem argumentos para usar o seletor.

### 12. Impor TDD pela configuração

**Quando usar:** você quer que todo roadmap seja test-first e que toda tarefa passe nas suas verificações antes de ser marcada como finalizada.

1. 💬 **Chat do agente:** (ou edite o arquivo você mesmo)
   ```text
   Em .specforce/config.yaml, adicione uma instrução de tasks exigindo tarefas de teste [RED] seguidas de tarefas [GREEN] na mesma fase, e os hooks on_task_finished "make lint" e "make test".
   ```
   O resultado fica assim:
   ```yaml
   instructions:
     tasks:
       - "Cada fase deve conter tarefas [RED] (testes para o caminho feliz e casos de borda) seguidas de tarefas [GREEN] (a implementação)."
   hooks:
     on_task_finished:
       - "make lint"
       - "make test"
   ```
2. 💬 **Chat do agente:**
   ```text
   Use spf.spec para planejar: permitir que o cliente baixe cada fatura em PDF pelo painel.
   ```

<details><summary>Por trás dos panos</summary>

- As tarefas recebem o token exato `[RED]` no título. Quando todas as tarefas finalizadas em um lote são `[RED]`, apenas os hooks `on_task_finished` são pulados (com um aviso); os hooks de fase e de todas as tarefas continuam rodando.

</details>

**Se der errado:**
- Uma entrada de hook com `|` ou `&&` nunca funciona: veja [Falha de hook bloqueia uma tarefa](#falha-de-hook-bloqueia-uma-tarefa).

## Situações comuns

### Arquivamento recusado: tarefas pendentes

**Sintoma:** o arquivamento para e lista tarefas que não foram finalizadas.

**Por quê:** uma especificação só é arquivada quando todas as tarefas estão finalizadas. A skill de arquivamento recusa mesmo que você peça para arquivar assim mesmo; `--force` pula apenas esta verificação e nunca as verificações de delta ou de QA.

**O que fazer:** o agente já lista as tarefas não finalizadas. Peça `Use spf.implement para finalizar <slug>`. Se uma tarefa não for mais necessária, peça ao `spf.spec` para removê-la da especificação.

### Arquivamento recusado: delta ausente ou pendente

**Sintoma:** o arquivamento é recusado com um estado de delta `missing`, `pending` ou `invalid`.

**Por quê:** a partir de `medium`, as especificações vivas dos módulos precisam refletir a especificação antes do arquivamento: o delta deve estar `applied` (ou declarado `none`). `missing` significa que nenhum delta foi escrito; `pending` significa que ele ainda não foi aplicado aos módulos.

**O que fazer:** o agente informa o estado do delta. Para `missing`, peça `Use spf.spec para escrever o delta de módulo de <slug>, ou declarar que não há impacto em módulos`. Para `pending`, peça `Use spf.archive para aplicar o delta e arquivar <slug>`. Para `invalid`, peça ao `spf.spec` para corrigir os problemas reportados.

### Conflito de base do delta

**Sintoma:** o agente reporta um conflito de delta em um bloco de módulo, e o arquivamento deixa o módulo intocado.

**Por quê:** depois que o seu delta foi selado, outra especificação (muitas vezes de uma branch paralela) alterou o mesmo bloco de módulo. Aplicar o seu sobrescreveria essa mudança silenciosamente.

**O que fazer:** peça `Use spf.spec para refazer o delta de <slug> com base no módulo atual e selá-lo de novo`. O agente verifica o delta de novo antes de arquivar.

### Arquivamento recusado: QA FAILED

**Sintoma:** o arquivamento para e lista critérios de QA que falharam.

**Por quê:** o veredito de QA armazenado é `FAILED`. O arquivamento nunca prossegue sobre um veredito de falha, nem mesmo com `--force`.

**O que fazer:** leia os critérios que falharam listados pelo agente e então siga o [fluxo 7](#7-qa-falhou).

### Aviso no arquivamento: sem relatório de QA

**Sintoma:** o arquivamento continua, mas avisa que não existe relatório de QA.

**Por quê:** nenhum veredito foi registrado para a especificação, normalmente porque a implementação foi feita ou finalizada fora do `spf.implement`, ou a sessão terminou antes do QA.

**O que fazer:** se você quer a evidência, interrompa e peça `Use spf.implement para executar a etapa de QA independente de <slug>`, depois arquive de novo. Caso contrário, o aviso é apenas informativo.

### Falha de hook bloqueia uma tarefa

**Sintoma:** o agente não consegue marcar uma tarefa como finalizada e mostra a saída de um hook.

**Por quê:** um hook de `.specforce/config.yaml` terminou com código diferente de zero. Os hooks rodam sem shell (cada entrada é dividida nos espaços), em paralelo, então `|`, `&&` e redirecionamentos nunca são interpretados. Apenas um lote só de tarefas `[RED]` pula os hooks `on_task_finished`.

**O que fazer:** peça `Corrija a causa do hook que está falhando; não o pule nem o remova`. Se a própria entrada do hook estiver errada, edite `.specforce/config.yaml` para que ela chame um script ou um alvo do `make` (veja [Configuração](configuration.md)).

### Agente quer codar sem especificação

**Sintoma:** você pede uma mudança e o agente se recusa a editar o código, ou para no meio da implementação pedindo uma atualização da especificação.

**Por quê:** isso é intencional. O código só é alterado por meio de uma especificação aprovada, e uma mudança de escopo ou de arquitetura precisa passar pelo planejamento primeiro.

**O que fazer:** peça `Use spf.spec para planejar: <a mudança>`, ou `Use spf.discovery para explorar: <a ideia>` quando ela ainda não estiver clara.

### Worker retorna uma pergunta BLOCKED

**Sintoma:** durante a implementação, o agente para e faz uma pergunta com opções e uma recomendação.

**Por quê:** o worker que executa um lote não pode perguntar diretamente a você. Quando precisa de uma decisão, ele para com um relatório BLOCKED e a sessão principal repassa a pergunta para você.

**O que fazer:** responda no chat, escolhendo uma opção ou a recomendação. Se a sua resposta mudar o escopo, espere que o agente atualize a especificação com o `spf.spec` primeiro.

### Redimensionamento durante a implementação

**Sintoma:** depois de um redimensionamento, a implementação se recusa a continuar e lista artefatos ausentes.

**Por quê:** promover uma especificação (por exemplo para `large`) marca os artefatos que passaram a ser obrigatórios como ausentes, e a implementação fica bloqueada até que eles existam. Rebaixar nunca apaga arquivos.

**O que fazer:** peça `Use spf.spec para escrever os artefatos que faltam em <slug>`. Depois de aprová-los, peça `Use spf.implement para continuar <slug>`.

### Editar uma especificação aprovada

**Sintoma:** você quer mudar requisitos ou tarefas de uma especificação que já aprovou.

**Por quê:** uma especificação aprovada é um contrato. Mudá-la significa atualizar todos os artefatos relacionados juntos (requisitos, design, delta, tarefas) e aprová-la de novo; editar apenas o `tasks.md` deixa a especificação inconsistente.

**O que fazer:** peça `Use spf.spec para atualizar <slug>: <a mudança>. Mantenha todos os artefatos consistentes`, revise, aprove e então siga o [fluxo 5](#5-mudar-depois-de-implementar-antes-de-arquivar).
