package service

import (
	"TicketApp/src/repository"
	"TicketApp/src/type/entity"
	"TicketApp/src/type/util"
	"TicketApp/src/type/util/rabbitmq"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/rabbitmq/amqp091-go"
	"strings"
	"time"
)

type TicketServiceType struct {
	TicketRepository repository.TicketRepository
	UserClient       util.Client
	CategoryClient   util.Client
	Channel          *amqp091.Channel
	Queue            amqp091.Queue
}

type TicketService interface {
	TicketServiceInsert(ticket entity.TicketPostRequestModel, creatorId string) (*util.PostResponseModel, *util.Error)
	TicketServiceGetById(id string) (*entity.Ticket, *util.Error)
	TicketServiceDeleteById(id string) (util.DeleteResponseType, *util.Error)
	TicketServiceGetAll(filter util.Filter) (*util.GetAllResponseType, *util.Error)
}

func NewTicketService(r repository.TicketRepository, userClient util.Client, categoryClient util.Client, channel *amqp091.Channel, queue amqp091.Queue) TicketServiceType {
	return TicketServiceType{
		TicketRepository: r,
		UserClient:       userClient,
		CategoryClient:   categoryClient,
		Channel:          channel,
		Queue:            queue,
	}
}

func (t TicketServiceType) TicketServiceInsert(ticketPostRequestModel entity.TicketPostRequestModel, token string) (*util.PostResponseModel, *util.Error) {

	creatorId := util.DecodeTokenReturnsUserId(token)
	ticket := entity.Ticket{
		CategoryId:     ticketPostRequestModel.CategoryId,
		Attachments:    ticketPostRequestModel.Attachments,
		Answers:        ticketPostRequestModel.Answers,
		Subject:        ticketPostRequestModel.Subject,
		Body:           ticketPostRequestModel.Body,
		CreatedBy:      creatorId,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		LastAnsweredAt: time.Now(),
		Status:         byte(entity.CREATED),
	}

	if ticket.Id == "" {
		isSuccess, err := util.CheckTicketModel(ticket, t.UserClient, t.CategoryClient)
		if !isSuccess {
			return nil, err
		}
	}

	ticket.Id = uuid.New().String()
	ticket.CreatedAt = time.Now()
	ticket.UpdatedAt = time.Now()

	result, err := t.TicketRepository.TicketRepoInsert(ticket)

	return result, err
}
func (t TicketServiceType) TicketServiceGetById(id string) (*entity.Ticket, *util.Error) {
	go t.CheckUserDeleteQueueForUpdate(t.Channel, t.Queue)

	result, err := t.TicketRepository.TicketRepoGetById(id)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (t TicketServiceType) CheckUserDeleteQueueForUpdate(channel *amqp091.Channel, queue amqp091.Queue) {
	msgS := rabbitmq.ConsumeMessage(channel, queue)
	for msg := range msgS {
		a := strings.ReplaceAll(string(msg.Body), "\"", "")
		a = strings.ReplaceAll(a, "\r\n", "")
		a = strings.ReplaceAll(a, "\r", "")
		a = strings.ReplaceAll(a, "\n", "")
		fmt.Println(a)
		_, err := t.TicketRepository.UpdateDeletedTicket(a)
		if err != nil {
			log.Error(err)
		}
	}
}
func (t TicketServiceType) TicketServiceDeleteById(id string) (util.DeleteResponseType, *util.Error) {
	result, err := t.TicketRepository.TicketRepoDeleteById(id)
	if err != nil || result.IsSuccess == false {
		return result, err
	}
	return result, nil
}
func (t TicketServiceType) TicketServiceGetAll(filter util.Filter) (*util.GetAllResponseType, *util.Error) {
	result, err := t.TicketRepository.TicketRepositoryGetAll(filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}
