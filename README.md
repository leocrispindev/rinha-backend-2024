# rinha-backend-2024

Este projeto demonstra como construir uma aplicação web simples usando a linguagem de programação Go, o banco de dados PostgreSQL e o servidor web Nginx. A aplicação permite que os usuários gerenciem uma lista de tarefas.

## Exemplo
A aplicação possui as seguintes funcionalidades:

### Realizar Transação
* POST /cliente/[id]transacoes
``` 
Payload:
{
    "valor": 10,
    "tipo" : "d",
    "descricao" : "EXEMPLO"
}
```
Resposta:
```
HTTP 200 OK

{
    "limite" : 100000,
    "saldo" : -100
}
```

### Consultar Extrato
* GET /cliente/[id]/extrato
```
{
  "saldo": {
    "total": -100,
    "data_extrato": "2024-02-27T03:30:40.217753Z",
    "limite": 100000
  },
  "ultimas_transacoes": [
    {
      "valor": 10,
      "tipo": "d",
      "descricao": "exemplo",
      "realizada_em": "2024-01-17T02:34:38.543030Z"
    }
  ]
}
``` 
### Execução

```
docker-compose up 
```