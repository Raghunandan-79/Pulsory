package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("authorization")

		if header == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header missing",
			})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			ctx.Abort()
			return 
		}

		claims := token.Claims.(jwt.MapClaims)

		ctx.Set("userId", claims["userId"])

		ctx.Next()
	}
}