// @title	Task API
// @version 1.0
// @description Это программа для создания и отслеживания задач

// @contact.name Шукурилло
// @contact.email karimov199872@gmail.com

// @host localhost:8181

// @securityDefinitions.basic BasicAuth
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
	"to_do_list/config"
	"to_do_list/models"
	"to_do_list/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка при загрузке переменных окружения")
	}

	config.ConnectDatabase()

	if err := config.DB.AutoMigrate(&models.User{}, &models.Task{}); err != nil {
		log.Fatal("Не удалось выполнить миграции")
	}

	r := gin.Default()

	routes.AuthRoutes(r)
	routes.TaskRoutes(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run("localhost:8181"); err != nil {
		log.Fatal(err)
	}
}
