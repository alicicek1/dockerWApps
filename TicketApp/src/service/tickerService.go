package service

import (
	"TicketApp/src/repository"
	"TicketApp/src/type/entity"
	util2 "TicketApp/src/type/util"
	"github.com/google/uuid"
	"time"
)

type TicketServiceType struct {
	TicketRepository repository.TicketRepository
}

type TicketService interface {
	TicketServiceInsert(user entity.Ticket) (*entity.TicketPostResponseModel, *util2.Error)
	TicketServiceGetById(id string) (*entity.Ticket, *util2.Error)
	TicketServiceDeleteById(id string) (util2.DeleteResponseType, *util2.Error)
	TicketServiceGetAll(filter util2.Filter) (*entity.TicketGetReponseModel, *util2.Error)
}

func NewTicketService(r repository.TicketRepository) TicketServiceType {
	return TicketServiceType{TicketRepository: r}
}

func (t TicketServiceType) TicketServiceInsert(ticket entity.Ticket) (*entity.TicketPostResponseModel, *util2.Error) {
	if ticket.Id == "" {
		isSuccess, err := util2.CheckTicketModel(ticket)
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
func (t TicketServiceType) TicketServiceGetById(id string) (*entity.Ticket, *util2.Error) {
	result, err := t.TicketRepository.TicketRepoGetById(id)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (t TicketServiceType) TicketServiceDeleteById(id string) (util2.DeleteResponseType, *util2.Error) {
	result, err := t.TicketRepository.TicketRepoDeleteById(id)
	if err != nil || result.IsSuccess == false {
		return util2.DeleteResponseType{IsSuccess: false}, err
	}
	return util2.DeleteResponseType{IsSuccess: true}, nil
}
func (t TicketServiceType) TicketServiceGetAll(filter util2.Filter) (*entity.TicketGetReponseModel, *util2.Error) {
	result, err := t.TicketRepository.TicketRepositoryGetAll(filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}
