package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AccommodationStore interface {
	Get(id primitive.ObjectID) (*Accommodation, error)
	GetAll() ([]*Accommodation, error)
	Insert(accommodation *Accommodation) error
	DeleteAll()
	Delete(id primitive.ObjectID) error
	Update(id primitive.ObjectID, accommodation *Accommodation) error
	UpdateDefaultPrice(id primitive.ObjectID, price *float64) error
	UpdateSpecialPrice(id primitive.ObjectID, price *float64, Start *time.Time, End *time.Time) error
	UpdateTypeOfPayment(id primitive.ObjectID, typeOfPayment *string) error
}
