package routes

import (
	"github.com/gin-gonic/gin"
	"to_do_list/controllers"
	"to_do_list/middlewares"
)

func TaskRoutes(r *gin.Engine) {
	tasks := r.Group("/tasks")
	tasks.Use(middlewares.AuthMiddleware())

	tasks.POST("/", controllers.CreateTask)
	tasks.GET("/", controllers.GetTasks)
	tasks.PUT("/:id", controllers.UpdateTask)
	tasks.DELETE("/:id", controllers.DeleteTask)
}
