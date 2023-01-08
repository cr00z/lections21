# Задание

Написать веб-сервер чата (без клиентской части) работающий по REST-архитектуре.

- Подразумевается регистрация и аутентификация пользователей (например /users/register).
- Пользователи должны иметь возможность писать в общий чат и получать сообщения из общего чата (например /messages).
Сообщения в общем чате должны идентифицировать пользователя для возможности отправки личных сообщений.
 - Пользователи должны иметь возможность отправлять/получать личные сообщения (например /users/{id}/messages или /users/me/messages).
Читать личные сообщения других пользователей, конечно-же, запрещено.

Конкретная имплементация общего чата в концепте RESTfull остается на усмотрение разработчика.
Хранение данных подразумевается in-memory.

Плюсами будут:
- пагинация сообщений
- возможность безболезненно переключить хранение данных на другой тип, например DB
- документация ручек

# Запуск

`go run cmd/server/main.go`

```
2022/12/28 10:13:09 server started
```

## Регистрация

`curl -v -H "Content-Type: application/json" -d '{"username": "test1", "password": "111"}' -XPOST localhost:5000/register`

```
{
 "id": 3
}
```
`curl -v -H "Content-Type: application/json" -d '{"username": "test2", "password": "111"}' -XPOST localhost:5000/register`
```
{
 "id": 4
}
```

## Авторизация

`curl -v -H "Content-Type: application/json" -d '{"username": "test1", "password": "111"}' -XPOST localhost:5000/login`
```
{
 "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzIyMDEwNjUsImlhdCI6MTY3MjE5NzQ2NSwiVXNlcklEIjozfQ.L3OUBIS5CGGOX0LrQFs-Zo00dhzLx8mGdib5kaGuhRU"
}
```
`curl -v -H "Content-Type: application/json" -d '{"username": "test2", "password": "111"}' -XPOST localhost:5000/login`
```
{
 "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzIyMDExMTUsImlhdCI6MTY3MjE5NzUxNSwiVXNlcklEIjo0fQ.B3B0Hfq8PCLKzX0ucknC-H--mIOoYPU321pJw5PEsts"
}
```

## Запрос сообщений чата

`curl -v -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzIyMDEwNjUsImlhdCI6MTY3MjE5NzQ2NSwiVXNlcklEIjozfQ.L3OUBIS5CGGOX0LrQFs-Zo00dhzLx8mGdib5kaGuhRU" -XGET localhost:5000/api/messages`

```
[
 {
  "id": 0,
  "from_user": 1,
  "to_user": 0,
  "text": "hello all"
 },
 {
  "id": 1,
  "from_user": 1,
  "to_user": 0,
  "text": "im User1"
 },
 {
  "id": 2,
  "from_user": 2,
  "to_user": 0,
  "text": "hello, User1"
 }
]
```

## Отправка сообщений в общий чат

`curl -v -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzIyMDEwNjUsImlhdCI6MTY3MjE5NzQ2NSwiVXNlcklEIjozfQ.L3OUBIS5CGGOX0LrQFs-Zo00dhzLx8mGdib5kaGuhRU" -d '{"text":"test text"}' -XPOST localhost:5000/api/messages`
```
{
 "status": "ok"
}
...
// запрос сообщений
...
 {
  "id": 4,
  "from_user": 3,
  "to_user": 0,
  "text": "test text"
 }
]
```

## Отправка личных сообщений

`curl -v -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzIyMDEwNjUsImlhdCI6MTY3MjE5NzQ2NSwiVXNlcklEIjozfQ.L3OUBIS5CGGOX0LrQFs-Zo00dhzLx8mGdib5kaGuhRU" -d '{"to_user":4, "text":"private message1"}' -XPOST localhost:5000/api/users/4/messages`
```
{
 "status": "ok"
}
```
`curl -v -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzIyMDEwNjUsImlhdCI6MTY3MjE5NzQ2NSwiVXNlcklEIjozfQ.L3OUBIS5CGGOX0LrQFs-Zo00dhzLx8mGdib5kaGuhRU" -d '{"to_user":4, "text":"private message2"}' -XPOST localhost:5000/api/users/4/messages`
```
{
 "status": "ok"
}
```

## Получение личных сообщений

`curl -v -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzIyMDExMTUsImlhdCI6MTY3MjE5NzUxNSwiVXNlcklEIjo0fQ.B3B0Hfq8PCLKzX0ucknC-H--mIOoYPU321pJw5PEsts" -XGET localhost:5000/api/users/me/messages`
```
[
 {
  "id": 6,
  "from_user": 3,
  "to_user": 4,
  "text": "private message1"
 },
 {
  "id": 7,
  "from_user": 3,
  "to_user": 4,
  "text": "private message2"
 }
]
```