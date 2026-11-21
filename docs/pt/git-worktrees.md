[ [English](../en/git-worktrees.md) | [Português](git-worktrees.md) | [Español](../es/git-worktrees.md) ]

# Suporte a Git Worktree

## Introdução

O Desenvolvimento Orientado por Especificação frequentemente envolve trabalhar em várias funcionalidades ou correções de bugs em paralelo, em diferentes branches git. O Specforce simplifica esse fluxo oferecendo suporte nativo a **Git Worktrees**. Esse recurso permite que o console TUI do Specforce descubra e exiba especificações e progresso de worktrees vinculadas, oferecendo uma visão unificada do status de todo o projeto sem a necessidade de trocar de branch constantemente.

## Mecanismo de Descoberta

O Specforce inclui um Mecanismo de Descoberta Unificado que detecta automaticamente as git worktrees vinculadas ao repositório atual. Esse mecanismo varre o diretório `.git/worktrees` para identificar outras worktrees ativas e seus respectivos caminhos.

O principal benefício dessa integração é a possibilidade de ver especificações de várias branches (por exemplo, `main`, `feature/X`, `bugfix/Y`) simultaneamente dentro do console TUI do Specforce. Isso dá aos desenvolvedores uma visão global de todo o trabalho em andamento no projeto, facilitando a coordenação e evitando esforços duplicados.

## Restrições

Para garantir a integridade dos dados e evitar modificações acidentais entre branches, o Specforce impõe as seguintes restrições:

- **Acesso somente leitura:** Todas as especificações descobertas em git worktrees externas (as que não são o seu diretório de trabalho atual) são estritamente **somente leitura**.
- **Ações bloqueadas:** Para essas especificações externas, todas as ações que modificam estado são bloqueadas. Isso inclui assumir ou finalizar tarefas, mudar o estado das especificações e quaisquer edições automáticas de arquivos.
- **Isolamento da Constituição:** O Specforce ignora o diretório `.specforce/docs/` (a Constituição do Projeto) das worktrees externas. Isso garante que a sua implementação atual sempre siga os princípios de arquitetura e engenharia definidos na sua branch local.

## Escopo da Varredura

O Mecanismo de Descoberta Unificado faz uma varredura profunda das worktrees vinculadas para reunir todos os artefatos de SDD (Spec-Driven Development) relevantes. Os seguintes itens são incluídos na varredura entre worktrees:

- **Especificações ativas:** Todos os arquivos em `.specforce/specs/**` em todas as worktrees descobertas.
- **Implementações:** Quaisquer arquivos de implementação ou roadmaps atualmente referenciados por especificações ativas.
- **Especificações arquivadas:** Especificações concluídas e arquivadas anteriormente, localizadas em `.specforce/archive/**`.

Ao incluir tanto o trabalho ativo quanto o arquivado, o Specforce oferece um contexto histórico e atual da evolução do projeto em todas as branches.

## Especificações em Paralelo e Deltas de Módulo

Quando especificações em worktrees diferentes alteram a mesma especificação viva de módulo em `.specforce/docs/modules/`, a selagem do delta impede que uma sobrescreva a outra silenciosamente. Durante o planejamento, `specforce spec delta seal <slug>` registra uma impressão digital de cada bloco de módulo que o delta modifica ou remove. Se outra especificação alterar um desses blocos primeiro (por exemplo, ela é arquivada e mesclada na sua branch), `specforce spec delta check <slug>` reporta um conflito para esse item: o delta precisa ser refeito com base no módulo atual e selado de novo antes que a especificação possa ser arquivada. Veja [Deltas de módulo](artifacts.md#6-deltas-de-módulo).
