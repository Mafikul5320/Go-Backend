package utils

import "github.com/labstack/echo/v4"

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

func SendSuccess(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SendError(c echo.Context, statusCode int, message string, errors interface{}) error {
	return c.JSON(statusCode, ErrorResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}
