package application

import (
	"github.com/ZMS-DevOps/hotel-service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	store domain.UserStore
}

func NewUserService(store domain.UserStore) *UserService {
	return &UserService{
		store: store,
	}
}

func (service *UserService) Get(id primitive.ObjectID) (*domain.User, error) {
	return service.store.Get(id)
}

func (service *UserService) Add(userId primitive.ObjectID) error {
	err := service.store.Insert(&domain.User{Id: userId})
	if err != nil {
		return err
	}
	return nil
}
