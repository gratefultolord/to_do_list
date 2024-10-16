package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
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

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен не предоставлен"})
			c.Abort()
			return
		}

		// Токен должен быть формата "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Некорректный формат токена"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Парсинг и проверка токена
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Невалидный или просроченный токен"})
			log.Println("Error validating token:", err)
			log.Println("Token:", tokenString)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Невалидные данные токена"})
			c.Abort()
			return
		}

		userID := claims["user_id"].(float64)
		c.Set("user_id", uint(userID))

		c.Next()
	}
}
