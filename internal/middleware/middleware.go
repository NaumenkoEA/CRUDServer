// Package middleware : file contains verification for authenticated user
package middleware

import (
	"awesomeProject/internal/service"

	"github.com/labstack/echo/v4/middleware"
)

// IsAuthenticated check for authenticated
var IsAuthenticated = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: service.JwtKey,
})
