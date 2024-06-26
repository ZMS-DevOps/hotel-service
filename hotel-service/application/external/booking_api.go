package external

import (
	"context"
	booking "github.com/ZMS-DevOps/booking-service/proto"
	"github.com/ZMS-DevOps/hotel-service/util"
	"github.com/afiskon/promtail-client/promtail"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/trace"
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

func CreateBookingUnavailability(bookingClient booking.BookingServiceClient, id primitive.ObjectID, reviewReservationRequestAutomatically bool, hostId string, name string, span trace.Span, loki promtail.Client) (*booking.AddUnavailabilityResponse, error) {
	util.HttpTraceInfo("Adding unavailability in booking service...", span, loki, "CreateBookingUnavailability", "")
	return bookingClient.AddUnavailability(
		context.TODO(),
		&booking.AddUnavailabilityRequest{
			Id:                id.Hex(),
			Automatically:     reviewReservationRequestAutomatically,
			HostId:            hostId,
			AccommodationName: name,
		})
}

func UpdateBookingUnavailability(bookingClient booking.BookingServiceClient, id primitive.ObjectID, reviewReservationRequestAutomatically bool, hostId string, name string, span trace.Span, loki promtail.Client) (*booking.EditAccommodationResponse, error) {
	util.HttpTraceInfo("Update unavailability in booking service...", span, loki, "UpdateBookingUnavailability", "")
	return bookingClient.EditAccommodation(
		context.TODO(),
		&booking.EditAccommodationRequest{
			Id:                id.Hex(),
			Automatically:     reviewReservationRequestAutomatically,
			HostId:            hostId,
			AccommodationName: name,
		})
}

func CheckAccommodationHasReservation(bookingClient booking.BookingServiceClient, accommodationId primitive.ObjectID, span trace.Span, loki promtail.Client) (*booking.CheckAccommodationHasReservationResponse, error) {
	util.HttpTraceInfo("Check if accommodation has reservation in booking service...", span, loki, "CheckAccommodationHasReservation", "")
	return bookingClient.CheckAccommodationHasReservation(
		context.TODO(),
		&booking.CheckAccommodationHasReservationRequest{
			AccommodationId: accommodationId.Hex(),
		})
}
