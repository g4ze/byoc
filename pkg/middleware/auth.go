package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/g4ze/byoc/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

// verifies the jwt and
// appends the user attributed
// to the token to the context params
func JwtMiddleware() gin.HandlerFunc {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	var JwtKey = []byte(os.Getenv("JWT_SECRET"))
	if JwtKey == nil {
		log.Fatalf("JWT_SECRET not found in .env")
	}
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization Header"})
			c.Abort()
			return
		}
		tokenString = tokenString[len("Bearer "):]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		log.Printf("Token: %+v", token)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"validity": token.Valid, "err": err})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "claims not ok"})
		}

		userName, ok := claims["userName"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "userName not ok"})
		}

		c.Params = append(c.Params, gin.Param{Key: "user", Value: userName})
		c.Next()
	}
}

var limiter = rate.NewLimiter(1, 5) // 1 request per second, burst of 5

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}
func GenerateJWT(password, userName string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	user, err := database.GetUser(userName, password)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", nil
	}

	var KEY = []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"password": password,
		"userName": userName,
	})
	tokenString, err := token.SignedString(KEY)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
