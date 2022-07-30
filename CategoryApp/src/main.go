package main

import (
	categoryConfig "CategoryApp/src/config"
	_ "CategoryApp/src/docs"
	"CategoryApp/src/handler"
	categoryRepository "CategoryApp/src/repository"
	categoryService "CategoryApp/src/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
	"net/http"
)

// @title Swagger Category API
// @version 1.0
// @description This is a sample ticket & category API server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host user.swagger.io
// @BasePath /api/categories
func main() {
	mCfg := categoryConfig.NewMongoConfig()
	client, _, cancel, cfg := mCfg.ConnectDatabase()
	defer cancel()

	e := echo.New()

	categoryCollection := mCfg.GetCollection(client, cfg.CategoryColName)
	categoryRepository := categoryRepository.NewCategoryRepository(categoryCollection)
	categoryService := categoryService.NewCategoryService(categoryRepository)
	categoryHandler := handler.NewCategoryHandler(categoryService, cfg)

	healthCheckHandler := handler.NewHealthCheckHandler()
	e.GET("/healthCheck", healthCheckHandler.HealthCheck)

	Route(e, categoryHandler)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	log.Fatal(e.Start(":8081"))

}

func Route(e *echo.Echo, categoryHandler handler.CategoryHandler) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut},
	}))

	categoryGroup := e.Group("/api/categories")
	categoryGroup.GET("/:id", categoryHandler.CategoryGetById)
	categoryGroup.GET("", categoryHandler.CategoryGetAll)
	categoryGroup.POST("", categoryHandler.CategoryInsert)
	categoryGroup.DELETE("/:id", categoryHandler.CategoryDeleteById)
	categoryGroup.GET("/isExist/:id", categoryHandler.CategoryIfExistById)

}
