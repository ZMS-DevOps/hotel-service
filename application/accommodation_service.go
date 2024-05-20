package application

import (
	"fmt"
	booking "github.com/ZMS-DevOps/booking-service/proto"
	"github.com/ZMS-DevOps/hotel-service/application/external"

	"github.com/ZMS-DevOps/hotel-service/domain"
	"github.com/ZMS-DevOps/hotel-service/infrastructure/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccommodationService struct {
	store         domain.AccommodationStore
	bookingClient booking.BookingServiceClient
}

// func NewAccommodationService(store domain.AccommodationStore) *AccommodationService {
func NewAccommodationService(store domain.AccommodationStore, bookingClient booking.BookingServiceClient) *AccommodationService {
	return &AccommodationService{
		store:         store,
		bookingClient: bookingClient,
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
	// todo create booking empty unavailability object
	fmt.Println(accommodation.Id)
	fmt.Println("Stigao sammmmmmmm")
	_, err = external.CreateBookingUnavailability(service.bookingClient, accommodation.Id)
	fmt.Println(err)
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
		currentSpecialPrices, err := service.store.GetSpecialPrices(id)
		if err != nil {
			return err
		}

		newSpecialPrice := domain.SpecialPrice{
			Price: *updatePriceDto.Price,
			DateRange: domain.DateRange{
				Start: updatePriceDto.DateRange.Start,
				End:   updatePriceDto.DateRange.End,
			},
		}

		newSpecialPrices := AddSpecialPrice(currentSpecialPrices, newSpecialPrice)

		if err := service.store.UpdateSpecialPrice(id, newSpecialPrices); err != nil {
			return err
		}
	}
	if updatePriceDto.Type != nil {
		if err := service.store.UpdateTypeOfPayment(id, updatePriceDto.Type); err != nil {
			return err
		}
	}
	return nil
}
