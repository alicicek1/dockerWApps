package handler

import (
	ticketConfig "TicketApp/src/config"
	"TicketApp/src/service"
	"TicketApp/src/type/entity"
	"TicketApp/src/type/util"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type TicketHandler struct {
	ticketService service.TicketService
	cfg           *ticketConfig.AppConfig
}

func NewTicketHandler(ticketService service.TicketService, cfg *ticketConfig.AppConfig) TicketHandler {
	return TicketHandler{ticketService: ticketService, cfg: cfg}
}

// TicketGetById godoc
// @Summary      Show a ticket
// @Description  get string by ID
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "id"
// @Success      200  {object}  entity.Ticket
// @Failure      400  {object}  util.Error
// @Failure      404  {object}  util.Error
// @Failure      500  {object}  util.Error
// @Router       /api/tickets/{id} [get]
func (h *TicketHandler) TicketGetById(ctx echo.Context) error {
	id := ctx.Param("id")

	if id == "" {
		return ctx.JSON(http.StatusBadRequest, util.PathVariableNotFound.ModifyApplicationName("ticket handler").ModifyErrorCode(4018))
	}

	if !util.IsValidUUID(id) {
		return ctx.JSON(http.StatusBadRequest, util.PathVariableIsNotValid.ModifyApplicationName("ticket handler").ModifyErrorCode(4019))
	}

	ticket, errSrv := h.ticketService.TicketServiceGetById(id)
	if errSrv != nil || ticket == nil {
		return ctx.JSON(http.StatusNotFound, util.NotFound.ModifyApplicationName("category handler").ModifyErrorCode(4018))
	}

	return ctx.JSON(http.StatusOK, ticket)
}

// TicketInsert godoc
// @Summary Insert a ticket
// @Description Insert a ticket by requested body
// @Tags tickets
// @Accept json
// @Produce json
// @Param ticketPostRequestModel body entity.TicketPostRequestModel true "ticketPostRequestModel"
// @Success 200 {object} util.PostResponseModel
// @Failure 400 {object} util.Error
// @Failure 404 {object} util.Error
// @Failure 500 {object} util.Error
// @Router /api/tickets [post]
func (h *TicketHandler) TicketInsert(ctx echo.Context) error {
	ticketPostRequestModel := entity.TicketPostRequestModel{}

	err := json.NewDecoder(ctx.Request().Body).Decode(&ticketPostRequestModel)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, util.InvalidBody.ModifyApplicationName("category handler").ModifyErrorCode(4022).ModifyOperation("POST"))
	}
	category := entity.Ticket{
		CategoryId:     ticketPostRequestModel.CategoryId,
		Attachments:    ticketPostRequestModel.Attachments,
		Answers:        ticketPostRequestModel.Answers,
		Subject:        ticketPostRequestModel.Subject,
		Body:           ticketPostRequestModel.Body,
		CreatedBy:      ticketPostRequestModel.CreatedBy,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		LastAnsweredAt: time.Now(),
		Status:         byte(entity.CREATED),
	}

	res, errSrv := h.ticketService.TicketServiceInsert(category)
	if errSrv != nil {
		return ctx.JSON(errSrv.StatusCode, errSrv)
	}
	return ctx.JSON(http.StatusOK, res)
}

// TicketDeleteById godoc
// @Summary      Delete a ticket
// @Description  Delete a ticket by id
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "id"
// @Success      200  {object}  util.DeleteResponseType
// @Failure      400  {object}  util.Error
// @Failure      404  {object}  util.Error
// @Failure      500  {object}  util.Error
// @Router       /api/tickets/{id} [delete]
func (h *TicketHandler) TicketDeleteById(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return ctx.JSON(http.StatusBadRequest, util.PathVariableNotFound.ModifyApplicationName("user handler").ModifyErrorCode(4020))
	}

	if !util.IsValidUUID(id) {
		return ctx.JSON(http.StatusBadRequest, util.PathVariableIsNotValid.ModifyApplicationName("user handler").ModifyErrorCode(4021))
	}

	res, errSrv := h.ticketService.TicketServiceDeleteById(id)
	if errSrv != nil {
		return ctx.JSON(errSrv.StatusCode, errSrv)
	}

	return ctx.JSON(http.StatusOK, res)
}

// TicketGetAll godoc
// @Summary Get list of tickets
// @Description Get list of tickets
// @Tags tickets
// @Accept json
// @Produce json
// @Param filter query util.Filter true "filter"
// @Success 200 {object} util.GetAllResponseType
// @Failure 400 {object} util.Error
// @Failure 404 {object} util.Error
// @Failure 500 {object} util.Error
// @Router /api/tickets [get]
func (h *TicketHandler) TicketGetAll(ctx echo.Context) error {
	filter := util.Filter{}
	page, pageSize := util.ValidatePaginationFilters(ctx.QueryParam("page"), ctx.QueryParam("pageSize"), h.cfg.MaxPageLimit)
	filter.Page = page
	filter.PageSize = pageSize

	sortingField, sortingDirection := util.ValidateSortingFilters(entity.Ticket{}, ctx.QueryParam("sort"), ctx.QueryParam("sDirection"))
	filter.SortingField = sortingField
	filter.SortingDirection = sortingDirection

	//filtering

	tickets, err := h.ticketService.TicketServiceGetAll(filter)
	if err != nil {
		return ctx.JSON(err.StatusCode, err)
	}

	if tickets == nil {
		return ctx.JSON(http.StatusNotFound, util.NotFound.ModifyApplicationName("user handler").ModifyErrorCode(5001))
	}
	return ctx.JSON(http.StatusOK, tickets)
}
