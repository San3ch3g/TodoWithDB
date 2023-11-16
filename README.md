# TodoWithDB
 Это простое приложение на языке программирования Go, представляющее собой пример использования базы данных MySQL для хранения списка задач (todo list). Приложение предоставляет API для получения, добавления, обновления и удаления задач в базе данных.
 ## Использование
 ### Установка зависимостей
 Перед запуском приложения убедитесь, что у вас установлены необходимые зависимости:
 ```go
 go get -u github.com/rs/cors
 go get -u github.com/go-sql-driver/mysql
```
### Конфигурация базы данных
Замените значения в строке подключения к базе данных в методах data, update, createReq, и delete согласно вашим настройкам MySQL.
```go
db, err := sql.Open("mysql", "your_user:your_password@tcp(host:port)/your_data_base")
```
### Запуск приложения
```go
go run main.go
```
Приложение будет доступно по адресу http://localhost:8080 .
## API
### Получение списка задач
* URL: /data
* Метод: GET
```json
[
    {"id": 1, "title": "Задача 1", "completed": true},
    {"id": 2, "title": "Задача 2", "completed": false},
]
```
### Обновление задачи
* URL: /update
* Метод: POST

```json
[
    {"id": 1, "title": "Обновленная задача", "completed": false}
]
```
### Добавление новой задачи
* URL: /create
* Метод: POST

```json
[
    {"title": "Новая задача", "completed": true}
]
```
### Удаление задачи
* URL: /delete
* Метод: POST

```json
[
    {"id": 1}
]
```
