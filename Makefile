.PHONY: build run test docker-build docker-run

# Сборка приложения
build:
	go build -o bin/api ./cmd/api

# Запуск приложения
run: build
	./bin/api

# Запуск тестов
test:
	go test -v ./...

# Сборка Docker образа
docker-build:
	docker build -t notes-service .

# Запуск Docker контейнеров
docker-run:
	docker-compose up --build

# Остановка Docker контейнеров
docker-stop:
	docker-compose down

# Применение миграций (пример, требуется установка migrate tool)
migrate-up:
	migrate -path migrations -database "postgres://user:password@localhost:5432/notesdb?sslmode=disable" up

# Откат миграций
migrate-down:
	migrate -path migrations -database "postgres://user:password@localhost:5432/notesdb?sslmode=disable" down