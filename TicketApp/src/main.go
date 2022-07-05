package main

import (
	ticketConfig "TicketApp/src/config"
	_ "TicketApp/src/docs"
	"TicketApp/src/handler"
	"TicketApp/src/repository"
	"TicketApp/src/service"
	"TicketApp/src/type/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
	"net/http"
)

// @title Swagger Ticket API
// @version 1.0
// @description This is a sample ticket & category API server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host user.swagger.io
// @BasePath /api/tickets
func main() {
	mCfg := ticketConfig.NewMongoConfig()
	client, _, cancel, cfg := mCfg.ConnectDatabase()
	defer cancel()

	e := echo.New()

	userClient := util.Client{BaseUrl: "http://user_service:8083/api/users/"}
	categoryClient := util.Client{BaseUrl: "http://category_service:8081/api/categories/"}

	ticketCollection := mCfg.GetCollection(client, cfg.TicketColName)
	ticketRepository := repository.NewTicketRepository(ticketCollection)
	ticketService := service.NewTicketService(ticketRepository, userClient, categoryClient)
	ticketHandler := handler.NewTicketHandler(ticketService, cfg)

	Route(e, ticketHandler)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	log.Fatal(e.Start(":8082"))
}

func Route(e *echo.Echo, ticketHandler handler.TicketHandler) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut},
	}))

	ticketGroup := e.Group("/api/tickets")
	ticketGroup.GET("/:id", ticketHandler.TicketGetById)
	ticketGroup.GET("", ticketHandler.TicketGetAll)
	ticketGroup.POST("", ticketHandler.TicketInsert)
	ticketGroup.DELETE("/:id", ticketHandler.TicketDeleteById)
}
