package handler

import (
	"strconv"

	"spotsync/dto"
	"spotsync/middleware"
	"spotsync/service"
	"spotsync/utils"

	"github.com/labstack/echo/v4"
)

type ReservationHandler struct {
	reservationService service.ReservationService
	validator          *utils.CustomValidator
}

func NewReservationHandler(e *echo.Echo, reservationService service.ReservationService, v *utils.CustomValidator) {
	handler := &ReservationHandler{
		reservationService: reservationService,
		validator:          v,
	}

	g := e.Group("/api/v1/reservations")
	g.Use(middleware.JWTMiddleware())

	// Authenticated routes
	g.POST("", handler.CreateReservation)
	g.GET("/my-reservations", handler.GetMyReservations)
	g.DELETE("/:id", handler.CancelReservation)

	// Admin routes
	adminGroup := g.Group("")
	adminGroup.Use(middleware.AdminMiddleware())
	adminGroup.GET("", handler.GetAllReservations)
}

func (h *ReservationHandler) CreateReservation(c echo.Context) error {
	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return utils.SendError(c, 400, "Bad Request", "Invalid JSON format")
	}

	if err := h.validator.Validate(&req); err != nil {
		return utils.SendError(c, 400, "Bad Request", err.Error())
	}

	userID := c.Get("user_id").(uint)

	reservation, err := h.reservationService.CreateReservation(userID, req)
	if err != nil {
		if err.Error() == "zone is at full capacity" {
			return utils.SendError(c, 409, "Conflict", err.Error())
		}
		return utils.SendError(c, 500, "Internal Server Error", err.Error())
	}

	return utils.SendSuccess(c, 201, "Reservation confirmed successfully", reservation)
}

func (h *ReservationHandler) GetMyReservations(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	reservations, err := h.reservationService.GetMyReservations(userID)
	if err != nil {
		return utils.SendError(c, 500, "Internal Server Error", err.Error())
	}

	// Ensure we return an empty array instead of null
	if reservations == nil {
		reservations = []dto.ReservationResponse{}
	}

	return utils.SendSuccess(c, 200, "My reservations retrieved successfully", reservations)
}

func (h *ReservationHandler) CancelReservation(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.SendError(c, 400, "Bad Request", "Invalid reservation ID")
	}

	userID := c.Get("user_id").(uint)

	err = h.reservationService.CancelReservation(uint(id), userID)
	if err != nil {
		if err.Error() == "reservation not found or you don't have permission" {
			return utils.SendError(c, 403, "Forbidden", err.Error())
		}
		if err.Error() == "reservation is already cancelled" {
			return utils.SendError(c, 400, "Bad Request", err.Error())
		}
		return utils.SendError(c, 500, "Internal Server Error", err.Error())
	}

	return utils.SendSuccess(c, 200, "Reservation cancelled successfully", nil)
}

func (h *ReservationHandler) GetAllReservations(c echo.Context) error {
	reservations, err := h.reservationService.GetAllReservations()
	if err != nil {
		return utils.SendError(c, 500, "Internal Server Error", err.Error())
	}

	if reservations == nil {
		reservations = []dto.ReservationResponse{}
	}

	return utils.SendSuccess(c, 200, "All reservations retrieved successfully", reservations)
}
