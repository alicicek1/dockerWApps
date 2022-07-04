package repository

import (
	"TicketApp/src/type/entity"
	util "TicketApp/src/type/util"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	totalCount, err := t.TicketCollection.CountDocuments(ctx, filter.Filters)
	if err != nil {
		return nil, util.CountGet.ModifyApplicationName("ticket repository").ModifyDescription(err.Error()).ModifyErrorCode(3001)
	}
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
		RowCount: totalCount,
		Models:   tickets,
	}, nil
}
