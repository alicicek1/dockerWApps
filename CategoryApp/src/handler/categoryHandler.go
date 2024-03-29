package handler

import (
	categoryConfig "CategoryApp/src/config"
	categoryService "CategoryApp/src/service"
	categoryType "CategoryApp/src/type"
	"CategoryApp/src/type/util"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CategoryHandler struct {
	categoryService categoryService.CategoryService
	cfg             *categoryConfig.AppConfig
}

func NewCategoryHandler(categoryService categoryService.CategoryService, cfg *categoryConfig.AppConfig) CategoryHandler {
	return CategoryHandler{
		categoryService: categoryService,
		cfg:             cfg,
	}
}

// CategoryGetById godoc
// @Summary      Show a category
// @Description  get string by ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "id"
// @Success      200  {object}  categoryType.Category
// @Failure      400  {object}  util.Error
// @Failure      404  {object}  util.Error
// @Failure      500  {object}  util.Error
// @Router       /api/categories/{id} [get]
func (h *CategoryHandler) CategoryGetById(ctx echo.Context) error {
	id := ctx.Param("id")
	err := util.ValidateForUserHandlerId(id)
	if err != nil {
		return ctx.JSON(err.StatusCode, err)
	}

	category, errSrv := h.categoryService.CategoryServiceGetById(id)
	if errSrv != nil || category == nil {
		return ctx.JSON(http.StatusNotFound, errSrv)
	}

	return ctx.JSON(http.StatusOK, category)
}

// CategoryInsert godoc
// @Summary Insert a category
// @Description Insert a category by requested body
// @Tags categories
// @Accept json
// @Produce json
// @Param categoryPostRequestModel body categoryType.CategoryPostRequestModel true "categoryPostRequestModel"
// @Param Authorization header string true "Authorization"
// @Success 200 {object} util.PostResponseModel
// @Failure 400 {object} util.Error
// @Failure 404 {object} util.Error
// @Failure 500 {object} util.Error
// @Router /api/categories [post]
func (h *CategoryHandler) CategoryInsert(ctx echo.Context) error {
	categoryPostRequestModel := categoryType.CategoryPostRequestModel{}

	err := json.NewDecoder(ctx.Request().Body).Decode(&categoryPostRequestModel)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, util.InvalidBody.ModifyApplicationName("category handler").ModifyErrorCode(4022))
	}

	res, errSrv := h.categoryService.CategoryServiceInsert(categoryPostRequestModel)
	if errSrv != nil {
		return ctx.JSON(errSrv.StatusCode, errSrv)
	}
	return ctx.JSON(http.StatusOK, res)
}

// CategoryDeleteById godoc
// @Summary      Delete a category
// @Description  Delete a category by id
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "id"
// @Success      200  {object}  util.DeleteResponseType
// @Failure      400  {object}  util.Error
// @Failure      404  {object}  util.Error
// @Failure      500  {object}  util.Error
// @Router       /api/categories/{id} [delete]
func (h *CategoryHandler) CategoryDeleteById(ctx echo.Context) error {
	id := ctx.Param("id")
	err := util.ValidateForUserHandlerId(id)
	if err != nil {
		return ctx.JSON(err.StatusCode, err)
	}

	res, errSrv := h.categoryService.CategoryServiceDeleteById(id)
	if errSrv != nil {
		return ctx.JSON(errSrv.StatusCode, errSrv)
	}

	return ctx.JSON(http.StatusOK, res)
}

// CategoryGetAll godoc
// @Summary Get list of categories
// @Description Get list of categories
// @Tags categories
// @Accept json
// @Produce json
// @Param filter query util.Filter true "filter"
// @Param name query string false "name"
// @Success 200 {object} util.GetAllResponseType
// @Failure 400 {object} util.Error
// @Failure 404 {object} util.Error
// @Failure 500 {object} util.Error
// @Router /api/categories [get]
func (h *CategoryHandler) CategoryGetAll(ctx echo.Context) error {
	filter := util.Filter{}
	page, pageSize := util.ValidatePaginationFilters(ctx.QueryParam("page"), ctx.QueryParam("pageSize"), h.cfg.MaxPageLimit)
	filter.Page = page
	filter.PageSize = pageSize

	sortingField, sortingDirection := util.ValidateSortingFilters(categoryType.Category{}, ctx.QueryParam("sort"), ctx.QueryParam("sDirection"))
	filter.SortingField = sortingField
	filter.SortingDirection = sortingDirection

	nameValue := ctx.QueryParam("name")
	filter.Filters = map[string]interface{}{}
	if nameValue != "" {
		filter.Filters["name"] = util.CreateEqualFilter(nameValue, "name")
	}

	categories, errSrv := h.categoryService.CategoryServiceGetAll(filter)
	if errSrv != nil {
		return ctx.JSON(errSrv.StatusCode, errSrv)
	}

	if categories == nil || categories.RowCount == 0 {
		return ctx.JSON(http.StatusNotFound, util.NotFound.ModifyApplicationName("category handler").ModifyErrorCode(5000))
	}
	return ctx.JSON(http.StatusOK, categories)
}

// CategoryIfExistById godoc
// @Summary      Checks if category exist by id
// @Description  Checks if category exist by id
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "id"
// @Success      200  {object}  bool
// @Failure      400  {object}  util.Error
// @Failure      404  {object}  util.Error
// @Failure      500  {object}  util.Error
// @Router       /isExist/{id} [get]
func (h *CategoryHandler) CategoryIfExistById(ctx echo.Context) error {
	id := ctx.Param("id")
	err := util.ValidateForUserHandlerId(id)
	if err != nil {
		return ctx.JSON(err.StatusCode, err)
	}

	res, errSrv := h.categoryService.CategoryIfExistById(id)
	if errSrv != nil {
		return ctx.JSON(errSrv.StatusCode, errSrv)
	}

	return ctx.JSON(http.StatusOK, res)
}
