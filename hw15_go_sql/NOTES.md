
### Команды curl для обработчика `handlerUser`

POST запрос для добавления нового пользователя:

```bash
curl -sL -X POST http://localhost:8888/user -H "Content-Type: application/json" -d '{"name":"John Doe","email":"john@example.me","password":"12345678"}'
```

GET запрос для получения списка всех пользователей:

```bash
curl -sL -X GET http://localhost:8888/user | jq -r .
```

GET запрос для получения пользователя по id:

```bash
curl -sL -X GET "http://localhost:8888/user?id=1" | jq -r .
```

PUT запрос для обновления существующего пользователя:

```bash
curl -X PUT http://localhost:8888/user -H "Content-Type: application/json" -d '{"id":1,"name":"John Doe","email":"john.doe@example.com"}'
```

DELETE запрос для удаления пользователя:

```bash
curl -X DELETE http://localhost:8888/user -H "Content-Type: application/json" -d '{"id":1}'
```

### Команды curl для обработчика `handlerProduct`

POST запрос для добавления нового продукта:

```bash
curl -sL -X POST http://localhost:8888/product -H "Content-Type: application/json" -d '{"name":"Laptop","price":1500.00}'
```

GET запрос для получения списка всех продуктов:

```bash
curl -sL -X GET http://localhost:8888/product | jq -r .
```

GET запрос для получения продукта по id:

```bash
curl -sL -X GET "http://localhost:8888/product?id=1" | jq -r .
```

PUT запрос для обновления существующего продукта:

```bash
curl -sL -X PUT http://localhost:8888/product -H "Content-Type: application/json" -d '{"id":1,"name":"Gaming Laptop","price":2000.00}'
```

DELETE запрос для удаления продукта:

```bash
curl -sL -X DELETE http://localhost:8888/product -H "Content-Type: application/json" -d '{"id":1}'
```

### Команды curl для `handlerOrder`

POST запрос для добавления нового заказа:

```bash
curl -X POST http://localhost:8888/order -H "Content-Type: application/json" -d '{
  "userID": 1,
  "orderDate": "2024-05-30T17:19:20.92627Z",
  "totalAmount": 200.00,
  "products": {
    "list": [
      {"id": 1, "name": "Phone", "price": 500.00},
      {"id": 2, "name": "Laptop", "price": 1500.00}
    ]
  }
}'
```

GET запрос для получения заказа по id:

```bash
curl -sL -X GET "http://localhost:8888/order?id=1" | jq -r .
```

GET запрос для получения заказов пользователя по user_id:

```bash
curl -sL -X GET "http://localhost:8888/order?user_id=1"
```

DELETE запрос для удаления заказа:

```bash
curl -X DELETE http://localhost:8888/order -H "Content-Type: application/json" -d '{
  "id": 1,
  "userID": 1,
  "orderDate": "2024-05-30T17:19:20.92627Z",
  "totalAmount": 200.00,
  "products": {
    "list": [
      {"id": 1, "name": "Phone", "price": 500.00},
      {"id": 2, "name": "Laptop", "price": 1500.00}
    ]
  }
}'
curl -X DELETE http://localhost:8888/order -H "Content-Type: application/json" -d '{"id": 4}'
```
