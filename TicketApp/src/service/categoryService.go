package service

import (
	"TicketApp/src/repository"
	"TicketApp/src/type/entity"
	util2 "TicketApp/src/type/util"
	"github.com/google/uuid"
	"time"
)

type CategoryServiceType struct {
	CategoryRepository repository.CategoryRepository
}

type CategoryService interface {
	CategoryServiceInsert(user entity.Category) (*entity.CategoryPostResponseModel, *util2.Error)
	CategoryServiceGetById(id string) (*entity.Category, *util2.Error)
	CategoryServiceDeleteById(id string) (util2.DeleteResponseType, *util2.Error)
	CategoryServiceGetAll(filter util2.Filter) (*entity.CategoryGetResponseModel, *util2.Error)
}

func NewCategoryService(r repository.CategoryRepository) CategoryServiceType {
	return CategoryServiceType{CategoryRepository: r}
}

func (c CategoryServiceType) CategoryServiceInsert(category entity.Category) (*entity.CategoryPostResponseModel, *util2.Error) {
	if category.Id == "" {
		isSuccess, err := util2.CheckCategoryModel(category)
		if !isSuccess {
			return nil, err
		}
	}

	category.Id = uuid.New().String()
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	result, err := c.CategoryRepository.CategoryRepoInsert(category)

	return result, err
}
func (c CategoryServiceType) CategoryServiceGetById(id string) (*entity.Category, *util2.Error) {
	result, err := c.CategoryRepository.CategoryRepoGetById(id)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (c CategoryServiceType) CategoryServiceDeleteById(id string) (util2.DeleteResponseType, *util2.Error) {
	result, err := c.CategoryRepository.CategoryRepoDeleteById(id)
	if err != nil || result.IsSuccess == false {
		return util2.DeleteResponseType{IsSuccess: false}, err
	}
	return util2.DeleteResponseType{IsSuccess: true}, nil
}
func (c CategoryServiceType) CategoryServiceGetAll(filter util2.Filter) (*entity.CategoryGetResponseModel, *util2.Error) {
	result, err := c.CategoryRepository.CategoryRepositoryGetAll(filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}
