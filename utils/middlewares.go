package utils

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil 
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			return
		}

		c.Set("userID", claims["user_id"].(string))
		c.Set("role", claims["role"].(string))
		c.Next()
	}
}

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role.(string) != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient permissions"})
			return
		}
		c.Next()
	}
}

type client struct {
	count       int
	lastRequest time.Time
}

var (
	clients = make(map[string]*client)
	mu      sync.Mutex 
)

func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		identifier := c.ClientIP() 

		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})
			if token != nil && token.Valid {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if userID, exists := claims["user_id"]; exists {
						identifier = userID.(string)
					}
				}
			}
		}

		mu.Lock()
		cl, exists := clients[identifier]
		if !exists {
			clients[identifier] = &client{count: 1, lastRequest: time.Now()}
			mu.Unlock()
			c.Next()
			return
		}

		if time.Since(cl.lastRequest) > window {
			cl.count = 1
			cl.lastRequest = time.Now()
			mu.Unlock()
			c.Next()
			return
		}

		cl.count++
		if cl.count > limit {
			mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
			return
		}

		mu.Unlock()
		c.Next()
	}
}