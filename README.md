# Сервис для вычисления арифметических выражений
## Описание

Сервис предназначен для вычисления арифметических выражений, отправленных через HTTP POST запрос. Пользователь отправляет выражение в теле запроса, а сервис возвращает результат вычисления или сообщение об ошибке.

## Структура проекта
```
.
│   go.mod
│   README.md
│
├───cmd
│       main.go
│
└───internal
    ├───api
    │       api.go
    │       api_test.go
    │
    └───math
            math.go
            math_test.go
```

## Как запустить

1. Убедитесь, что у вас установлен Go версии 1.19 или выше.

2. Клонируйте этот репозиторий и перейдите в него:
   ```cmd
   git clone https://github.com/schmalz302/Calc
   cd Calc
   ```

3. Запустите сервис с помощью команды:
   ```cmd
   go run cmd/main.go
   ```

4. Сервис будет доступен по адресу: [http://localhost:8080](http://localhost:8080)

## Использование
Рекомендуется использовать `curl`, Postman или аналогичный инструмент для проверки работы сервиса. Проверьте все сценарии: корректные выражения, некорректные данные и симуляцию внутренних ошибок.

## Сценарии использования

| **Request Method** | **Endpoint** | **Request Body**                                           | **Response Body**                                    | **HTTP Status Code** |
|--------------------|--------------|------------------------------------------------------------|------------------------------------------------------|----------------------|
| POST               | `/api/v1/calculate`  | `{ "expression": "2 + 2" }`                               | `{ "result": 4 }`                                    | 200 OK               |
| POST               | `/api/v1/calculate`  | `{ "expression": "2 / 0" }`                               | `{"error": "Internal server error"}`                 | 500 Internal Server Error |
| POST               | `/api/v1/calculate`  | `{ "expression": "invalid expression" }`                  | `{ "error": "Invalid expression" }`                  | 422 Unprocessable Entity |
| POST               | `/api/v1/calculate`  | `{ "expression": }`                                       | `{ "error": "Invalid request body" }`                | 400 Bad Request      |
| POST               | `/api/v1/calculate`  | `non-json string`                                         | `{ "error": "Invalid request body" }`                | 400 Bad Request      |
| GET                | `/api/v1/calculate`  | N/A                                                       | `{ "error":"Method not allowed" }`                | 405 Method Not Allowed |

## Коды ответов
- 200: Успешное вычисление
- 422: Ошибка в выражении (неверный формат)
- 500: Внутренняя ошибка сервера
- 405: Неверный метод запроса (только POST разрешен)
## Тестирование
тесты math 
```bash
go test ./internal/math/
```
тесты api 
```bash
go test  ./internal/api/
```
тесты math с подробной информацией  
```bash
go test -v ./internal/math/
```
тесты api с подробной информацией  
```bash
go test -v ./internal/api/
```
запуск всех тестов
```bash
go test ./...
```
запуск всех тестов c подробной информацией
```bash
go -v test ./...
```
запуск тестов с информации о покрытии
```bash
go test -cover ./...
```

### Уточнения 

- файл main не тестировал, тк при запуске сервера мы не вводим порт и другие данные, валидировать особо нечего, при запуске сервера все должно отработать корректно

- покрытие тестами 95.2% - пакет math, 90.0% - пакет api
- сценарии использования оформил в виде таблице, тк это удобно

- тг ```@bll_nev_egor```
