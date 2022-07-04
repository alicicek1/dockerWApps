package categoryService

import (
	categoryRepository "CategoryApp/src/repository"
	categoryType "CategoryApp/src/type"
	"CategoryApp/src/type/util"
	"github.com/google/uuid"
	"time"
)

type CategoryServiceType struct {
	CategoryRepository categoryRepository.CategoryRepository
}

type CategoryService interface {
	CategoryServiceInsert(user categoryType.Category) (*util.PostResponseModel, *util.Error)
	CategoryServiceGetById(id string) (*categoryType.Category, *util.Error)
	CategoryServiceDeleteById(id string) (util.DeleteResponseType, *util.Error)
	CategoryServiceGetAll(filter util.Filter) (*util.GetAllResponseType, *util.Error)
	CategoryIfExistById(id string) (bool, *util.Error)
}

func NewCategoryService(r categoryRepository.CategoryRepository) CategoryServiceType {
	return CategoryServiceType{CategoryRepository: r}
}

func (c CategoryServiceType) CategoryServiceInsert(category categoryType.Category) (*util.PostResponseModel, *util.Error) {
	if category.Id == "" {
		isSuccess, err := util.CheckCategoryModel(category)
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
func (c CategoryServiceType) CategoryServiceGetById(id string) (*categoryType.Category, *util.Error) {
	result, err := c.CategoryRepository.CategoryRepoGetById(id)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (c CategoryServiceType) CategoryServiceDeleteById(id string) (util.DeleteResponseType, *util.Error) {
	result, err := c.CategoryRepository.CategoryRepoDeleteById(id)
	if err != nil || result.IsSuccess == false {
		return util.DeleteResponseType{IsSuccess: false}, err
	}
	return util.DeleteResponseType{IsSuccess: true}, nil
}
func (c CategoryServiceType) CategoryServiceGetAll(filter util.Filter) (*util.GetAllResponseType, *util.Error) {
	result, err := c.CategoryRepository.CategoryRepositoryGetAll(filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c CategoryServiceType) CategoryIfExistById(id string) (bool, *util.Error) {
	result, err := c.CategoryRepository.CategoryIfExistById(id)
	if err != nil {
		return false, err
	}
	return result, nil
}
