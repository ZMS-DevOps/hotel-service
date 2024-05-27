package application

import (
	"github.com/ZMS-DevOps/hotel-service/domain"
)

func AddSpecialPrice(specialPrices []domain.SpecialPrice, newSpecialPrice domain.SpecialPrice) []domain.SpecialPrice {

	var updatedSpecialPrices []domain.SpecialPrice
	for _, sp := range specialPrices {
		if overlap(sp.DateRange, newSpecialPrice.DateRange) {
			if sp.DateRange.Start.Before(newSpecialPrice.DateRange.Start) {
				updatedSpecialPrices = append(updatedSpecialPrices, domain.SpecialPrice{
					Price: sp.Price,
					DateRange: domain.DateRange{
						Start: sp.DateRange.Start,
						End:   newSpecialPrice.DateRange.Start,
					},
				})
			}
			if sp.DateRange.End.After(newSpecialPrice.DateRange.End) {
				updatedSpecialPrices = append(updatedSpecialPrices, domain.SpecialPrice{
					Price: sp.Price,
					DateRange: domain.DateRange{
						Start: newSpecialPrice.DateRange.End,
						End:   sp.DateRange.End,
					},
				})
			}
		} else {
			updatedSpecialPrices = append(updatedSpecialPrices, sp)
		}
	}

	updatedSpecialPrices = append(updatedSpecialPrices, newSpecialPrice)
	return updatedSpecialPrices
}

func overlap(a, b domain.DateRange) bool {
	return a.Start.Before(b.End) && b.Start.Before(a.End)
}
