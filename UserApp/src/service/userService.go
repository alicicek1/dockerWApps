package userService

import (
	userRepository "UserApp/src/repository"
	userEntity "UserApp/src/type/entity"
	"UserApp/src/type/util"
	"UserApp/src/type/util/client"
	"UserApp/src/type/util/rabbitmq"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/rabbitmq/amqp091-go"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type UserServiceType struct {
	UserRepository userRepository.UserRepository
	Queue          amqp091.Queue
	Channel        *amqp091.Channel
	Client         client.Client
}

func NewUserService(r userRepository.UserRepository, channel *amqp091.Channel, queue amqp091.Queue, ticketClient client.Client) UserServiceType {
	return UserServiceType{UserRepository: r, Queue: queue, Channel: channel, Client: ticketClient}
}

type UserService interface {
	UserServiceInsert(user userEntity.UserPostRequestModel, userId string) (*util.PostResponseModel, *util.Error)
	UserServiceGetById(id string) (*userEntity.User, *util.Error)
	UserServiceDeleteById(id string) (util.DeleteResponseType, *util.Error)
	UserServiceGetAll(filter util.Filter) (*util.GetAllResponseType, *util.Error)
	UserServiceLogin(loginRequestModel userEntity.LoginRequestModel) (*userEntity.LoginResponseModel, *util.Error)
	UserIfExistById(id string) (bool, *util.Error)
}

func (u UserServiceType) UserServiceInsert(userPostRequestModel userEntity.UserPostRequestModel, userId string) (*util.PostResponseModel, *util.Error) {
	user := userEntity.User{Username: userPostRequestModel.Username,
		Password: userPostRequestModel.Password,
		Email:    userPostRequestModel.Email,
		Type:     userPostRequestModel.Type,
		Age:      userPostRequestModel.Age,
	}

	if userId != "" {
		user.Id = userId
	}

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
		return result, err
	}

	go CheckUserTicketCountAndPublishToQueue(u.Client, id, u.Channel, u.Queue)

	return result, nil
}

func CheckUserTicketCountAndPublishToQueue(c client.Client, id string, channel *amqp091.Channel, queue amqp091.Queue) {
	count := GetUserTicketCounts(c, id)
	for i := 0; i < count; i++ {
		val := rabbitmq.PublishToQueue(channel, queue, id)
		if !val {
			log.Error("An error occurred while publishing message", id)
		}
	}
}

func GetUserTicketCounts(c client.Client, id string) int {
	path := "getCountByCreatedId/" + id
	resp, err := c.Get(path, nil)
	if err != nil {
		return 0
	}

	resp = strings.Trim(resp, "\n")
	if res, e := strconv.Atoi(resp); e != nil {
		return 0
	} else {
		return res
	}
}
func (u UserServiceType) UserServiceGetAll(filter util.Filter) (*util.GetAllResponseType, *util.Error) {
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

	res, errTkn := util.CreateToken(result.Id)
	if errTkn != nil {
		return nil, util.NewError("user service", "TOKEN CREATE", errTkn.Error(), http.StatusBadRequest, 5003)
	}

	return res, nil
}
func (u UserServiceType) UserIfExistById(id string) (bool, *util.Error) {
	result, err := u.UserRepository.UserIfExistById(id)
	if err != nil {
		return false, err
	}
	return result, nil
}
