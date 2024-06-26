package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccommodationStore interface {
	Get(id primitive.ObjectID) (*Accommodation, error)
	GetAll() ([]*Accommodation, error)
	GetByHostId(ownerId string) ([]*Accommodation, error)
	Insert(accommodation *Accommodation) error
	InsertWithId(accommodation *Accommodation) error
	DeleteAll()
	Delete(id primitive.ObjectID) error
	Update(id primitive.ObjectID, accommodation *Accommodation) error
	UpdateDefaultPrice(id primitive.ObjectID, price *float32) error
	UpdateSpecialPrice(id primitive.ObjectID, newSpecialPrices []SpecialPrice) error
	UpdateTypeOfPayment(id primitive.ObjectID, typeOfPayment *string) error
	GetSpecialPrices(id primitive.ObjectID) ([]SpecialPrice, error)
	DeleteByHostId(hostId string) error
}
