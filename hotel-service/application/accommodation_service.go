package application

import (
	"encoding/base64"
	"errors"
	booking "github.com/ZMS-DevOps/booking-service/proto"
	"github.com/ZMS-DevOps/hotel-service/application/external"
	"github.com/ZMS-DevOps/hotel-service/domain"
	"github.com/ZMS-DevOps/hotel-service/infrastructure/dto"
	"github.com/ZMS-DevOps/hotel-service/util"
	search "github.com/ZMS-DevOps/search-service/proto"
	"github.com/afiskon/promtail-client/promtail"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/trace"
	"io/ioutil"
	"log"
)

type AccommodationService struct {
	store         domain.AccommodationStore
	bookingClient booking.BookingServiceClient
	searchClient  search.SearchServiceClient
	loki          promtail.Client
}

func NewAccommodationService(store domain.AccommodationStore, bookingClient booking.BookingServiceClient, searchClient search.SearchServiceClient, loki promtail.Client) *AccommodationService {
	return &AccommodationService{
		store:         store,
		bookingClient: bookingClient,
		searchClient:  searchClient,
		loki:          loki,
	}
}

func (service *AccommodationService) Get(id primitive.ObjectID, span trace.Span, loki promtail.Client) (*domain.Accommodation, error) {
	util.HttpTraceInfo("Fetching accommodation by id...", span, loki, "Add", "")
	return service.store.Get(id)
}

func (service *AccommodationService) GetAll(span trace.Span, loki promtail.Client) ([]*domain.Accommodation, error) {
	util.HttpTraceInfo("Fetching all accommodations...", span, loki, "Add", "")
	return service.store.GetAll()
}

func (service *AccommodationService) GetByHostId(ownerId string, span trace.Span, loki promtail.Client) ([]*domain.Accommodation, error) {
	util.HttpTraceInfo("Fetching accommodations by host id...", span, loki, "Add", "")
	return service.store.GetByHostId(ownerId)
}

func (service *AccommodationService) Add(accommodation *domain.Accommodation, span trace.Span, loki promtail.Client) error {
	util.HttpTraceInfo("Inserting accommodation...", span, loki, "Add", "")
	err := service.store.Insert(accommodation)
	if err != nil {
		return err
	}
	accommodation.SpecialPrice = []domain.SpecialPrice{}
	_, err = external.CreateBookingUnavailability(service.bookingClient, accommodation.Id, accommodation.ReviewReservationRequestAutomatically, accommodation.HostId, accommodation.Name, span, loki)
	if err != nil {
		return err
	}
	_, err = external.AddSearchAccommodation(service.searchClient, dto.MapToSearchAccommodation(accommodation), span, loki)
	if err != nil {
		return err
	}
	return nil
}

func (service *AccommodationService) Update(id primitive.ObjectID, accommodation *domain.Accommodation, span trace.Span, loki promtail.Client) error {
	util.HttpTraceInfo("Fetching accommodation by id...", span, loki, "Add", "")
	_, err := service.store.Get(id)
	if err != nil {
		util.HttpTraceError(err, "can't decode login payload", span, service.loki, "Login", "")
		return err
	}
	util.HttpTraceInfo("Updating accommodation...", span, loki, "Add", "")
	err = service.store.Update(id, accommodation)
	if err != nil {
		return err
	}
	_, err = external.UpdateBookingUnavailability(service.bookingClient, accommodation.Id, accommodation.ReviewReservationRequestAutomatically, accommodation.HostId, accommodation.Name, span, loki)
	_, err = external.EditSearchAccommodation(service.searchClient, dto.MapToSearchAccommodation(accommodation), span, loki)
	if err != nil {
		return err
	}
	return nil
}

func (service *AccommodationService) Delete(id primitive.ObjectID, span trace.Span, loki promtail.Client) error {
	util.HttpTraceInfo("Deleting accommodation...", span, loki, "Add", "")
	canDelete, err := external.CheckAccommodationHasReservation(service.bookingClient, id, span, loki)
	if err != nil {
		return err
	}
	if canDelete.Success {
		return service.deleteAccommodation(id, span, loki)
	}
	return errors.New("accommodation could not be deleted")
}

func (service *AccommodationService) deleteAccommodation(id primitive.ObjectID, span trace.Span, loki promtail.Client) error {
	util.HttpTraceInfo("Deleting accommodation from collection...", span, loki, "Add", "")
	if err := service.store.Delete(id); err != nil {
		return err
	}
	if err := service.notifySearchServiceWhenDeletingAccommodation(id, span, loki); err != nil {
		return err
	}
	return nil
}

