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

Регистрация / авторизация

`curl -v -H "Content-Type: application/json" -d '{"username": "test", "password": "111"}' -XPOST localhost:5000/register`

`curl -v -H "Content-Type: application/json" -d '{"username": "test", "password": "111"}' -XPOST localhost:5000/login`

Запрос сообщений чата

`curl -v -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzIxNDIyMTksImlhdCI6MTY3MjEzODYxOSwiVXNlcklEIjozfQ.wtPjl4kFgQV_UE9eY9XCIP5a93JI4ACD1Xrnon80zCc" -XGET localhost:5000/api/messages`

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

Отправка сообщений в общий чат

curl -v -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzIxNDIyMTksImlhdCI6MTY3MjEzODYxOSwiVXNlcklEIjozfQ.wtPjl4kFgQV_UE9eY9XCIP5a93JI4ACD1Xrnon80zCc" -d '{"text":"test text"}' -XPOST localhost:5000/api/messages