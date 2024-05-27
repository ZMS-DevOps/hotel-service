package dto

import (
	"github.com/ZMS-DevOps/hotel-service/domain"
	search "github.com/ZMS-DevOps/search-service/proto"
	"time"
)

func MapAccommodation(accommodation *AccommodationDto) *domain.Accommodation {
	accommodationPb := &domain.Accommodation{
		Name:         accommodation.Name,
		Location:     accommodation.Location,
		Benefits:     accommodation.Benefits,
		Photos:       accommodation.Photos,
		GuestNumber:  mapGuestNumber(&accommodation.GuestNumber),
		DefaultPrice: mapDefaultPrice(&accommodation.DefaultPrice),
	}
	return accommodationPb
}

func mapGuestNumber(guestNumber *GuestNumberDto) domain.GuestNumber {
	return domain.GuestNumber{
		Min: guestNumber.Min,
		Max: guestNumber.Max,
	}
}

func mapDefaultPrice(defaultPrice *DefaultPriceDto) domain.DefaultPrice {
	return domain.DefaultPrice{
		Price: defaultPrice.Price,
		Type:  *MapPricingType(&defaultPrice.Type),
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
	return &search.Accommodation{
		AccommodationId: accommodation.Id.Hex(),
		Name:            accommodation.Name,
		Location:        accommodation.Location,
		MainPhoto:       accommodation.Photos[0],
		MinGuestNumber:  int32(accommodation.GuestNumber.Min),
		MaxGuestNumber:  int32(accommodation.GuestNumber.Max),
		DefaultPrice:    accommodation.DefaultPrice.Price,
		PriceType:       string(rune(accommodation.DefaultPrice.Type)),
		SpecialPrice:    mapSearchSpecialPrice(accommodation.SpecialPrice),
	}
}

func mapSearchSpecialPrice(price []domain.SpecialPrice) []*search.SpecialPrice {
	var result []*search.SpecialPrice

	for _, p := range price {
		sp := &search.SpecialPrice{
			Price:     p.Price,
			StartDate: p.DateRange.Start.Format(time.RFC3339),
			EndDate:   p.DateRange.End.Format(time.RFC3339),
		}
		result = append(result, sp)
	}

	return result
}
