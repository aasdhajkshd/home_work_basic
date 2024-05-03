# HTTP Server in Go

- [Understanding HTTP Server in Go](https://surajincloud.com/understanding-http-server-in-go-basic)
- [Разработка веб-серверов на Golang — от простого к сложному](https://habr.com/ru/companies/skillbox/articles/446454/)
- [Testing Golang with httptest](https://speedscale.com/blog/testing-golang-with-httptest/)

## Сторона сервера

```go
$ go run . -mode=server
Starting 'server'... by default
2024/05/02 15:06:26 for interruption, press Ctrl+C...
2024/05/02 15:07:33 Received request: GET /hello
2024/05/02 15:07:38 Received request: POST /hello
2024/05/02 15:07:45 Received request: POST /users
2024/05/02 15:07:45 Received POST content...
Reveived JSON data, decoded: {4 Dmitry Dimovich dd@mail.ru 11111111}
2024/05/02 15:07:51 Received request: GET /users
```

## Сторона клиента

### Проверка

```go
$ go run . -mode=client -url=localhost:8081 -path=hello
map[Content-Length:[30] Content-Type:[text/plain; charset=utf-8] Date:[Thu, 02 May 2024 12:07:33 GMT]]:30
Привет, мой друг!

$ go run . -mode=client -url=localhost:8081 -path=hello -method=POST
map[Content-Length:[40] Content-Type:[text/plain; charset=utf-8] Date:[Thu, 02 May 2024 12:07:38 GMT]]:40
Привет, пишете письмо?
```

### Передача HTTP GET и POST запросов от клиента с ответами с данными 

```go
$ go run . -mode=client -url=localhost:8081 -path=users -method=POST
map[Content-Length:[2] Content-Type:[text/plain] Date:[Thu, 02 May 2024 12:07:45 GMT]]:2
OK

$ go run . -mode=client -url=localhost:8081 -path=users -method=GET
map[Content-Length:[151] Content-Type:[application/json] Date:[Thu, 02 May 2024 12:07:51 GMT]]:151
[{"id":1,"name":"Vasiliy Ivanov","email":"vi@mail.ru","password":"12345678"},{"id":2,"name":"Pavel Petrov","email":"pp@mail.ru","password":"87654321"}]
```
