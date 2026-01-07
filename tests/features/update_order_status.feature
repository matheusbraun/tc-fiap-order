# language: pt
Funcionalidade: Atualizar Status do Pedido
  Como um operador
  Eu quero atualizar o status de um pedido
  Para refletir o progresso do pedido

  Cenário: Atualizar status de Recebido para Em Preparação
    Dado que existe um pedido com ID 1 com status 1
    Quando eu atualizar o status do pedido 1 para 2
    Então o status deve ser atualizado com sucesso
    E o pedido deve ter status 2

  Cenário: Atualizar status de Em Preparação para Pronto
    Dado que existe um pedido com ID 1 com status 2
    Quando eu atualizar o status do pedido 1 para 3
    Então o status deve ser atualizado com sucesso

  Cenário: Atualizar status de Pronto para Finalizado
    Dado que existe um pedido com ID 1 com status 3
    Quando eu atualizar o status do pedido 1 para 4
    Então o status deve ser atualizado com sucesso
