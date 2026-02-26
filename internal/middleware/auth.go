package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/akatakan/nobetgo/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates the JWT token in the Authorization header
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			util.JSONError(c, http.StatusUnauthorized, "Authorization header is required", nil)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			util.JSONError(c, http.StatusUnauthorized, "Invalid authorization header format", nil)
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing algorithm to prevent "alg:none" bypass
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			util.JSONError(c, http.StatusUnauthorized, "Invalid or expired token", err)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			util.JSONError(c, http.StatusUnauthorized, "Invalid token claims", nil)
			return
		}

		// Safe type assertions to prevent panics from malformed tokens
		userIDVal, ok1 := claims["user_id"].(float64)
		roleVal, ok2 := claims["role"].(string)
		if !ok1 || !ok2 {
			util.JSONError(c, http.StatusUnauthorized, "Invalid token claims", nil)
			return
		}

		// Set user info to context (key "userID" matches handler expectations)
		c.Set("userID", uint(userIDVal))
		c.Set("role", roleVal)

		c.Next()
	}
}

// RoleMiddleware restricts access to specific roles
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			util.JSONError(c, http.StatusUnauthorized, "Role not found in context", nil)
			return
		}

		userRole := role.(string)
		isAllowed := false
		for _, r := range allowedRoles {
			if userRole == r {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			util.JSONError(c, http.StatusForbidden, "Forbidden: You don't have permission to access this resource", nil)
			return
		}

		c.Next()
	}
}
