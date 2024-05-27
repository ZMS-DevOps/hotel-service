package application

import (
	"github.com/ZMS-DevOps/hotel-service/domain"
	"time"
)

func AddSpecialPrice(specialPrices []domain.SpecialPrice, newSpecialPrice domain.SpecialPrice) []domain.SpecialPrice {

	var updatedSpecialPrices []domain.SpecialPrice
	for _, sp := range specialPrices {
		if overlap(sp.DateRange, newSpecialPrice.DateRange) {
			if sp.DateRange.Start.Before(newSpecialPrice.DateRange.Start) {
				updatedSpecialPrices = updateSpecialPrices(updatedSpecialPrices, sp.Price, sp.DateRange.Start, newSpecialPrice.DateRange.Start)
			}
			if sp.DateRange.End.After(newSpecialPrice.DateRange.End) {
				updatedSpecialPrices = updateSpecialPrices(updatedSpecialPrices, sp.Price, newSpecialPrice.DateRange.End, sp.DateRange.End)
			}
		} else {
			updatedSpecialPrices = append(updatedSpecialPrices, sp)
		}
	}

	updatedSpecialPrices = append(updatedSpecialPrices, newSpecialPrice)
	return updatedSpecialPrices
}

func updateSpecialPrices(updatedSpecialPrices []domain.SpecialPrice, price float32, start time.Time, end time.Time) []domain.SpecialPrice {
	updatedSpecialPrices = append(updatedSpecialPrices, domain.SpecialPrice{
		Price: price,
		DateRange: domain.DateRange{
			Start: start,
			End:   end,
		},
	})
	return updatedSpecialPrices
}

func overlap(a, b domain.DateRange) bool {
	return a.Start.Before(b.End) && b.Start.Before(a.End)
}
