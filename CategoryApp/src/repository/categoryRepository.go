package categoryRepository

import (
	categoryType "CategoryApp/src/type"
	"CategoryApp/src/type/util"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

type CategoryRepositoryType struct {
	CategoryCollection *mongo.Collection
}

func NewCategoryRepository(categoryCollection *mongo.Collection) CategoryRepositoryType {
	return CategoryRepositoryType{CategoryCollection: categoryCollection}
}

type CategoryRepository interface {
	CategoryRepoInsert(category categoryType.Category) (*util.PostResponseModel, *util.Error)
	CategoryRepoGetById(id string) (*categoryType.Category, *util.Error)
	CategoryRepoDeleteById(id string) (util.DeleteResponseType, *util.Error)
	CategoryRepositoryGetAll(filter util.Filter) (*util.GetAllResponseType, *util.Error)
	CategoryIfExistById(id string) (bool, *util.Error)
}

func (c CategoryRepositoryType) CategoryRepoInsert(category categoryType.Category) (*util.PostResponseModel, *util.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := c.CategoryCollection.InsertOne(ctx, category)
	if err != nil {
		return nil, util.UpsertFailed.ModifyApplicationName("user repository").ModifyErrorCode(4015)
	}
	return &util.PostResponseModel{Id: category.Id}, nil
}
func (c CategoryRepositoryType) CategoryRepoGetById(id string) (*categoryType.Category, *util.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var category categoryType.Category
	filter := bson.M{"_id": id}
	if err := c.CategoryCollection.FindOne(ctx, filter).Decode(&category); err != nil {
		return nil, util.NotFound.ModifyApplicationName("category repository").ModifyErrorCode(4028)
	}
	return &category, nil
}
func (c CategoryRepositoryType) CategoryRepoDeleteById(id string) (util.DeleteResponseType, *util.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	result, err := c.CategoryCollection.DeleteOne(ctx, filter)
	if err != nil || result.DeletedCount <= 0 {
		return util.DeleteResponseType{IsSuccess: false}, util.DeleteFailed.ModifyApplicationName("category repository").ModifyErrorCode(4029)
	}
	return util.DeleteResponseType{IsSuccess: true}, nil
}
func (c CategoryRepositoryType) CategoryRepositoryGetAll(filter util.Filter) (*util.GetAllResponseType, *util.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	totalCount, err := c.CategoryCollection.CountDocuments(ctx, filter.Filters)
	if err != nil {
		return nil, util.CountGet.ModifyApplicationName("category repository").ModifyDescription(err.Error()).ModifyErrorCode(3000)
	}
	opts := options.Find().SetSkip(filter.Page).SetLimit(filter.PageSize)
	if filter.SortingField != "" && filter.SortingDirection != 0 {
		opts.SetSort(bson.D{{filter.SortingField, filter.SortingDirection}})
	}

	cur, err := c.CategoryCollection.Find(ctx, filter.Filters, opts)
	if err != nil {

	}
	var categories []categoryType.Category
	err = cur.All(ctx, &categories)
	return &util.GetAllResponseType{
		RowCount: totalCount,
		Models:   categories,
	}, nil
}

func (c CategoryRepositoryType) CategoryIfExistById(id string) (bool, *util.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	count, err := c.CategoryCollection.CountDocuments(ctx, bson.M{"_id": id})
	if err != nil {
		return false, util.UnKnownError.ModifyApplicationName("Category Repository").ModifyOperation("Count document by id").ModifyDescription(err.Error()).ModifyErrorCode(3032).ModifyStatusCode(http.StatusBadRequest)
	}

	if count > 0 {
		return true, nil
	}

	return false, util.UnKnownError.ModifyApplicationName("Category Repository").ModifyOperation("Count document by id").ModifyDescription("There is no user with provided identified.").ModifyErrorCode(3033).ModifyStatusCode(http.StatusBadRequest)
}
