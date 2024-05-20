package dto

import (
	"github.com/ZMS-DevOps/hotel-service/domain"
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
