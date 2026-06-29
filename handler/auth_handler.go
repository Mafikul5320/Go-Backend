package handler

import (
	"spotsync/dto"
	"spotsync/service"
	"spotsync/utils"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService service.AuthService
	validator   *utils.CustomValidator
}

func NewAuthHandler(e *echo.Echo, authService service.AuthService, v *utils.CustomValidator) {
	handler := &AuthHandler{
		authService: authService,
		validator:   v,
	}

	g := e.Group("/api/v1/auth")
	g.POST("/register", handler.Register)
	g.POST("/login", handler.Login)
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.SendError(c, 400, "Bad Request", "Invalid JSON format")
	}

	if err := h.validator.Validate(&req); err != nil {
		return utils.SendError(c, 400, "Bad Request", err.Error())
	}

	user, err := h.authService.Register(req)
	if err != nil {
		if err.Error() == "email already in use" {
			return utils.SendError(c, 400, "Bad Request", err.Error())
		}
		return utils.SendError(c, 500, "Internal Server Error", err.Error())
	}

	return utils.SendSuccess(c, 201, "User registered successfully", user)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.SendError(c, 400, "Bad Request", "Invalid JSON format")
	}

	if err := h.validator.Validate(&req); err != nil {
		return utils.SendError(c, 400, "Bad Request", err.Error())
	}

	token, user, err := h.authService.Login(req)
	if err != nil {
		return utils.SendError(c, 401, "Unauthorized", err.Error())
	}

	response := map[string]interface{}{
		"token": token,
		"user":  user,
	}

	return utils.SendSuccess(c, 200, "Login successful", response)
}
