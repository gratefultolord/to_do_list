package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"to_do_list/config"
	"to_do_list/models"
)

// CreateTaskInput структура для создания задачи
type CreateTaskInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// CreateTask godoc
// @Summary Создание задачи
// @Description Создает задачу
// @tags Задачи
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен"
// @Param body body CreateTaskInput true "Данные для создания задачи"
// @Success 201 {object} models.Task
// @Failure 401 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /tasks [post]
func CreateTask(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	var input CreateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		UserID:      userID.(uint),
		Status:      models.StatusStarted,
	}

	//Сохраняем задачу в БД
	if err := config.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании задачи"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"task": task})
}

// GetTasks godoc
// @Summary Список задач
// @Description Даёт список задач пользователя
// @tags Задачи
// @Param Authorization header string true "Bearer токен"
// @Success 200 {array} models.Task
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /tasks [get]
func GetTasks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	var tasks []models.Task
	if err := config.DB.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ошибка при получении задач"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// UpdateTaskInput структура задачи для обновления
type UpdateTaskInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// UpdateTask godoc
// @Summary Обновление задачи
// @Description Обновляет название, описание и статус задачи
// @tags Задачи
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен"
// @Param id path int true "ID задачи"
// @Param body body UpdateTaskInput true "Данные задачи для обновления"
// @Success 200 {object} models.Task
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /tasks/{id} [put]
func UpdateTask(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	var task models.Task
	if err := config.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
		return
	}

	var input UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Обновляем задачу
	if input.Title != "" {
		task.Title = input.Title
	}

	if input.Description != "" {
		task.Description = input.Description
	}

	if input.Status != "" {
		task.Status = models.Status(input.Status)
	}

	// Сохраняем изменения в БД
	if err := config.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении задачи"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// DeleteTask godoc
// @Summary Удаление задачи
// @Description Удаляет задачу
// @tags Задачи
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен"
// @Param id path int true "ID задачи"
// @Success 204 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
func DeleteTask(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	// Получаем задачу по ID из URL
	var task models.Task
	if err := config.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
		return
	}

	// Удаляем задачу
	if err := config.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении задачи"})
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Задача успешно удалена"})
}
