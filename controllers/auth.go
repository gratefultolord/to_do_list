package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
	"to_do_list/config"
	"to_do_list/models"
)

var jwtSecret []byte

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка при загрузке переменных окружения: ", err)
	}
	jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(jwtSecret) == 0 {
		log.Fatal("JWT_SECRET_KEY не установлен или пустой")
	}
}

// AuthDetails структура для передачи данных аутентификации
type AuthDetails struct {
	Username string `json:"username" binding:"required" example:"userman"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// SuccessResponse структура для успешного ответа
type SuccessResponse struct {
	Message string `json:"message" example:"Пользователь успешно зарегистрирован"`
}

// TokenResponse структура получения токена
type TokenResponse struct {
	Token string `json:"token" example:"your_token"`
}

// Register godoc
// @Summary 	Регистрация пользователя
// @Description Регистрирует пользователя с указанными username и password
// @Tags 		Аутентификация
// @Accept 		json
// @Produce 	json
// @Param		body body models.User true "Данные пользователя для регистрации"
// @Success 	201 {object} SuccessResponse
// @Failure		400 {object} gin.H
// @Failure		500 {object} gin.H
// @Router		/register [post]
func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("Полученные данные при регистрации:", user)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при хешировании пароля"})
		return
	}
	user.Password = string(hashedPassword)

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании пользователя"})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{Message: "Пользователь успешно зарегистрирован"})
}

// Login godoc
// @Summary 	Авторизация
// @Description Авторизация пользователя и выдача ему токена
// @tags 		Аутентификация
// @Accept		json
// @Produce		json
// @Param		body body AuthDetails true "Данные пользователя для авторизации"
// @Success 	200 {object} TokenResponse
// @Failure 	400 {object} gin.H
// @Failure		401 {object} gin.H
// @Failure		500 {object} gin.H
// @Router		/login [post]
func Login(c *gin.Context) {
	var authDetails AuthDetails
	var user models.User

	if err := c.ShouldBindJSON(&authDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Поиск пользователя по имени
	if err := config.DB.Where("username = ?", authDetails.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неправильное имя пользователя"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authDetails.Password)); err != nil {
		log.Println("LOOK HERE", user.Password, authDetails.Password)
		log.Println("Error while comparing passes:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный пароль"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации токена"})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{Token: tokenString})
}
