package handler

import (
	ticketConfig "TicketApp/src/config"
	"TicketApp/src/service"
	"TicketApp/src/type/entity"
	"TicketApp/src/type/util"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
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
func (t *TicketHandler) TicketGetById(ctx echo.Context) error {
	id := ctx.Param("id")
	err := util.ValidateForUserHandlerId(id)
	if err != nil {
		return ctx.JSON(err.StatusCode, err)
	}

	ticket, errSrv := t.ticketService.TicketServiceGetById(id)
	if errSrv != nil || ticket == nil {
		return ctx.JSON(http.StatusNotFound, errSrv)
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
// @Param Authorization header string true "Authorization"
// @Success 200 {object} util.PostResponseModel
// @Failure 400 {object} util.Error
// @Failure 404 {object} util.Error
// @Failure 500 {object} util.Error
// @Router /api/tickets [post]
func (t *TicketHandler) TicketInsert(ctx echo.Context) error {
	ticketPostRequestModel := entity.TicketPostRequestModel{}

	err := json.NewDecoder(ctx.Request().Body).Decode(&ticketPostRequestModel)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, util.InvalidBody.ModifyApplicationName("category handler").ModifyErrorCode(4022).ModifyOperation("POST"))
	}

	authorization := ctx.Request().Header["Authorization"][0]

	res, errSrv := t.ticketService.TicketServiceInsert(ticketPostRequestModel, authorization)
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
func (t *TicketHandler) TicketDeleteById(ctx echo.Context) error {
	id := ctx.Param("id")
	err := util.ValidateForUserHandlerId(id)
	if err != nil {
		return ctx.JSON(err.StatusCode, err)
	}

	res, errSrv := t.ticketService.TicketServiceDeleteById(id)
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
// @Param categoryId query string false "categoryId"
// @Success 200 {object} util.GetAllResponseType
// @Failure 400 {object} util.Error
// @Failure 404 {object} util.Error
// @Failure 500 {object} util.Error
// @Router /api/tickets [get]
func (t *TicketHandler) TicketGetAll(ctx echo.Context) error {
	filter := util.Filter{}
	page, pageSize := util.ValidatePaginationFilters(ctx.QueryParam("page"), ctx.QueryParam("pageSize"), t.cfg.MaxPageLimit)
	filter.Page = page
	filter.PageSize = pageSize

	sortingField, sortingDirection := util.ValidateSortingFilters(entity.Ticket{}, ctx.QueryParam("sort"), ctx.QueryParam("sDirection"))
	filter.SortingField = sortingField
	filter.SortingDirection = sortingDirection

	//filtering
	//CategoryId
	filters := map[string]interface{}{}
	if categoryId := ctx.QueryParam("categoryId"); categoryId != "" {
		filters["categoryId"] = bson.M{"$regex": primitive.Regex{
			Pattern: categoryId,
			Options: "i",
		}}
	}

	filter.Filters = filters

	tickets, err := t.ticketService.TicketServiceGetAll(filter)
	if err != nil {
		return ctx.JSON(err.StatusCode, err)
	}

	if tickets == nil || tickets.RowCount == 0 {
		return ctx.JSON(http.StatusNotFound, util.NotFound.ModifyApplicationName("user handler").ModifyErrorCode(5001))
	}
	return ctx.JSON(http.StatusOK, tickets)
}
