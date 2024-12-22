# Сервис для вычисления арифметических выражений

## Описание

Сервис предназначен для вычисления арифметических выражений, отправленных через HTTP POST запрос. Пользователь отправляет выражение в теле запроса, а сервис возвращает результат вычисления или сообщение об ошибке.

## Как запустить

1. Убедитесь, что у вас установлен Go версии 1.19 или выше.

2. Клонируйте этот репозиторий:
   ```cmd
   git clone <repository_url>
   cd <repository_directory>
   ```

3. Запустите сервис с помощью команды:
   ```cmd
   go run main.go
   ```

4. Сервис будет доступен по адресу: [http://localhost:8080](http://localhost:8080)

## Использование

Сервис предоставляет один эндпоинт:

- **URL**: `/api/v1/calculate`
- **Метод**: `POST`
- **Тело запроса** (JSON):
  ```json
  {
    "expression": "арифметическое выражение"
  }
  ```

Пример запроса с использованием `curl` в Windows:

### Однострочная команда
```cmd
curl --location http://localhost:8080/api/v1/calculate --header "Content-Type: application/json" --data "{ \"expression\": \"2+2*2\" }"
```

### Многострочная команда с использованием `^` для переноса строк
```cmd
curl --location http://localhost:8080/api/v1/calculate ^
--header "Content-Type: application/json" ^
--data "{ \"expression\": \"2+2*2\" }"
```

### Возможные ответы

#### Успешное вычисление
- **Код ответа**: `200 OK`
- **Пример тела ответа**:
  ```json
  {
    "result": "6"
  }
  ```

#### Некорректное выражение
- **Код ответа**: `422 Unprocessable Entity`
- **Пример тела ответа**:
  ```json
  {
    "error": "Expression is not valid"
  }
  ```

#### Внутренняя ошибка сервера
- **Код ответа**: `500 Internal Server Error`
- **Пример тела ответа**:
  ```json
  {
    "error": "Internal server error"
  }
  ```

## Примеры запросов

- Успешный запрос:
  ```cmd
  curl --location http://localhost:8080/api/v1/calculate --header "Content-Type: application/json" --data "{ \"expression\": \"10/2+5\" }"
  ```

  Ответ:
  ```json
  {
    "result": "10"
  }
  ```

- Некорректное выражение:
  ```cmd
  curl --location http://localhost:8080/api/v1/calculate --header "Content-Type: application/json" --data "{ \"expression\": \"2+2a\" }"
  ```

  Ответ:
  ```json
  {
    "error": "Expression is not valid"
  }
  ```

- Внутренняя ошибка сервера:
  Чтобы смоделировать, можно отправить некорректные данные, которые вызывают ошибку обработки.

## Тестирование

Рекомендуется использовать `curl`, Postman или аналогичный инструмент для проверки работы сервиса. Проверьте все сценарии: корректные выражения, некорректные данные и симуляцию внутренних ошибок.


