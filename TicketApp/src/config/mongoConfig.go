package ticketConfig

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/net/context"
	"os"
	"time"
)

type AppConfig struct {
	Env                   string
	MongoClientUri        string
	DBName                string
	UserColName           string
	TicketColName         string
	CategoryColName       string
	MongoDuration         int16
	MaxPageLimit          int
	RabbitUsername        string
	RabbitPassword        string
	RabbitUri             string
	UserDeleteCheckQName  string
	CategoryContainerName string
	TicketContainerName   string
	UserContainerName     string
}

var EnvConfig = map[string]AppConfig{
	"local": {
		Env:                   "local",
		MongoClientUri:        "mongodb://localhost:27017",
		DBName:                "TicketApp",
		UserColName:           "User",
		TicketColName:         "Ticket",
		CategoryColName:       "Category",
		MongoDuration:         5,
		MaxPageLimit:          100,
		RabbitUsername:        "guest",
		RabbitPassword:        "guest",
		RabbitUri:             "localhost:5672",
		UserDeleteCheckQName:  "userDelCheck",
		CategoryContainerName: "localhost",
		TicketContainerName:   "localhost",
		UserContainerName:     "localhost",
	},
	"qa": {
		Env:                   "qa",
		MongoClientUri:        "mongodb://mongo:27017",
		DBName:                "TicketApp",
		UserColName:           "User",
		TicketColName:         "Ticket",
		CategoryColName:       "Category",
		MongoDuration:         5,
		MaxPageLimit:          100,
		RabbitUsername:        "guest",
		RabbitPassword:        "guest",
		RabbitUri:             "rabbitmq:5672",
		UserDeleteCheckQName:  "userDelCheck",
		CategoryContainerName: "category_service",
		TicketContainerName:   "ticket_service",
		UserContainerName:     "user_service",
	},
	"prod": {
		Env:                   "qa",
		MongoClientUri:        "mongodb+srv://admin:1@cluster0.ymrmq.mongodb.net/?retryWrites=true&w=majority",
		DBName:                "TicketApp",
		UserColName:           "User",
		TicketColName:         "Ticket",
		CategoryColName:       "Category",
		MongoDuration:         5,
		MaxPageLimit:          100,
		RabbitUsername:        "admin",
		RabbitPassword:        "1",
		RabbitUri:             "rabbitmq:5672",
		UserDeleteCheckQName:  "userDelCheck",
		CategoryContainerName: "category_service",
		TicketContainerName:   "ticket_service",
		UserContainerName:     "user_service",
	},
}

func NewAppConfig() AppConfig {
	return AppConfig{}
}

func (mCfg *AppConfig) CloseConnection(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func (mCfg *AppConfig) ConnectDatabase() (*mongo.Client, context.Context, context.CancelFunc, *AppConfig) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	cfg := GetConfigModel()
	//conStr := "mongodb+srv://admin:1@cluster0.ymrmq.mongodb.net/?retryWrites=true&w=majority"
	if cfg.MongoClientUri == "" {
		panic("Connection string was not found. Check the .env file.")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoClientUri))
	if err != nil {
		panic(err)
	}

	return client, ctx, cancel, &cfg
}

func (mCfg *AppConfig) Ping(client *mongo.Client, ctx context.Context) error {
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("Connected successfully.")
	return nil
}

func (mCfg *AppConfig) GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("TicketApp").Collection(collectionName)
}

func GetConfigModel() AppConfig {

	env := os.Getenv("Env")
	if env == "" {
		panic("Env was not found.")
	}

	model, exist := EnvConfig[env]
	if !exist {
		panic("There is no model with provided environment.")
	}

	return model
}

func (mCfg *AppConfig) GetRabbitMqDialConnectionUri() string {
	return "amqp://" + mCfg.RabbitUsername + ":" + mCfg.RabbitPassword + "@" + mCfg.RabbitUri + "/"
}
