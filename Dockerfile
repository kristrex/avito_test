# Используем официальный образ Go
FROM golang:1.23-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы проекта в контейнер
COPY . .

# Устанавливаем зависимости
RUN go mod download

WORKDIR /app/src/main

# Собираем приложение
RUN go build -o main .

WORKDIR /app

RUN cp /app/src/main/main /app/main

# Экспонируем порт, на котором будет работать приложение
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]