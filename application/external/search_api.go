package external

import (
	"context"
	search "github.com/ZMS-DevOps/search-service/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func AddSearchAccommodation(searchClient search.SearchServiceClient, accommodation *search.Accommodation) (*search.AddAccommodationResponse, error) {
	return searchClient.AddAccommodation(context.TODO(), &search.AddAccommodationRequest{Accommodation: accommodation})
}

func EditSearchAccommodation(searchClient search.SearchServiceClient, accommodation *search.Accommodation) (*search.EditAccommodationResponse, error) {
	return searchClient.EditAccommodation(context.TODO(), &search.EditAccommodationRequest{Accommodation: accommodation})
}

func DeleteSearchAccommodation(searchClient search.SearchServiceClient, id primitive.ObjectID) (*search.DeleteAccommodationResponse, error) {
	return searchClient.DeleteAccommodation(context.TODO(), &search.DeleteAccommodationRequest{AccommodationId: id.Hex()})
}