func (service *AccommodationService) notifySearchServiceWhenDeletingAccommodation(id primitive.ObjectID, span trace.Span, loki promtail.Client) error {
	_, err := external.DeleteSearchAccommodation(service.searchClient, id, span, loki)
	if err != nil {
		return err
	}
	return nil
}

func (service *AccommodationService) UpdatePrice(id primitive.ObjectID, updatePriceDto dto.UpdatePriceDto, span trace.Span, loki promtail.Client) error {
	util.HttpTraceInfo("Fetching accommodation...", span, loki, "Add", "")
	_, err := service.store.Get(id)
	if err != nil {
		return err
	}

	if updatePriceDto.DateRange == nil && updatePriceDto.Price != nil {
		util.HttpTraceInfo("Updating accommodation default price...", span, loki, "Add", "")
		if err := updateDefaultPrice(id, updatePriceDto, service); err != nil {
			return err
		}
	} else if updatePriceDto.Price != nil {
		util.HttpTraceInfo("Updating accommodation special price...", span, loki, "Add", "")
		if err := updateSpecialPrice(id, updatePriceDto, service, span, loki); err != nil {
			return err
		}
	}
	if updatePriceDto.Type != nil {
		util.HttpTraceInfo("Updating accommodation payment type...", span, loki, "Add", "")
		if err := service.store.UpdateTypeOfPayment(id, updatePriceDto.Type); err != nil {
			return err
		}
	}

	util.HttpTraceInfo("Fetching accommodation...", span, loki, "Add", "")
	updatedAccommodation, err := service.store.Get(id)
	if err != nil {
		return err
	}
	_, err = external.EditSearchAccommodation(service.searchClient, dto.MapToSearchAccommodation(updatedAccommodation), span, loki)
	if err != nil {
		return err
	}

	return nil
}

func (service *AccommodationService) GetImages(accommodationIds []dto.GetImagesRequest, span trace.Span, loki promtail.Client) ([]dto.ImageResponse, error) {
	var images []dto.ImageResponse
	for _, accommodationId := range accommodationIds {
		id, err := primitive.ObjectIDFromHex(accommodationId.Id)
		if err != nil {
			return nil, err
		}
		util.HttpTraceInfo("Fetching accommodation...", span, loki, "Add", "")
		accommodation, err := service.Get(id, span, loki)
		if err != nil {
			return nil, err
		}
		encodedImaged, err := service.base64Encode(accommodation.Photos, span, loki)
		if err != nil {
			return nil, err
		}
		images = append(images, dto.ImageResponse{Id: id, Images: encodedImaged})
	}
	return images, nil
}

func (service *AccommodationService) base64Encode(photos []string, span trace.Span, loki promtail.Client) ([]string, error) {
	util.HttpTraceInfo("Encoding image base64...", span, loki, "Add", "")
	var base64Photos []string
	for _, photoPath := range photos {
		base64Photo, err := encodeFileToBase64(photoPath)
		if err != nil {
			continue
		}
		base64Photos = append(base64Photos, base64Photo)
	}

	return base64Photos, nil
}

func encodeFileToBase64(filePath string) (string, error) {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(fileBytes), nil
}

func (service *AccommodationService) OnDeleteAccommodations(hostId string, span trace.Span, loki promtail.Client) {
	util.HttpTraceInfo("Fetching accommodation by host id...", span, loki, "Add", "")
	accommodations, err := service.store.GetByHostId(hostId)
	if err != nil {
		log.Println(err)
		return
	}
	for _, accom := range accommodations {
		if err := service.deleteAccommodation(accom.Id, span, loki); err != nil {
			log.Println(err)
			return
		}
	}
}

func updateSpecialPrice(id primitive.ObjectID, updatePriceDto dto.UpdatePriceDto, service *AccommodationService, span trace.Span, loki promtail.Client) error {
	util.HttpTraceInfo("Fetching special prices by accommodation id...", span, loki, "Add", "")
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
	util.HttpTraceInfo("Updating special prices...", span, loki, "Add", "")
	if err := service.store.UpdateSpecialPrice(id, newSpecialPrices); err != nil {
		return err
	}
	return nil
}

func updateDefaultPrice(id primitive.ObjectID, updatePriceDto dto.UpdatePriceDto, service *AccommodationService) error {
	if err := service.store.UpdateDefaultPrice(id, updatePriceDto.Price); err != nil {
		return err
	}
	return nil
}
