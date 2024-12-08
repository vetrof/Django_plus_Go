package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"context"
	"github.com/gin-gonic/gin"
)

type Shop struct {
	ID        int     `json:"id"`
	ShopID    string  `json:"shop_id"`
	Title     string  `json:"title"`
	City      string  `json:"city"`
	Address   string  `json:"address"`
	Enabled   bool    `json:"enabled"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Points    string  `json:"points"`
}

func getShopsHandler(c *gin.Context) {
	// Формируем строку подключения из ENV
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)

	// Логируем строку подключения
	log.Printf("Connecting to database with DSN: %s", dsn)

	// Подключение к базе данных
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer conn.Close(context.Background())

	// Выполнение SQL-запроса с JOIN для получения информации о городе
	log.Println("Executing query to fetch shops with city...")
	rows, err := conn.Query(context.Background(), `
        SELECT shop.shop_id, shop.title, city.name AS city, shop.address,
               shop.enabled, shop.latitude, shop.longitude, shop.points
        FROM catalog_shop shop
        JOIN catalog_city city ON shop.city_id = city.id
    `)
	if err != nil {
		log.Printf("Query failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query failed"})
		return
	}
	defer rows.Close()

	// Чтение данных
	var shops []Shop
	for rows.Next() {
		var shop Shop
		err = rows.Scan(&shop.ShopID, &shop.Title, &shop.City, &shop.Address, &shop.Enabled, &shop.Latitude, &shop.Longitude, &shop.Points)
		if err != nil {
			log.Printf("Failed to scan row: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		log.Printf("Fetched shop: %+v", shop) // Логируем информацию о каждом магазине
		shops = append(shops, shop)
	}

	// Проверка на отсутствие данных
	if len(shops) == 0 {
		log.Println("No shops found.")
		c.JSON(http.StatusNotFound, gin.H{"error": "No shops found"})
		return
	}

	// Отправка ответа
	log.Println("Returning shops data.")
	c.JSON(http.StatusOK, shops)
}

func main() {
	// Создание роутера
	r := gin.Default()

	// Регистрация маршрута
	r.GET("/shops", getShopsHandler)

	// Запуск сервера
	port := os.Getenv("GO_SERVICE_PORT")
	if port == "" {
		port = "8080" // Установите порт по умолчанию
	}
	log.Printf("Server running on port %s", port)
	r.Run(":" + port)
}
