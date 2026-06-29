package middleware

import (
	"os"
	"strings"

	"spotsync/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return utils.SendError(c, 401, "Unauthorized", "Missing token")
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return utils.SendError(c, 401, "Unauthorized", "Invalid token format")
			}

			tokenString := parts[1]
			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				secret = "supersecretkey"
			}

			token, err := jwt.ParseWithClaims(tokenString, &utils.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				return utils.SendError(c, 401, "Unauthorized", "Invalid or expired token")
			}

			claims, ok := token.Claims.(*utils.JwtCustomClaims)
			if !ok {
				return utils.SendError(c, 401, "Unauthorized", "Invalid token claims")
			}

			c.Set("user_id", claims.ID)
			c.Set("user_role", claims.Role)

			return next(c)
		}
	}
}

func AdminMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role := c.Get("user_role")
			if role != "admin" {
				return utils.SendError(c, 403, "Forbidden", "Admin access required")
			}
			return next(c)
		}
	}
}
