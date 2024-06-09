package external

import (
	"context"
	booking "github.com/ZMS-DevOps/booking-service/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func NewBookingClient(address string) booking.BookingServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to Catalogue service: %v", err)
	}
	return booking.NewBookingServiceClient(conn)
}

func getConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func CreateBookingUnavailability(bookingClient booking.BookingServiceClient, id primitive.ObjectID, reviewReservationRequestAutomatically bool, hostId string, name string) (*booking.AddUnavailabilityResponse, error) {
	return bookingClient.AddUnavailability(
		context.TODO(),
		&booking.AddUnavailabilityRequest{
			Id:                id.Hex(),
			Automatically:     reviewReservationRequestAutomatically,
			HostId:            hostId,
			AccommodationName: name,
		})
}

func UpdateBookingUnavailability(bookingClient booking.BookingServiceClient, id primitive.ObjectID, reviewReservationRequestAutomatically bool, hostId string, name string) (*booking.EditAccommodationResponse, error) {
	return bookingClient.EditAccommodation(
		context.TODO(),
		&booking.EditAccommodationRequest{
			Id:                id.Hex(),
			Automatically:     reviewReservationRequestAutomatically,
			HostId:            hostId,
			AccommodationName: name,
		})
}
