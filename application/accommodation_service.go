package application

import (
	"fmt"
	"github.com/mmmajder/zms-devops-hotel-service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccommodationService struct {
	store domain.AccommodationStore
}

func NewAccommodationService(store domain.AccommodationStore) *AccommodationService {
	return &AccommodationService{
		store: store,
	}
}

func (service *AccommodationService) Get(id primitive.ObjectID) (*domain.Accommodation, error) {
	return service.store.Get(id)
}

func (service *AccommodationService) GetAll() ([]*domain.Accommodation, error) {
	return service.store.GetAll()
}

func (service *AccommodationService) Add(accommodation *domain.Accommodation) error {
	err := service.store.Insert(accommodation)
	if err != nil {
		return err
	}
	return nil
}

func (service *AccommodationService) Update(id primitive.ObjectID, accommodation *domain.Accommodation) error {
	fmt.Println(id)
	_, err := service.store.Get(id)
	if err != nil {
		return err // Return error if accommodation does not exist
	}
	err = service.store.Update(id, accommodation)
	if err != nil {
		return err
	}

	return nil
}

func (service *AccommodationService) Delete(id primitive.ObjectID) error {
	_, err := service.store.Get(id)
	if err != nil {
		return err // Return error if accommodation does not exist
	}
	err = service.store.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
