package db

import (
	"fmt"
	"log"
	"os"
	"zadanie-6105/src/main/models"
)

func MigrateDB() error {
	// Открываем файл migration.sql
	file, err := os.Open("./src/main/db/migration.sql")
	if err != nil {
		return fmt.Errorf("error opening migration.sql file: %v", err)
	}
	defer file.Close()

	// Читаем содержимое файла
	sqlBytes, err := os.ReadFile("./src/main/db/migration.sql")
	if err != nil {
		return fmt.Errorf("error reading migration.sql file: %v", err)
	}

	// Преобразуем содержимое файла в строку
	sqlQuery := string(sqlBytes)

	// Получаем объект sql.DB из gorm.DB
	sqlDB, err := models.Db.DB()
	if err != nil {
		return fmt.Errorf("error getting sql.DB object: %v", err)
	}

	// Выполняем SQL-запросы
	_, err = sqlDB.Exec(sqlQuery)
	if err != nil {
		return fmt.Errorf("error executing SQL query: %v", err)
	}

	log.Println("Migration completed successfully")
	return nil
}
