package handler

import (
    "strconv"

    "spotsync/dto"
    "spotsync/middleware"
    "spotsync/service"
    "spotsync/utils"

    "github.com/labstack/echo/v4"
)

type UserHandler struct {
    userService service.UserService
    validator   *utils.CustomValidator
}

func NewUserHandler(e *echo.Echo, userService service.UserService, v *utils.CustomValidator) {
    handler := &UserHandler{
        userService: userService,
        validator:   v,
    }

    g := e.Group("/api/v1/users")
    g.Use(middleware.JWTMiddleware())
    g.Use(middleware.AdminMiddleware())

    g.GET("", handler.GetAllUsers)
    g.GET("/:id", handler.GetUserByID)
    g.DELETE("/:id", handler.DeleteUser)
}

func (h *UserHandler) GetAllUsers(c echo.Context) error {
    users, err := h.userService.GetAllUsers()
    if err != nil {
        return utils.SendError(c, 500, "Internal Server Error", err.Error())
    }

    if users == nil {
        users = []dto.UserResponse{}
    }

    return utils.SendSuccess(c, 200, "Users retrieved successfully", users)
}

func (h *UserHandler) GetUserByID(c echo.Context) error {
    idParam := c.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        return utils.SendError(c, 400, "Bad Request", "Invalid user ID")
    }

    user, err := h.userService.GetUserByID(uint(id))
    if err != nil {
        if err.Error() == "record not found" {
            return utils.SendError(c, 404, "Not Found", "User not found")
        }
        return utils.SendError(c, 500, "Internal Server Error", err.Error())
    }

    return utils.SendSuccess(c, 200, "User retrieved successfully", user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
    idParam := c.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        return utils.SendError(c, 400, "Bad Request", "Invalid user ID")
    }

    err = h.userService.DeleteUser(uint(id))
    if err != nil {
        return utils.SendError(c, 500, "Internal Server Error", err.Error())
    }

    return utils.SendSuccess(c, 200, "User deleted successfully", nil)
}
