package dto

import (
	"time"
)

type AccommodationResponse struct {
	Id                                    string            `json:"id"`
	HostId                                string            `json:"host_id"`
	Name                                  string            `json:"name"`
	Location                              string            `json:"location"`
	Benefits                              []string          `json:"benefits"`
	Photos                                []string          `json:"photos"`
	GuestNumber                           GuestNumberDto    `json:"guest_number"`
	DefaultPrice                          DefaultPriceDto   `json:"default_price"`
	SpecialPrice                          []SpecialPriceDto `json:"special_price"`
	ReviewReservationRequestAutomatically bool              `json:"review_reservation_request_automatically"`
}

type GuestNumber struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type SpecialPriceDto struct {
	Price     float32      `json:"price"`
	DateRange DateRangeDto `json:"date_range"`
}

type DateRangeDto struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}
