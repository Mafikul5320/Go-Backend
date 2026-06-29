package handler

import (
	"strconv"

	"spotsync/dto"
	"spotsync/middleware"
	"spotsync/service"
	"spotsync/utils"

	"github.com/labstack/echo/v4"
)

type ZoneHandler struct {
	zoneService service.ZoneService
	validator   *utils.CustomValidator
}

func NewZoneHandler(e *echo.Echo, zoneService service.ZoneService, v *utils.CustomValidator) {
	handler := &ZoneHandler{
		zoneService: zoneService,
		validator:   v,
	}

	g := e.Group("/api/v1/zones")
	
	// Public routes
	g.GET("", handler.GetAllZones)
	g.GET("/:id", handler.GetZoneByID)

	// Admin routes
	adminGroup := g.Group("")
	adminGroup.Use(middleware.JWTMiddleware())
	adminGroup.Use(middleware.AdminMiddleware())
	adminGroup.POST("", handler.CreateZone)
}

func (h *ZoneHandler) CreateZone(c echo.Context) error {
	var req dto.CreateZoneRequest
	if err := c.Bind(&req); err != nil {
		return utils.SendError(c, 400, "Bad Request", "Invalid JSON format")
	}

	if err := h.validator.Validate(&req); err != nil {
		return utils.SendError(c, 400, "Bad Request", err.Error())
	}

	zone, err := h.zoneService.CreateZone(req)
	if err != nil {
		return utils.SendError(c, 500, "Internal Server Error", err.Error())
	}

	return utils.SendSuccess(c, 201, "Parking zone created successfully", zone)
}

func (h *ZoneHandler) GetAllZones(c echo.Context) error {
	zones, err := h.zoneService.GetAllZones()
	if err != nil {
		return utils.SendError(c, 500, "Internal Server Error", err.Error())
	}

	// Ensure we return an empty array instead of null if empty
	if zones == nil {
		zones = []dto.ZoneResponse{}
	}

	return utils.SendSuccess(c, 200, "Parking zones retrieved successfully", zones)
}

func (h *ZoneHandler) GetZoneByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.SendError(c, 400, "Bad Request", "Invalid zone ID")
	}

	zone, err := h.zoneService.GetZoneByID(uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			return utils.SendError(c, 404, "Not Found", "Parking zone not found")
		}
		return utils.SendError(c, 500, "Internal Server Error", err.Error())
	}

	return utils.SendSuccess(c, 200, "Parking zone retrieved successfully", zone)
}
