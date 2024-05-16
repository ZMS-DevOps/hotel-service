package application

import (
	"github.com/mmmajder/zms-devops-hotel-service/domain"
	"github.com/mmmajder/zms-devops-hotel-service/infrastructure/dto"
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

func (service *AccommodationService) UpdatePrice(id primitive.ObjectID, updatePriceDto dto.UpdatePriceDto) interface{} {
	_, err := service.store.Get(id)
	if err != nil {
		return err
	}
	if updatePriceDto.DateRange == nil && updatePriceDto.Price != nil {
		if err := service.store.UpdateDefaultPrice(id, updatePriceDto.Price); err != nil {
			return err
		}
	} else if updatePriceDto.Price != nil {
		// todo check if there are no reservations if that date range
		if err := service.store.UpdateSpecialPrice(id, updatePriceDto.Price, &updatePriceDto.DateRange.Start, &updatePriceDto.DateRange.End); err != nil {
			return err
		}
	}
	if updatePriceDto.Type != nil {
		// todo check if there are no reservations for accommodation at all
		if err := service.store.UpdateTypeOfPayment(id, updatePriceDto.Type); err != nil {
			return err
		}
	}
	return nil
}
