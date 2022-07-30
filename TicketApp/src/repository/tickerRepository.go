package repository

import (
	"TicketApp/src/type/entity"
	util "TicketApp/src/type/util"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

type TicketRepositoryType struct {
	TicketCollection *mongo.Collection
}

func NewTicketRepository(ticketCollection *mongo.Collection) TicketRepositoryType {
	return TicketRepositoryType{TicketCollection: ticketCollection}
}

type TicketRepository interface {
	TicketRepoInsert(ticket entity.Ticket) (*util.PostResponseModel, *util.Error)
	TicketRepoGetById(id string) (*entity.Ticket, *util.Error)
	TicketRepoDeleteById(id string) (util.DeleteResponseType, *util.Error)
	TicketRepositoryGetAll(filter util.Filter) (*util.GetAllResponseType, *util.Error)
	UpdateDeletedTicket(msg string) (interface{}, *util.Error)
	TicketRepoGetCountByCreatedId(id string) (int64, *util.Error)
}

func (t TicketRepositoryType) UpdateDeletedTicket(msg string) (interface{}, *util.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	opts := options.Update().SetUpsert(true)
	filter := bson.D{
		{"createdBy", strings.TrimRight(msg, "\n")},
		{"isDeleted", bson.D{
			{"$exists", false},
		}},
	}

	update := bson.D{{"$set", bson.D{{"isDeleted", true}}}}
	_, err := t.TicketCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, util.UpsertFailed.ModifyApplicationName("ticket repository").ModifyErrorCode(4055)
	}
	return &util.PostResponseModel{Id: msg}, nil
}

func (t TicketRepositoryType) TicketRepoInsert(ticket entity.Ticket) (*util.PostResponseModel, *util.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := t.TicketCollection.InsertOne(ctx, ticket)
	if err != nil {
		return nil, util.UpsertFailed.ModifyApplicationName("user repository").ModifyErrorCode(4015)
	}
	return &util.PostResponseModel{Id: ticket.Id}, nil
}

func (t TicketRepositoryType) TicketRepoGetCountByCreatedId(id string) (int64, *util.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"createdBy": id}
	if res, err := t.TicketCollection.CountDocuments(ctx, filter); err != nil {
		return 0, util.NotFound.ModifyApplicationName("ticket repository").ModifyErrorCode(4030)
	} else {
		return res, nil
	}
}

func (t TicketRepositoryType) TicketRepoGetById(id string) (*entity.Ticket, *util.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var ticket entity.Ticket
	filter := bson.M{"_id": id}
	if err := t.TicketCollection.FindOne(ctx, filter).Decode(&ticket); err != nil {
		return nil, util.NotFound.ModifyApplicationName("ticket repository").ModifyErrorCode(4030)
	}
	return &ticket, nil
}
func (t TicketRepositoryType) TicketRepoDeleteById(id string) (util.DeleteResponseType, *util.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	result, err := t.TicketCollection.DeleteOne(ctx, filter)
	if err != nil || result.DeletedCount <= 0 {
		return util.DeleteResponseType{IsSuccess: false}, util.DeleteFailed.ModifyApplicationName("ticket repository").ModifyErrorCode(4031)
	}
	return util.DeleteResponseType{IsSuccess: true}, nil
}
func (t TicketRepositoryType) TicketRepositoryGetAll(filter util.Filter) (*util.GetAllResponseType, *util.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	ticketCountChannel := make(chan int64)
	errorChannel := make(chan *util.Error)
	go t.GetTotalCount(ctx, filter.Filters, ticketCountChannel, errorChannel)

	select {
	case ticketCountVal, ok := <-ticketCountChannel:
		if ok {
			opts := options.Find().SetSkip(filter.Page).SetLimit(filter.PageSize)
			if filter.SortingField != "" && filter.SortingDirection != 0 {
				opts.SetSort(bson.D{{filter.SortingField, filter.SortingDirection}})
			}

			cur, err := t.TicketCollection.Find(ctx, filter.Filters, opts)
			if err != nil {

			}
			var tickets []entity.Ticket
			err = cur.All(ctx, &tickets)
			return &util.GetAllResponseType{
				RowCount: ticketCountVal,
				Models:   tickets,
			}, nil
		}
	case errorVal, ok := <-errorChannel:
		if ok {
			return nil, errorVal
		}
	}
	return nil, nil
}

func (t TicketRepositoryType) GetTotalCount(ctx context.Context, filters map[string]interface{}, countChannel chan int64, errorChannel chan *util.Error) {
	totalCount, err := t.TicketCollection.CountDocuments(ctx, filters)
	if err != nil {
		errorChannel <- util.CountGet.ModifyApplicationName("ticket repository").ModifyDescription(err.Error()).ModifyErrorCode(3001)
	} else {
		countChannel <- totalCount
	}
}
