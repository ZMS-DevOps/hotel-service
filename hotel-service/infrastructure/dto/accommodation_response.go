package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AccommodationResponse struct {
	Id                                    primitive.ObjectID `bson:"id"`
	HostId                                primitive.ObjectID `bson:"host_id"`
	Name                                  string             `bson:"name"`
	Location                              string             `bson:"location"`
	Benefits                              []string           `bson:"benefits"`
	Photos                                []string           `bson:"photos"`
	GuestNumber                           GuestNumberDto     `bson:"guest_number"`
	DefaultPrice                          DefaultPriceDto    `bson:"default_price"`
	SpecialPrice                          []SpecialPriceDto  `bson:"special_price"`
	ReviewReservationRequestAutomatically bool               `bson:"review_reservation_request_automatically"`
}

type GuestNumber struct {
	Min int `bson:"min"`
	Max int `bson:"max"`
}

type SpecialPriceDto struct {
	Price     float32      `bson:"price"`
	DateRange DateRangeDto `bson:"date_range"`
}

type DateRangeDto struct {
	Start time.Time `bson:"start"`
	End   time.Time `bson:"end"`
}
