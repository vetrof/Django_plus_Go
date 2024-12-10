package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os"
	"sync"
)

type Product struct {
	ID          int
	Name        string
	Description string
	Price       int
}

// Инициализация пула соединений
var pool *pgxpool.Pool

// Инициализация базы данных
func initDB() {
	dsn := "postgres://test_user:QWErty42@db:5432/test_pgdb"
	var err error
	// Создаем пул соединений
	pool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
}

func getShopsHandler(c *gin.Context) {
	// Пытаемся получить соединение из пула
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Failed to acquire connection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer conn.Release()

	// Выполнение SQL-запроса
	rows, err := conn.Query(context.Background(), "SELECT * FROM catalog_catalog")
	if err != nil {
		log.Printf("Query failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query failed"})
		return
	}
	defer rows.Close()

	// Создаем WaitGroup для ожидания завершения горутин
	var wg sync.WaitGroup
	var mu sync.Mutex // Защита от гонок данных
	var products []Product

	// Обработка строк из результата запроса в горутинах
	for rows.Next() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var product Product
			err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price)
			if err != nil {
				log.Printf("Failed to scan row: %v\n", err)
				return
			}
			// Защищаем от гонок при добавлении в срез
			mu.Lock()
			products = append(products, product)
			mu.Unlock()
		}()
	}

	// Ожидаем завершения всех горутин
	wg.Wait()

	// Проверка на отсутствие данных
	if len(products) == 0 {
		log.Println("No products found.")
		c.JSON(http.StatusNotFound, gin.H{"error": "No products found"})
		return
	}

	// Отправка ответа
	c.JSON(http.StatusOK, products)
}

func emptyHandler(c *gin.Context) {
	// Отправка ответа
	var products string
	c.JSON(http.StatusOK, products)
}

func main() {
	// Инициализация базы данных
	initDB()

	// Создание роутера
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Регистрация маршрута
	router.GET("/product", getShopsHandler)
	router.GET("/empty", emptyHandler)

	// Запуск сервера
	port := os.Getenv("GO_SERVICE_PORT")
	if port == "" {
		port = "8080" // Установите порт по умолчанию
	}
	log.Printf("Server running on port %s", port)
	router.Run(":" + port)
}
