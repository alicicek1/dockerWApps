package userHandler

import (
	userConfig "UserApp/src/config"
	userService "UserApp/src/service"
	userEntity "UserApp/src/type/entity"
	"UserApp/src/type/util"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type UserHandler struct {
	userService userService.UserService
	cfg         *userConfig.AppConfig
}

func NewUserHandler(userService userService.UserService, cfg *userConfig.AppConfig) UserHandler {
	return UserHandler{userService: userService, cfg: cfg}
}

// UserGetById godoc
// @Summary      Show a user
// @Description  get string by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "id"
// @Success      200  {object}  entity.User
// @Failure      400  {object}  util.Error
// @Failure      404  {object}  util.Error
// @Failure      500  {object}  util.Error
// @Router       /api/users/{id} [get]
func (h *UserHandler) UserGetById(ctx echo.Context) error {
	id := ctx.Param("id")
	err := util.ValidateForUserHandlerId(id)
	if err != nil {
		return ctx.JSON(err.StatusCode, err)
	}

	user, errSrv := h.userService.UserServiceGetById(id)
	if errSrv != nil || user == nil {
		return ctx.JSON(http.StatusNotFound, util.NotFound.ModifyApplicationName("user handler").ModifyErrorCode(4010))
	}

	return ctx.JSON(http.StatusOK, user)
}

// UserUpsert godoc
// @Summary Upsert a user
// @Description Upsert a user by requested body
// @Tags users
// @Accept json
// @Produce json
// @Param id path string false "id"
// @Param userPostRequestBody body entity.UserPostRequestModel true "userPostRequestBody"
// @Success 200 {object} util.PostResponseModel
// @Failure 400 {object} util.Error
// @Failure 404 {object} util.Error
// @Failure 500 {object} util.Error
// @Router /api/users [post]
func (h *UserHandler) UserUpsert(ctx echo.Context) error {
	userPostRequestModel := userEntity.UserPostRequestModel{}

	err := json.NewDecoder(ctx.Request().Body).Decode(&userPostRequestModel)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, util.InvalidBody.ModifyApplicationName("user handler").ModifyErrorCode(4012))
	}
	user := userEntity.User{Username: userPostRequestModel.Username,
		Password: userPostRequestModel.Password,
		Email:    userPostRequestModel.Email,
		Type:     userPostRequestModel.Type,
		Age:      userPostRequestModel.Age,
	}

	id := ctx.QueryParam("id")
	if id != "" {
		if !util.IsValidUUID(id) {
			return ctx.JSON(http.StatusBadRequest, util.PathVariableIsNotValid.ModifyApplicationName("user handler").ModifyOperation("POST").ModifyErrorCode(4011))
		}
		user.Id = id
	}

	res, errSrv := h.userService.UserServiceInsert(user)
	if errSrv != nil {
		return ctx.JSON(errSrv.StatusCode, errSrv)
	}
	return ctx.JSON(http.StatusOK, res)
}

// UserDeleteById godoc
// @Summary      Delete a user
// @Description  Delete a user by id
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "id"
// @Success      200  {object}  util.DeleteResponseType
// @Failure      400  {object}  util.Error
// @Failure      404  {object}  util.Error
// @Failure      500  {object}  util.Error
// @Router       /api/users/{id} [delete]
func (h *UserHandler) UserDeleteById(ctx echo.Context) error {
	id := ctx.Param("id")
	err := util.ValidateForUserHandlerId(id)
	if err != nil {
		return ctx.JSON(err.StatusCode, err)
	}

	res, errSrv := h.userService.UserServiceDeleteById(id)
	if errSrv != nil {
		return ctx.JSON(errSrv.StatusCode, errSrv)
	}

	return ctx.JSON(http.StatusOK, res)
}

// UserGetAll godoc
// @Summary Get list of users
// @Description Get list of users
// @Tags users
// @Accept json
// @Produce json
// @Param filter query util.Filter true "filter"
// @Success 200 {object} util.GetAllResponseType
// @Failure 400 {object} util.Error
// @Failure 404 {object} util.Error
// @Failure 500 {object} util.Error
// @Router /api/users [get]
func (h *UserHandler) UserGetAll(ctx echo.Context) error {
	filter := util.Filter{}
	page, pageSize := util.ValidatePaginationFilters(ctx.QueryParam("page"), ctx.QueryParam("pageSize"), h.cfg.MaxPageLimit)
	filter.Page = page
	filter.PageSize = pageSize

	sortingField, sortingDirection := util.ValidateSortingFilters(userEntity.User{}, ctx.QueryParam("sort"), ctx.QueryParam("sDirection"))
	filter.SortingField = sortingField
	filter.SortingDirection = sortingDirection

	filters := map[string]interface{}{}
	if username := ctx.QueryParam("username"); username != "" && len(username) < 30 {
		filters["username"] = bson.M{"$regex": primitive.Regex{
			Pattern: username,
			Options: "i",
		}}
	}

	if mingAgeStr := ctx.QueryParam("minAge"); mingAgeStr != "" {
		if minAge, err := strconv.Atoi(mingAgeStr); err == nil {
			filters["age"] = bson.M{"$gte": minAge}
		}
	}

	if maxAgeStr := ctx.QueryParam("maxAge"); maxAgeStr != "" {
		if maxAge, err := strconv.Atoi(maxAgeStr); err == nil {
			minFilter, exist := filters["age"]
			if exist {
				delete(filters, "age")
				filters["$and"] = bson.A{
					bson.M{"age": minFilter},
					bson.M{"age": bson.M{"$lte": maxAge}},
				}
			} else {
				filters["age"] = bson.M{"$lte": maxAge}
			}
		}
	}

	filter.Filters = filters

	res, err := h.userService.UserServiceGetAll(filter)
	if err != nil {
		return ctx.JSON(err.StatusCode, err)
	}

	if res.Models == nil {
		return ctx.JSON(http.StatusNotFound, util.NotFound.ModifyApplicationName("user handler").ModifyErrorCode(5001))
	}
	ctx.Response().Header().Add("x-total-count", strconv.FormatInt(res.RowCount, 10))
	return ctx.JSON(http.StatusOK, res)
}

