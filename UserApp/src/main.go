package main

import (
	userConfig "UserApp/src/config"
	_ "UserApp/src/docs"
	userHandler "UserApp/src/handler"
	userRepository "UserApp/src/repository"
	userService "UserApp/src/service"
	_ "UserApp/src/type/util/client"
	client2 "UserApp/src/type/util/client"
	"UserApp/src/type/util/rabbitmq"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rabbitmq/amqp091-go"
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
	mCfg := userConfig.NewAppConfig()
	client, _, cancel, cfg := mCfg.ConnectDatabase()
	defer cancel()

	ticketClient := client2.Client{BaseUrl: "http://ticket_service:8082/api/tickets/"}
	channel, userDeleteCheckQueue := OpenRabbitConnection(cfg)

	e := echo.New()
	//e.HTTPErrorHandler = util.NewHttpErrorHandler(util.NewErrorStatusCodeMaps()).Handler

	userCollection := mCfg.GetCollection(client, cfg.UserColName)
	userRepository := userRepository.NewUserRepository(userCollection)
	userService := userService.NewUserService(userRepository, channel, userDeleteCheckQueue, ticketClient)
	userHandler := userHandler.NewUserHandler(userService, cfg)

	Route(e, userHandler)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	log.Fatal(e.Start(":8083"))

}

func OpenRabbitConnection(cfg *userConfig.AppConfig) (*amqp091.Channel, amqp091.Queue) {
	conStr := userConfig.GetRabbitMqDialConnectionUri(cfg)
	fmt.Println("conStr->", conStr)
	conn := rabbitmq.Connect(conStr)
	//defer conn.Close()

	channel := rabbitmq.OpenChannel(conn)
	//defer channel.Close()

	qName := cfg.UserDeleteCheckQName
	userDeleteCheckQueue := rabbitmq.DeclareAQueue(channel, qName)

	return channel, userDeleteCheckQueue
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
