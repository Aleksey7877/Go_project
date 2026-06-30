# Go Junior Training

Домашние задания по дисциплине «Язык Go».

## Структура

### `/webinars/hw1`

Задание 1. Вывод в консоль имени пользователя, аргументов командной строки и версии Go.

---

### `/webinars/hw2/gateway`

Задание 2. Простой HTTP-сервер.

Реализован endpoint:

```text
GET /ping
```

Запуск:

```bash
cd webinars/hw2/gateway
go run .
```

Проверка:

```bash
curl http://localhost:8080/ping
```

---

### `/webinars/hw2/ledger`

Задание 3. Работа с транзакциями и бюджетами.

Реализовано:

* создание транзакций;
* хранение транзакций в памяти;
* создание и обновление бюджетов;
* чтение бюджетов из `budgets.json`;
* проверка транзакций на соответствие лимиту бюджета, категории и году.

---

### `/webinars/hw2/ledger`

Задание 4. Валидация данных через интерфейс.

Ledger дополнен интерфейсом:

```go
type Validatable interface {
    Validate() error
}
```

Метод `Validate()` реализован для транзакций и бюджетов.

Структура пакета:

```text
demo.go        - сценарии проверки работы ledger
io.go          - чтение бюджетов из JSON-файла
models.go      - структуры данных Transaction и Budget
storage.go     - бизнес-логика работы с транзакциями и бюджетами в памяти
validation.go  - интерфейс Validatable и методы Validate()
```

---

## `/webinars/hw2/gateway`

Задание 5. Минимальный REST API для Gateway с использованием бизнес-логики Ledger.

В текущей версии `ledger` используется как импортируемый пакет, а запуск приложения выполняется через `gateway`.

Запуск:

```bash
cd webinars/hw2/gateway
go run .
```

Сервер запускается на:

```text
http://localhost:8080
```

### Доступные endpoints

```text
GET  /api/budgets
POST /api/budgets

GET  /api/transactions
POST /api/transactions
```

### Проверка API

Получить список бюджетов:

```bash
curl http://localhost:8080/api/budgets
```

Создать бюджет:

```bash
curl -X POST http://localhost:8080/api/budgets \
-H "Content-Type: application/json" \
-d '{"category":"food","limit":1000,"period":"2026"}'
```

Получить список транзакций:

```bash
curl http://localhost:8080/api/transactions
```

Создать транзакцию:

```bash
curl -X POST http://localhost:8080/api/transactions \
-H "Content-Type: application/json" \
-d '{"amount":450,"category":"food","description":"lunch","date":"2026-09-10"}'
```

Проверка превышения бюджета:

```bash
curl -i -X POST http://localhost:8080/api/transactions \
-H "Content-Type: application/json" \
-d '{"amount":10000,"category":"food","description":"expensive dinner","date":"2026-09-10"}'
```

Ожидаемый ответ:

```json
{
  "error": "budget exceeded"
}
```

Статус:

```text
409 Conflict
```

### Структура Gateway

```text
main.go                  - запуск HTTP-сервера

internal/api/dto.go       - DTO-модели запросов и ответов
internal/api/handlers.go  - HTTP-обработчики
internal/api/mapper.go    - преобразование DTO в модели ledger и обратно
internal/api/middleware.go - middleware для логирования метода, пути и времени обработки
internal/api/response.go  - единый JSON-ответ и JSON-ошибки
internal/api/router.go    - регистрация маршрутов
```

### Особенности реализации

* Gateway принимает и возвращает JSON.
* Для ошибок используется единый формат:

```json
{
  "error": "message"
}
```

* Для успешного создания транзакций и бюджетов возвращается статус `201 Created`.
* Ошибки валидации возвращаются со статусом `400 Bad Request`.
* Превышение бюджета возвращается со статусом `409 Conflict`.
* Внутренняя бизнес-логика находится в пакете `ledger`.
* Middleware логирует HTTP-метод, путь и длительность обработки запроса.


### Задание 6. Тесты.


Для запуска тестов Ledger:


cd webinars/hw2/ledger
go test -v -cover -coverprofile=cover.out
go tool cover -html=cover.out -o cover.html


Для запуска тестов Gateway:

cd webinars/hw2/gateway
go test -v ./... -cover -coverprofile=cover.out
go tool cover -html=cover.out -o cover.html

Текущее покрытие:

Ledger: ~61.4%
Gateway internal/api: ~70.7%