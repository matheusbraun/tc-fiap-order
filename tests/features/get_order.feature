# language: pt
Funcionalidade: Obter Pedido
  Como um usuário
  Eu quero buscar um pedido por ID
  Para visualizar seus detalhes

  Contexto:
    Dado que o serviço de clientes está disponível
    E o serviço de produtos está disponível

  Cenário: Buscar pedido existente com sucesso
    Dado que existe um pedido com ID 1
    E o pedido pertence ao cliente 1
    E o cliente com ID 1 existe com nome "João Silva"
    E o pedido contém o produto 101
    E o produto com ID 101 existe com nome "Pizza"
    Quando eu buscar o pedido 1
    Então o pedido deve ser retornado
    E o pedido deve conter dados do cliente "João Silva"
    E o pedido deve conter dados do produto "Pizza"

  Cenário: Buscar pedido que não existe
    Quando eu buscar o pedido 999
    Então deve retornar erro de pedido não encontrado
