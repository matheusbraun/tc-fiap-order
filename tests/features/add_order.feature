# language: pt
Funcionalidade: Adicionar Pedido
  Como um cliente
  Eu quero criar um pedido com produtos
  Para que eu possa comprar itens

  Contexto:
    Dado que o serviço de clientes está disponível
    E o serviço de produtos está disponível

  Cenário: Criar pedido com sucesso com um único produto
    Dado que o cliente com ID 1 existe
    E o produto com ID 101 existe com preço 25.50
    Quando eu criar um pedido para o cliente 1 com os seguintes produtos:
      | product_id | quantity | price |
      | 101        | 2        | 25.50 |
    Então o pedido deve ser criado com sucesso
    E o total do pedido deve ser 51.00
    E o status do pedido deve ser 1

  Cenário: Criar pedido com sucesso com múltiplos produtos
    Dado que o cliente com ID 1 existe
    E o produto com ID 101 existe com preço 25.50
    E o produto com ID 102 existe com preço 15.00
    Quando eu criar um pedido para o cliente 1 com os seguintes produtos:
      | product_id | quantity | price |
      | 101        | 1        | 25.50 |
      | 102        | 3        | 15.00 |
    Então o pedido deve ser criado com sucesso
    E o total do pedido deve ser 70.50

  Cenário: Falhar ao criar pedido quando cliente não existe
    Dado que o cliente com ID 999 não existe
    Quando eu criar um pedido para o cliente 999 com produto 101
    Então a criação do pedido deve falhar
    E o erro deve conter "customer 999 not found"

  Cenário: Falhar ao criar pedido quando produto não existe
    Dado que o cliente com ID 1 existe
    E o produto com ID 999 não existe
    Quando eu criar um pedido para o cliente 1 com produto inexistente 999
    Então a criação do pedido deve falhar
    E o erro deve conter "products not found"
