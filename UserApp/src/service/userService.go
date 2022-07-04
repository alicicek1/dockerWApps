package userService

import (
	userRepository "UserApp/src/repository"
	userEntity "UserApp/src/type/entity"
	"UserApp/src/type/util"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type UserServiceType struct {
	UserRepository userRepository.UserRepository
}

func NewUserService(r userRepository.UserRepository) UserServiceType {
	return UserServiceType{UserRepository: r}
}

type UserService interface {
	UserServiceInsert(user userEntity.User) (*userEntity.UserPostResponseModel, *util.Error)
	UserServiceGetById(id string) (*userEntity.User, *util.Error)
	UserServiceDeleteById(id string) (util.DeleteResponseType, *util.Error)
	UserServiceGetAll(filter util.Filter) (*userEntity.UserGetResponseModel, *util.Error)
	UserServiceLogin(loginRequestModel userEntity.LoginRequestModel) (*userEntity.LoginResponseModel, *util.Error)
}

func (u UserServiceType) UserServiceInsert(user userEntity.User) (*userEntity.UserPostResponseModel, *util.Error) {
	if user.Id == "" {
		isSuccess, err := util.CheckUserModel(user)
		if !isSuccess {
			return nil, err
		}
	}

	if user.Password != "" {
		user.Password = util.GetMD5Hash(user.Password)
	}
	user.Type = userEntity.DEFAULT

	user.CreatedAt = time.Now()
	if user.Id == "" {
		user.Id = uuid.New().String()
	}
	user.UpdatedAt = time.Now()

	result, err := u.UserRepository.UserRepoInsert(user)

	return result, err
}
func (u UserServiceType) UserServiceGetById(id string) (*userEntity.User, *util.Error) {
	result, err := u.UserRepository.UserRepoGetById(id)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (u UserServiceType) UserServiceDeleteById(id string) (util.DeleteResponseType, *util.Error) {
	result, err := u.UserRepository.UserRepoDeleteById(id)
	if err != nil || result.IsSuccess == false {
		return util.DeleteResponseType{IsSuccess: false}, err
	}
	return util.DeleteResponseType{IsSuccess: true}, nil
}
func (u UserServiceType) UserServiceGetAll(filter util.Filter) (*userEntity.UserGetResponseModel, *util.Error) {
	result, err := u.UserRepository.UserRepositoryGetAll(filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (u UserServiceType) UserServiceLogin(loginRequestModel userEntity.LoginRequestModel) (*userEntity.LoginResponseModel, *util.Error) {
	*loginRequestModel.Password = util.GetMD5Hash(*loginRequestModel.Password)
	result, err := u.UserRepository.UserRepositoryFindByUsernameAndPassword(loginRequestModel)
	if err != nil || result == nil {
		return nil, err
	}

	res, errTkn := util.CreateToken(loginRequestModel)
	if errTkn != nil {
		return nil, util.NewError("user service", "TOKEN CREATE", errTkn.Error(), http.StatusBadRequest, 5003)
	}

	return res, nil
}
