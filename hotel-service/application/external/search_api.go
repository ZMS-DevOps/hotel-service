package external

import (
	"context"
	"github.com/ZMS-DevOps/hotel-service/util"
	search "github.com/ZMS-DevOps/search-service/proto"
	"github.com/afiskon/promtail-client/promtail"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func NewSearchClient(address string) search.SearchServiceClient {
	conn, err := getSearchConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to Search service: %v", err)
	}
	return search.NewSearchServiceClient(conn)
}

func getSearchConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func AddSearchAccommodation(searchClient search.SearchServiceClient, accommodation *search.Accommodation, span trace.Span, loki promtail.Client) (*search.AddAccommodationResponse, error) {
	util.HttpTraceInfo("Adding accommodation in search service...", span, loki, "Add", "")
	return searchClient.AddAccommodation(context.TODO(), &search.AddAccommodationRequest{Accommodation: accommodation})
}

func EditSearchAccommodation(searchClient search.SearchServiceClient, accommodation *search.Accommodation, span trace.Span, loki promtail.Client) (*search.EditAccommodationResponse, error) {
	util.HttpTraceInfo("Edit accommodation in search service...", span, loki, "Add", "")
	return searchClient.EditAccommodation(context.TODO(), &search.EditAccommodationRequest{Accommodation: accommodation})
}

func DeleteSearchAccommodation(searchClient search.SearchServiceClient, id primitive.ObjectID, span trace.Span, loki promtail.Client) (*search.DeleteAccommodationResponse, error) {
	util.HttpTraceInfo("Delete accommodation in search service...", span, loki, "Add", "")
	return searchClient.DeleteAccommodation(context.TODO(), &search.DeleteAccommodationRequest{AccommodationId: id.Hex()})
}
