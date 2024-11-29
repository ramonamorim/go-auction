# go-auction : 

Este projeto é um sistema de leilão online desenvolvido em Go. O sistema permite criar leilões, fazer lances e verificar o vencedor. O sistema também possui uma funcionalidade de fechamento automático de leilões após um tempo especificado.

## Pré-requisitos

- Docker
- Docker Compose

## **Build e start dos serviços**
  Para rodar os servicos:
  ```sh
  docker-compose up --build
  ```

## **Parar os serviços**

  Utilize o seguinte comando:
  ```sh
  docker-compose down
  ```


## A aplicação expõe os seguintes endpoints:

- `GET /auction`: Lista todos os leilões
- `GET /auction/:auctionId`: Busca um leilão pelo ID
- `POST /auction`: Cria um novo leilão
- `GET /auction/winner/:auctionId`: Busca o lance vencedor de um leilão pelo ID
- `POST /bid`: Cria um novo lance
- `GET /bid/:auctionId`: Lista todos os lances de um leilão
- `GET /user/:userId`: Busca um usuário pelo ID


## Para rodar os testes, siga os passos abaixo:

1. Garanta que o MongoDB está em execução na sua máquina local.
2. Execute o comando a seguir para rodar os testes:

    ```sh
    go test ./internal/infra/database/auction -v
    ```

## Para testar a API:

Após o contêiner estar em execução, você pode testar a API usando os seguintes comandos curl:

1. Para criar um novo leilão:
```
curl -X POST http://localhost:8080/auction \
-H "Content-Type: application/json" \
-d '{
"product_name": "iPhone ",
"category": "celular",
"description": "aparelho celular smart.",
"condition": 0
}'
```

Resposta esperada: Um código de status 204 No Content, sem corpo de resposta, indicando que a solicitação foi processada com sucesso, sem erros.

2. Para listar leilões com status concluído:

Observe que você deve aguardar o tempo definido na variável de ambiente para que os leilões sejam automaticamente marcados como concluídos.

```
curl http://localhost:8080/auction?status=1
```
Resposta esperada: Um objeto JSON contendo os dados do(s) leilão(ões) criado(s).