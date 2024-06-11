package startup

import (
	"github.com/ZMS-DevOps/hotel-service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var accommodations = []*domain.Accommodation{
	{
		Id:       getObjectId("6643a56c9dea1760db469b7b"),
		Name:     "Wailea Beach Resort",
		Location: "Hawaii",
		HostId:   "04d19820-6340-4c93-84f9-2ffda959a0d9",
		Benefits: []string{"BBQ facilities", "Balcony", "No smoke"},
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
		Location: "Bali",
		HostId:   "57325353-5469-4930-8ec9-35c003e1b967",
		Benefits: []string{"Kitchen", "View", "TV"},
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
	{
		Id:       getObjectId("ff43bdc7240f80f13b5d23e8"),
		Name:     "Luxury Villa",
		Location: "Bali",
		HostId:   "57325353-5469-4930-8ec9-35c003e1b967",
		Benefits: []string{"Free WiFi", "Air condition", "Bathtub"},
		Photos:   []string{"photo31.jpg", "photo32.jpg", "photo33.jpg"},
		GuestNumber: domain.GuestNumber{
			Min: 2,
			Max: 8,
		},
		DefaultPrice: domain.DefaultPrice{
			Price: 1000.00,
			Type:  domain.PerApartmentUnit,
		},
		SpecialPrice:                          []domain.SpecialPrice{{}},
		ReviewReservationRequestAutomatically: true,
	},
	{
		Id:       getObjectId("gg43bdc7240f80f13b5d24e9"),
		Name:     "Balinese Bungalow",
		Location: "Bali",
		HostId:   "57325353-5469-4930-8ec9-35c003e1b967",
		Benefits: []string{"Phone", "Kings bed", "Laundry"},
		Photos:   []string{"photo34.jpg", "photo35.jpg", "photo36.jpg"},
		GuestNumber: domain.GuestNumber{
			Min: 1,
			Max: 4,
		},
		DefaultPrice: domain.DefaultPrice{
			Price: 200.00,
			Type:  domain.PerGuest,
		},
		SpecialPrice:                          []domain.SpecialPrice{{}},
		ReviewReservationRequestAutomatically: true,
	},
	{
		Id:       getObjectId("hh43bdc7240f80f13b5d25e0"),
		Name:     "Cliffside Retreat",
		Location: "Bali",
		HostId:   "04d19820-6340-4c93-84f9-2ffda959a0d9",
		Benefits: []string{"Parking", "Swimming pool", "Fitness center"},
		Photos:   []string{"photo37.jpg", "photo38.jpg", "photo39.jpg"},
		GuestNumber: domain.GuestNumber{
			Min: 2,
			Max: 6,
		},
		DefaultPrice: domain.DefaultPrice{
			Price: 800.00,
			Type:  domain.PerApartmentUnit,
		},
		SpecialPrice:                          []domain.SpecialPrice{{}},
		ReviewReservationRequestAutomatically: false,
	},
	{
		Id:       getObjectId("ii43bdc7240f80f13b5d26e1"),
		Name:     "Ubud Villa",
		Location: "Bali",
		HostId:   "04d19820-6340-4c93-84f9-2ffda959a0d9",
		Benefits: []string{"BBQ facilities", "Balcony", "No smoke"},
		Photos:   []string{"photo40.jpg", "photo41.jpg", "photo42.jpg"},
		GuestNumber: domain.GuestNumber{
			Min: 2,
			Max: 5,
		},
		DefaultPrice: domain.DefaultPrice{
			Price: 650.00,
			Type:  domain.PerApartmentUnit,
		},
		SpecialPrice:                          []domain.SpecialPrice{{}},
		ReviewReservationRequestAutomatically: false,
	},
	{
		Id:       getObjectId("7743a6dc9dea1760db469b8c"),
		Name:     "Mountain Retreat",
		Location: "Switzerland",
		HostId:   "04d19820-6340-4c93-84f9-2ffda959a0d9",
		Benefits: []string{"BBQ facilities", "Balcony", "No smoke"},
		Photos:   []string{"photo7.jpg", "photo8.jpg", "photo9.jpg"},
		GuestNumber: domain.GuestNumber{
			Min: 4,
			Max: 10,
		},
		DefaultPrice: domain.DefaultPrice{
			Price: 1000.00,
			Type:  domain.PerApartmentUnit,
		},
		SpecialPrice:                          []domain.SpecialPrice{{}},
		ReviewReservationRequestAutomatically: false,
	},
	{
		Id:       getObjectId("8843bdc7240f80f13b5d19e7"),
		Name:     "Urban Loft",
		Location: "New York",
		HostId:   "04d19820-6340-4c93-84f9-2ffda959a0d9",
		Benefits: []string{"Phone", "Kings bed", "Laundry"},
		Photos:   []string{"photo10.jpg", "photo11.jpg", "photo12.jpg"},
		GuestNumber: domain.GuestNumber{
			Min: 1,
			Max: 3,
		},
		DefaultPrice: domain.DefaultPrice{
			Price: 200.00,
			Type:  domain.PerGuest,
		},
		ReviewReservationRequestAutomatically: true,
	},
	{
		Id:       getObjectId("jj43bdc7240f80f13b5d27e2"),
		Name:     "Central Park Apartment",
		Location: "New York",
		HostId:   "04d19820-6340-4c93-84f9-2ffda959a0d9",
		Benefits: []string{"Parking", "Swimming pool", "Fitness center"},
		Photos:   []string{"photo43.jpg", "photo44.jpg", "photo45.jpg"},
		GuestNumber: domain.GuestNumber{
			Min: 2,
			Max: 5,
		},
		DefaultPrice: domain.DefaultPrice{
			Price: 350.00,
			Type:  domain.PerApartmentUnit,
		},
		ReviewReservationRequestAutomatically: false,
	},
	{
		Id:       getObjectId("kk43bdc7240f80f13b5d28e3"),
		Name:     "Brooklyn Brownstone",
		Location: "New York",
		HostId:   "04d19820-6340-4c93-84f9-2ffda959a0d9",
		Benefits: []string{"Free WiFi", "Air condition", "Bathtub"},
		Photos:   []string{"photo46.jpg", "photo47.jpg", "photo48.jpg"},
		GuestNumber: domain.GuestNumber{
			Min: 1,
			Max: 6,
		},
		DefaultPrice: domain.DefaultPrice{
			Price: 300.00,
			Type:  domain.PerGuest,
		},
		ReviewReservationRequestAutomatically: true,
	},
	{
		Id:       getObjectId("ll43bdc7240f80f13b5d29e4"),
		Name:     "Times Square Studio",
		Location: "New York",
		HostId:   "04d19820-6340-4c93-84f9-2ffda959a0d9",
		Benefits: []string{"Parking", "Swimming pool", "Fitness center"},
		Photos:   []string{"photo49.jpg", "photo50.jpg", "photo51.jpg"},
		GuestNumber: domain.GuestNumber{
			Min: 1,
			Max: 2,
		},
		DefaultPrice: domain.DefaultPrice{
			Price: 250.00,
			Type:  domain.PerGuest,
		},
		ReviewReservationRequestAutomatically: true,
	},
	{
		Id:       getObjectId("9943a7dc9dea1760db469b9d"),
		Name:     "Beachfront Villa",
		Location: "Maldives",
		HostId:   "04d19820-6340-4c93-84f9-2ffda959a0d9",
		Benefits: []string{"Free WiFi", "Air condition", "Bathtub"},
		Photos:   []string{"photo13.jpg", "photo14.jpg", "photo15.jpg"},
		GuestNumber: domain.GuestNumber{
			Min: 2,
			Max: 8,
		},
		DefaultPrice: domain.DefaultPrice{
			Price: 1500.00,
			Type:  domain.PerApartmentUnit,
		},
		SpecialPrice:                          []domain.SpecialPrice{{}},
		ReviewReservationRequestAutomatically: false,
	},
	{
		Id:       getObjectId("aa43bdc7240f80f13b5d20d7"),
		Name:     "Country Farmhouse",
		Location: "France",
		HostId:   "04d19820-6340-4c93-84f9-2ffda959a0d9",
		Benefits: []string{"Parking", "Swimming pool", "Fitness center"},
		Photos:   []string{"photo16.jpg", "photo17.jpg", "photo18.jpg"},
		GuestNumber: domain.GuestNumber{
			Min: 2,
			Max: 5,
		},
		DefaultPrice: domain.DefaultPrice{
			Price: 300.00,
			Type:  domain.PerGuest,
		},
		ReviewReservationRequestAutomatically: true,
	},
	{
		Id:       getObjectId("bb43a8dc9dea1760db469c1e"),
		Name:     "Desert Oasis",
		Location: "Dubai",
		HostId:   "04d19820-6340-4c93-84f9-2ffda959a0d9",
		Benefits: []string{"Parking", "Swimming pool", "Fitness center"},
		Photos:   []string{"photo19.jpg", "photo20.jpg", "photo21.jpg"},
		GuestNumber: domain.GuestNumber{
			Min: 2,
			Max: 4,
		},
		DefaultPrice: domain.DefaultPrice{
			Price: 750.00,
			Type:  domain.PerGuest,
		},
		SpecialPrice:                          []domain.SpecialPrice{{}},
		ReviewReservationRequestAutomatically: false,
	},
	{
		Id:       getObjectId("cc43bdc7240f80f13b5d21e7"),
		Name:     "Tropical Paradise",
		Location: "Thailand",
		HostId:   "04d19820-6340-4c93-84f9-2ffda959a0d9",
		Benefits: []string{"Parking", "Swimming pool", "Fitness center"},
		Photos:   []string{"photo22.jpg", "photo23.jpg", "photo24.jpg"},
		GuestNumber: domain.GuestNumber{
			Min: 1,
			Max: 6,
		},
		DefaultPrice: domain.DefaultPrice{
			Price: 600.00,
			Type:  domain.PerApartmentUnit,
		},
		ReviewReservationRequestAutomatically: true,
	},
	{
		Id:       getObjectId("dd43a9dc9dea1760db469c2f"),
		Name:     "Safari Lodge",
		Location: "Kenya",
		HostId:   "04d19820-6340-4c93-84f9-2ffda959a0d9",
		Benefits: []string{"Parking", "Swimming pool", "Fitness center"},
		Photos:   []string{"photo25.jpg", "photo26.jpg", "photo27.jpg"},
		GuestNumber: domain.GuestNumber{
			Min: 2,
			Max: 12,
		},
		DefaultPrice: domain.DefaultPrice{
			Price: 1200.00,
			Type:  domain.PerApartmentUnit,
		},
		SpecialPrice:                          []domain.SpecialPrice{{}},
		ReviewReservationRequestAutomatically: false,
	},
	{
		Id:       getObjectId("ee43bdc7240f80f13b5d22e7"),
		Name:     "Lake House",
		Location: "Canada",
		HostId:   "04d19820-6340-4c93-84f9-2ffda959a0d9",
		Benefits: []string{"Parking", "Swimming pool", "Fitness center"},
		Photos:   []string{"photo28.jpg", "photo29.jpg", "photo30.jpg"},
		GuestNumber: domain.GuestNumber{
			Min: 1,
			Max: 6,
		},
		DefaultPrice: domain.DefaultPrice{
			Price: 400.00,
			Type:  domain.PerApartmentUnit,
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
