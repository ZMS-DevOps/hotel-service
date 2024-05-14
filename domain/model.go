package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Accommodation struct {
	Id           primitive.ObjectID `bson:"_id"`
	Name         string             `bson:"name"`
	Location     string             `bson:"location"`
	Benefits     []string           `bson:"benefits"`
	Photos       []string           `bson:"photos"`
	GuestNumber  GuestNumber        `bson:"guest_number"`
	DefaultPrice DefaultPrice       `bson:"default_price"`
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
	Price float64     `bson:"price"`
	Type  PricingType `bson:"type"`
}
