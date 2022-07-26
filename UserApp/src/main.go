package main

import (
	userConfig "UserApp/src/config"
	_ "UserApp/src/docs"
	userHandler "UserApp/src/handler"
	userRepository "UserApp/src/repository"
	userService "UserApp/src/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
	"net/http"
)

// @title Swagger User API
// @version 1.0
// @description This is a sample user API server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host user.swagger.io
// @BasePath /api/users
func main() {
	mCfg := userConfig.NewMongoConfig()
	client, _, cancel, cfg := mCfg.ConnectDatabase()
	defer cancel()

	e := echo.New()
	//e.HTTPErrorHandler = util.NewHttpErrorHandler(util.NewErrorStatusCodeMaps()).Handler

	userCollection := mCfg.GetCollection(client, cfg.UserColName)
	userRepository := userRepository.NewUserRepository(userCollection)
	userService := userService.NewUserService(userRepository)
	userHandler := userHandler.NewUserHandler(userService, cfg)

	Route(e, userHandler)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	log.Fatal(e.Start(":8083"))

}

func Route(e *echo.Echo, userHandler userHandler.UserHandler) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut},
	}))

	userGroup := e.Group("/api/users")
	userGroup.GET("/:id", userHandler.UserGetById)
	userGroup.GET("", userHandler.UserGetAll)
	userGroup.POST("", userHandler.UserUpsert)
	userGroup.POST("/login", userHandler.Login)
	userGroup.DELETE("/:id", userHandler.UserDeleteById)
	userGroup.GET("/isExist/:id", userHandler.UserIfExistById)
	userGroup.POST("/readCsv", userHandler.ReadCsv)
}
