package main

import (
	"fmt"
	"log"
	"os"
	"zadanie-6105/src/main/controllers"
	"zadanie-6105/src/main/db"
	"zadanie-6105/src/main/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	// Загружаем переменные окружения из файла .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// Получаем параметры подключения к базе данных из переменных окружения
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USERNAME"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DATABASE"),
	)
	// Подключаемся к базе данных
	var errDB error
	models.Db, errDB = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if errDB != nil {
		log.Fatalf("Error connecting to database: %v", errDB)
	}

	// Выполняем миграцию
	err = db.MigrateDB()
	if err != nil {
		log.Fatalf("Error running migration: %v", err)
	}
}

func main() {
	r := gin.Default()

	// Регистрация эндпоинтов
	api := r.Group("/api")
	{
		api.GET("/ping", controllers.PingHandler)
		api.GET("/tenders", controllers.TenderHandler)
		tender := api.Group("/tenders")
		{
			tender.POST("/new", controllers.CreateTender)
			tender.GET("/my", controllers.MyTender)
			tender.GET("/:tenderId/status", controllers.GetStatus)
			tender.PUT("/:tenderId/status", controllers.PutStatus)
		}
		bids := api.Group("/bids")
		{
			bids.POST("/new", controllers.CreateBids)
			bids.GET("/my", controllers.MyBids)
			bids.GET(":tenderId/list", controllers.BidsList)
		}
	}

	// Запуск сервера
	log.Printf("Starting server on %s", os.Getenv("SERVER_ADDRESS"))
	if err := r.Run(os.Getenv("SERVER_ADDRESS")); err != nil {
		log.Fatal(err)
	}
}
