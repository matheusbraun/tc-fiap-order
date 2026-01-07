# language: pt
Funcionalidade: Listar Pedidos
  Como um usuário
  Eu quero listar todos os pedidos ativos
  Para visualizar o status atual

  Contexto:
    Dado que o serviço de clientes está disponível
    E o serviço de produtos está disponível

  Cenário: Listar pedidos ativos
    Dado que existem 3 pedidos não finalizados
    E existe 1 pedido finalizado
    Quando eu listar os pedidos
    Então deve retornar 3 pedidos
    E nenhum pedido finalizado deve ser retornado

  Cenário: Listar quando não há pedidos
    Dado que não existem pedidos
    Quando eu listar os pedidos
    Então deve retornar uma lista vazia
