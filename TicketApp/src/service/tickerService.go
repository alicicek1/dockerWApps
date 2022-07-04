package service

import (
	"TicketApp/src/repository"
	"TicketApp/src/type/entity"
	"TicketApp/src/type/util"
	"github.com/google/uuid"
	"time"
)

type TicketServiceType struct {
	TicketRepository repository.TicketRepository
}

type TicketService interface {
	TicketServiceInsert(user entity.Ticket) (*util.PostResponseModel, *util.Error)
	TicketServiceGetById(id string) (*entity.Ticket, *util.Error)
	TicketServiceDeleteById(id string) (util.DeleteResponseType, *util.Error)
	TicketServiceGetAll(filter util.Filter) (*util.GetAllResponseType, *util.Error)
}

func NewTicketService(r repository.TicketRepository) TicketServiceType {
	return TicketServiceType{TicketRepository: r}
}

func (t TicketServiceType) TicketServiceInsert(ticket entity.Ticket) (*util.PostResponseModel, *util.Error) {
	if ticket.Id == "" {
		isSuccess, err := util.CheckTicketModel(ticket)
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
	result, err := t.TicketRepository.TicketRepoGetById(id)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (t TicketServiceType) TicketServiceDeleteById(id string) (util.DeleteResponseType, *util.Error) {
	result, err := t.TicketRepository.TicketRepoDeleteById(id)
	if err != nil || result.IsSuccess == false {
		return util.DeleteResponseType{IsSuccess: false}, err
	}
	return util.DeleteResponseType{IsSuccess: true}, nil
}
func (t TicketServiceType) TicketServiceGetAll(filter util.Filter) (*util.GetAllResponseType, *util.Error) {
	result, err := t.TicketRepository.TicketRepositoryGetAll(filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}
