package dto

import (
	"encoding/base64"
	"github.com/ZMS-DevOps/hotel-service/domain"
	search "github.com/ZMS-DevOps/search-service/proto"
	"io/ioutil"
	"time"
)

func MapAccommodation(accommodation *AccommodationDto) *domain.Accommodation {
	accommodationPb := &domain.Accommodation{
		HostId:                                accommodation.HostId,
		Name:                                  accommodation.Name,
		Location:                              accommodation.Location,
		Benefits:                              accommodation.Benefits,
		Photos:                                accommodation.Photos,
		GuestNumber:                           mapGuestNumberDto(&accommodation.GuestNumber),
		DefaultPrice:                          mapDefaultPriceDto(&accommodation.DefaultPrice),
		ReviewReservationRequestAutomatically: accommodation.ReviewReservationRequestAutomatically,
	}
	return accommodationPb
}

func mapGuestNumberDto(guestNumber *GuestNumberDto) domain.GuestNumber {
	return domain.GuestNumber{
		Min: guestNumber.Min,
		Max: guestNumber.Max,
	}
}

func mapGuestNumber(guestNumber *domain.GuestNumber) GuestNumberDto {
	return GuestNumberDto{
		Min: guestNumber.Min,
		Max: guestNumber.Max,
	}
}

func mapDefaultPriceDto(defaultPrice *DefaultPriceDto) domain.DefaultPrice {
	return domain.DefaultPrice{
		Price: defaultPrice.Price,
		Type:  *MapPricingType(&defaultPrice.Type),
	}
}

func mapDefaultPrice(defaultPrice *domain.DefaultPrice) DefaultPriceDto {
	return DefaultPriceDto{
		Price: defaultPrice.Price,
		Type:  defaultPrice.Type.String(),
	}
}

func MapPricingType(typeOfPayment *string) *domain.PricingType {
	var pricingType domain.PricingType
	switch *typeOfPayment {
	case "PerApartmentUnit":
		pricingType = domain.PerApartmentUnit
	case "PerGuest":
		pricingType = domain.PerGuest
	default:
		return nil
	}
	return &pricingType
}

func MapToSearchAccommodation(accommodation *domain.Accommodation) *search.Accommodation {
	var mainPhoto = ""
	if (accommodation.Photos != nil) && (len(accommodation.Photos) != 0) {
		mainPhoto = accommodation.Photos[0]
	}
	return &search.Accommodation{
		AccommodationId: accommodation.Id.Hex(),
		Name:            accommodation.Name,
		Location:        accommodation.Location,
		MainPhoto:       mainPhoto,
		MinGuestNumber:  int32(accommodation.GuestNumber.Min),
		MaxGuestNumber:  int32(accommodation.GuestNumber.Max),
		DefaultPrice:    accommodation.DefaultPrice.Price,
		PriceType:       accommodation.DefaultPrice.Type.String(),
		SpecialPrice:    mapSearchSpecialPrice(accommodation.SpecialPrice),
		HostId:          accommodation.HostId,
	}
}

func mapSearchSpecialPrice(price []domain.SpecialPrice) []*search.SpecialPrice {
	var result []*search.SpecialPrice

	for _, p := range price {
		result = addSpecialPrice(p, result)
	}

	return result
}

func addSpecialPrice(p domain.SpecialPrice, result []*search.SpecialPrice) []*search.SpecialPrice {
	sp := &search.SpecialPrice{
		Price:     p.Price,
		StartDate: p.DateRange.Start.Format(time.RFC3339),
		EndDate:   p.DateRange.End.Format(time.RFC3339),
	}
	result = append(result, sp)
	return result
}

func toDateRangeDto(dataRange domain.DateRange) DateRangeDto {
	return DateRangeDto{
		Start: dataRange.Start,
		End:   dataRange.End,
	}
}

func toSpecialPriceDto(specialPrice []domain.SpecialPrice) []SpecialPriceDto {
	var result []SpecialPriceDto
	for _, p := range specialPrice {
		result = append(result, SpecialPriceDto{
			Price:     p.Price,
			DateRange: toDateRangeDto(p.DateRange),
		})
	}
	return result
}

func encodeFileToBase64(filePath string) (string, error) {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(fileBytes), nil
}

func MapAccommodationResponse(accommodation domain.Accommodation) *AccommodationResponse {
	base64Photos := make([]string, 0, len(accommodation.Photos))
	for _, photoPath := range accommodation.Photos {
		base64Photo, err := encodeFileToBase64(photoPath)
		if err != nil {
			continue
		}
		base64Photos = append(base64Photos, base64Photo)
	}

	return &AccommodationResponse{
		Id:                                    accommodation.Id.Hex(),
		Name:                                  accommodation.Name,
		Location:                              accommodation.Location,
		Benefits:                              accommodation.Benefits,
		Photos:                                base64Photos,
		GuestNumber:                           mapGuestNumber(&accommodation.GuestNumber),
		DefaultPrice:                          mapDefaultPrice(&accommodation.DefaultPrice),
		SpecialPrice:                          toSpecialPriceDto(accommodation.SpecialPrice),
		HostId:                                accommodation.HostId,
		ReviewReservationRequestAutomatically: accommodation.ReviewReservationRequestAutomatically,
	}
}
