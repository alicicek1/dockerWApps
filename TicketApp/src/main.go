package main

import (
	ticketConfig "TicketApp/src/config"
	_ "TicketApp/src/docs"
	"TicketApp/src/handler"
	"TicketApp/src/repository"
	"TicketApp/src/service"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
)

// @title Swagger Ticket & Category API
// @version 1.0
// @description This is a sample ticket & category API server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host user.swagger.io
// @BasePath /api/
func main() {
	mCfg := ticketConfig.NewMongoConfig()
	client, _, cancel, cfg := mCfg.ConnectDatabase()
	defer cancel()

	e := echo.New()
	//e.HTTPErrorHandler = util.NewHttpErrorHandler(util.NewErrorStatusCodeMaps()).Handler

	categoryCollection := mCfg.GetCollection(client, cfg.CategoryColName)
	categoryRepository := repository.NewCategoryRepository(categoryCollection)
	categoryService := service.NewCategoryService(categoryRepository)
	categoryHandler := handler.NewCategoryHandler(categoryService, cfg)
	categoryGroup := e.Group("/api/categories")
	categoryGroup.GET("/:id", categoryHandler.CategoryGetById)
	categoryGroup.GET("", categoryHandler.CategoryGetAll)
	categoryGroup.POST("", categoryHandler.CategoryInsert)
	categoryGroup.DELETE("/:id", categoryHandler.CategoryDeleteById)
	//categoryGroup.GET("/swagger/*", echoSwagger.WrapHandler)

	ticketCollection := mCfg.GetCollection(client, cfg.TicketColName)
	ticketRepository := repository.NewTicketRepository(ticketCollection)
	ticketService := service.NewTicketService(ticketRepository)
	ticketHandler := handler.NewTicketHandler(ticketService, cfg)
	ticketGroup := e.Group("/api/tickets")
	ticketGroup.GET("/:id", ticketHandler.TicketGetById)
	ticketGroup.GET("", ticketHandler.TicketGetAll)
	ticketGroup.POST("", ticketHandler.TicketInsert)
	ticketGroup.DELETE("/:id", ticketHandler.TicketDeleteById)
	//ticketGroup.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	log.Fatal(e.Start(":8084"))
}