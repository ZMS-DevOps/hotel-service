package startup

import (
	"github.com/ZMS-DevOps/hotel-service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var accommodations = []*domain.Accommodation{
	{
		Id:       getObjectId("6643a56c9dea1760db469b7b"),
		Name:     "Luxury Villa",
		Location: "Tropical Paradise",
		HostId:   getObjectId(""),
		Benefits: []string{"Private Pool", "Ocean View", "Spa Services"},
		Photos:   []string{"photo1.jpg", "photo2.jpg", "photo3.jpg"},
		GuestNumber: domain.GuestNumber{
			Min: 2,
			Max: 6,
		},
		DefaultPrice: domain.DefaultPrice{
			Price: 500.00,
			Type:  domain.PerApartmentUnit,
		},
		SpecialPrice:                          []domain.SpecialPrice{{}},
		ReviewReservationRequestAutomatically: false,
	},
	{
		Id:       getObjectId("6643bdc7240f80f13b5d18d7"),
		Name:     "Cozy Cottage",
		Location: "Mountain Retreat",
		HostId:   getObjectId(""),
		Benefits: []string{"Fireplace", "Scenic Views", "Hiking Trails"},
		Photos:   []string{"photo4.jpg", "photo5.jpg", "photo6.jpg"},
		GuestNumber: domain.GuestNumber{
			Min: 1,
			Max: 4,
		},
		DefaultPrice: domain.DefaultPrice{
			Price: 75.00,
			Type:  domain.PerGuest,
		},
		ReviewReservationRequestAutomatically: true,
	},
}

func getObjectId(id string) primitive.ObjectID {
	if objectId, err := primitive.ObjectIDFromHex(id); err == nil {
		return objectId
	}
	return primitive.NewObjectID()
}