// Login godoc
// @Summary Login
// @Description Login - Besides response body token sets response header and cookie.
// @Tags users
// @Accept json
// @Produce json
// @Param loginRequestModel body entity.LoginRequestModel true "loginRequestModel"
// @Success 200 {object} entity.LoginResponseModel
// @Failure 400 {object} util.Error
// @Failure 404 {object} util.Error
// @Failure 500 {object} util.Error
// @Router /api/users/login [post]
func (h *UserHandler) Login(ctx echo.Context) error {
	loginRequestModel := userEntity.LoginRequestModel{}

	err := json.NewDecoder(ctx.Request().Body).Decode(&loginRequestModel)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, util.InvalidBody.ModifyApplicationName("user handler").ModifyErrorCode(4054).ModifyOperation("LOGIN"))
	}

	if loginRequestModel.Username == "" || *loginRequestModel.Password == "" {
		return ctx.JSON(http.StatusBadRequest, util.InvalidBody.ModifyApplicationName("user handler").ModifyErrorCode(4055).ModifyOperation("LOGIN"))
	}

	result, errSrv := h.userService.UserServiceLogin(loginRequestModel)
	if errSrv != nil {
		return ctx.JSON(errSrv.StatusCode, errSrv)
	}

	ctx.Response().Header().Add("Token", result.Token)
	ctx.SetCookie(&http.Cookie{
		Name:    reflect.TypeOf(result.Token).Name(),
		Value:   result.Token,
		Expires: result.ExpiresDate,
	})
	return ctx.JSON(http.StatusCreated, result)
}

// UserIfExistById godoc
// @Summary UserIfExistById
// @Description UserIfExistById - Validation endpoint for ticket post.
// @Tags users
// @Accept json
// @Produce json
// @Param loginRequestModel body entity.LoginRequestModel true "loginRequestModel"
// @Success 200 {object} bool
// @Failure 400 {object} util.Error
// @Failure 404 {object} util.Error
// @Failure 500 {object} util.Error
// @Router /api/users/isExist/{id} [get]
func (h *UserHandler) UserIfExistById(ctx echo.Context) error {
	id := ctx.Param("id")
	err := util.ValidateForUserHandlerId(id)
	if err != nil {
		return ctx.JSON(err.StatusCode, err)
	}

	res, errSrv := h.userService.UserIfExistById(id)
	if errSrv != nil {
		return ctx.JSON(errSrv.StatusCode, errSrv)
	}

	return ctx.JSON(http.StatusOK, res)
}

// ReadCsv godoc
// @Summary ReadCsv
// @Description ReadCsv
// @Tags users
// @Accept  json
// @Produce  json
// @Param fileCsv formData file true "Body with file csv"
// @Success 200 {object} bool
// @Failure 400 {object} util.Error
// @Failure 404 {object} util.Error
// @Failure 500 {object} util.Error
// @Router /api/users/readCsv [post]
func (h *UserHandler) ReadCsv(ctx echo.Context) error {
	fileHeader, fileHeaderErr := ctx.FormFile("file")
	CheckError(fileHeaderErr)

	if !strings.Contains(fileHeader.Filename, "csv") {
		return ctx.JSON(http.StatusBadRequest, util.NewError("user handler", "reading csv file", "file must be cvs", http.StatusBadRequest, 5000))
	}

	file, fileOpenErr := fileHeader.Open()
	defer file.Close()
	CheckError(fileOpenErr)

	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()
	CheckError(err)
	fmt.Println(data)

	//var wg sync.WaitGroup
	//go func() {
	//	data, err := csvReader.ReadAll()
	//	CheckError(err)
	//	for _, datum := range data {
	//		wg.Add(1)
	//		fmt.Println(datum)
	//		wg.Done()
	//	}
	//}()
	//wg.Wait()
	return ctx.JSON(http.StatusOK, "File reading. Check console.")
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
