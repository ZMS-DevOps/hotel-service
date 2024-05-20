package external

import (
	//"context"
	//booking "github.com/ZMS-DevOps/hotel-service/proto"

	"context"
	booking "github.com/ZMS-DevOps/booking-service/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	//"github.com/mmmajder/devops-search-service/proto"
	//"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials/insecure"
	//"log"
)

func NewBookingClient(address string) booking.BookingServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to Catalogue service: %v", err)
	}
	//defer conn.Close()
	return booking.NewBookingServiceClient(conn)
}

func getConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func CreateBookingUnavailability(bookingClient booking.BookingServiceClient, id primitive.ObjectID) (*booking.AddUnavailabilityResponse, error) {
	return bookingClient.AddUnavailability(context.TODO(), &booking.AddUnavailabilityRequest{Id: id.Hex()})
}
