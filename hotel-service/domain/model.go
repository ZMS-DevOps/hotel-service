package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Accommodation struct {
	Id           primitive.ObjectID `bson:"_id"`
	Name         string             `bson:"name"`
	Location     string             `bson:"location"`
	Benefits     []string           `bson:"benefits"`
	Photos       []string           `bson:"photos"`
	GuestNumber  GuestNumber        `bson:"guest_number"`
	DefaultPrice DefaultPrice       `bson:"default_price"`
	SpecialPrice []SpecialPrice     `bson:"special_price"`
}

type GuestNumber struct {
	Min int `bson:"min"`
	Max int `bson:"max"`
}

type PricingType int

const (
	PerApartmentUnit PricingType = iota
	PerGuest
)

type DefaultPrice struct {
	Price float32     `bson:"price"`
	Type  PricingType `bson:"type"`
}

type SpecialPrice struct {
	Price     float32   `bson:"price"`
	DateRange DateRange `bson:"date_range"`
}

type DateRange struct {
	Start time.Time `bson:"start"`
	End   time.Time `bson:"end"`
}

func (p PricingType) String() string {
	switch p {
	case PerApartmentUnit:
		return "PerApartmentUnit"
	case PerGuest:
		return "PerGuest"
	default:
		return "Unknown"
	}
}
