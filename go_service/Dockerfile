# Используем официальный образ Go
FROM golang:1.23-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем остальные файлы
COPY . .

# Сборка приложения
RUN go build -o go_service .

# Запуск приложения
CMD ["./go_service"]
