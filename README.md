# Notes Service

Cервис на Go для управления заметками с REST API, интеграцией Яндекс.Спеллер и системой аутентификации.

## Требования

- Docker и Docker Compose
- Make (опционально)

## Установка и запуск

1. Клонируйте репозиторий:
   ```
   git clone https://github.com/alinurmyrzakhanov/notes-service
   cd notes-service
   ```

2. Создайте файл `.env` в корне проекта и добавьте необходимые переменные окружения:
   ```
   DATABASE_URL=postgres://user:password@db:5432/notesdb?sslmode=disable
   JWT_SECRET=your-secret-key
   ```

3. Запустите приложение с помощью Docker Compose:
   ```
   make docker-run
   ```
   или
   ```
   docker-compose up --build
   ```

4. Приложение будет доступно по адресу `http://localhost:8080`

## API Endpoints

- `POST /register`: Регистрация нового пользователя
```
curl -X POST http://localhost:8080/register -H "Content-Type: application/json" -d '{
  "username": "admin",
  "password": "admin"
}'
```

- `POST /login`: Вход пользователя и получение JWT токена
```
curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{
  "username": "admin",
  "password": "admin"
}' 
```

- `POST /notes`: Создание новой заметки (требуется аутентификация)
```
curl -X POST http://localhost:8080/notes -H "Authorization: Bearer your-jwt-token" -H "Content-Type: application/json" -d '{
  "title": "My First Note",
  "content": "This is the content of my first note."
}'
```
- `GET /notes`: Получение списка заметок пользователя (требуется аутентификация)
```
curl -X GET http://localhost:8080/notes -H "Authorization: Bearer your-jwt-token"
```
## Разработка

- Для сборки приложения: `make build`
- Для запуска тестов: `make test`
- Для локального запуска: `make run`

## Структура проекта

- `cmd/api`: Точка входа в приложение
- `internal`: Внутренние пакеты приложения
  - `auth`: Аутентификация и авторизация
  - `config`: Конфигурация приложения
  - `handlers`: Обработчики HTTP-запросов
  - `models`: Модели данных
  - `repository`: Работа с базой данных
  - `spellcheck`: Интеграция с Яндекс.Спеллер
- `migrations`: SQL-скрипты для миграций базы данных
- `tests`: Автотесты

