package dto

import (
	"github.com/mmmajder/zms-devops-hotel-service/domain"
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
		Type:  domain.PerApartmentUnit, // todo fix
	}
}
