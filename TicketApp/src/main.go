package main

import (
	ticketConfig "TicketApp/src/config"
	_ "TicketApp/src/docs"
	"TicketApp/src/handler"
	"TicketApp/src/repository"
	"TicketApp/src/service"
	"TicketApp/src/type/util"
	"TicketApp/src/type/util/rabbitmq"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rabbitmq/amqp091-go"
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
	mCfg := ticketConfig.NewAppConfig()
	client, _, cancel, cfg := mCfg.ConnectDatabase()
	defer cancel()

	channel, userDeleteCheckQueue := OpenRabbitConnection(cfg)

	e := echo.New()
	e.Use(AuthorizationMiddleware)

	userClient := util.Client{BaseUrl: "http://" + cfg.UserContainerName + ":8083/api/users/"}
	categoryClient := util.Client{BaseUrl: "http://" + cfg.CategoryContainerName + ":8081/api/categories/"}

	ticketCollection := mCfg.GetCollection(client, cfg.TicketColName)
	ticketRepository := repository.NewTicketRepository(ticketCollection)
	ticketService := service.NewTicketService(ticketRepository, userClient, categoryClient, channel, userDeleteCheckQueue)
	ticketHandler := handler.NewTicketHandler(ticketService, cfg)

	go ticketService.CheckUserDeleteQueueForUpdate(channel, userDeleteCheckQueue)

	Route(e, ticketHandler)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	log.Fatal(e.Start(":8082"))
}

func OpenRabbitConnection(cfg *ticketConfig.AppConfig) (*amqp091.Channel, amqp091.Queue) {
	conStr := cfg.GetRabbitMqDialConnectionUri()

	conn := rabbitmq.Connect(conStr)
	//defer conn.Close()

	channel := rabbitmq.OpenChannel(conn)
	//defer channel.Close()

	qName := cfg.UserDeleteCheckQName
	userDeleteCheckQueue := rabbitmq.DeclareAQueue(channel, qName)

	return channel, userDeleteCheckQueue
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

func AuthorizationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Method == http.MethodPost {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return c.JSON(http.StatusUnauthorized, util.PostRequestsMustHaveATokenHeader)
			} else {
				userId := util.DecodeTokenReturnsUserId(token)
				if userId == "" {
					return c.JSON(http.StatusUnauthorized, util.TokenIsUnauthorized)

				}
			}
		}
		return next(c)
	}
}
