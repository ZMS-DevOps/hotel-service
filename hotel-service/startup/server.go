package startup

import (
	"fmt"
	booking "github.com/ZMS-DevOps/booking-service/proto"
	"github.com/ZMS-DevOps/hotel-service/application/external"
	search "github.com/ZMS-DevOps/search-service/proto"
	"github.com/gorilla/mux"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/ZMS-DevOps/hotel-service/application"
	"github.com/ZMS-DevOps/hotel-service/domain"
	"github.com/ZMS-DevOps/hotel-service/infrastructure/api"
	"github.com/ZMS-DevOps/hotel-service/infrastructure/persistence"
	"github.com/ZMS-DevOps/hotel-service/startup/config"
	"github.com/afiskon/promtail-client/promtail"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

type Server struct {
	config               *config.Config
	router               *mux.Router
	AccommodationHandler *api.AccommodationHandler
	traceProvider        *sdktrace.TracerProvider
	loki                 promtail.Client
}

func NewServer(config *config.Config, traceProvider *sdktrace.TracerProvider, loki promtail.Client) *Server {
	server := &Server{
		config:        config,
		router:        mux.NewRouter(),
		traceProvider: traceProvider,
		loki:          loki,
	}
	server.AccommodationHandler = server.setupHandlers()
	return server
}

func (server *Server) setupHandlers() *api.AccommodationHandler {
	mongoClient := server.initMongoClient()
	bookingClient := external.NewBookingClient(server.getBookingAddress())
	searchClient := external.NewSearchClient(server.getSearchAddress())
	accommodationStore := server.initAccommodationStore(mongoClient)
	accommodationService := server.initAccommodationService(accommodationStore, bookingClient, searchClient)
	accommodationHandler := server.initAccommodationHandler(accommodationService)
	accommodationHandler.Init(server.router)
	return accommodationHandler
}

func (server *Server) Start() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), server.router))
}

func (server *Server) getBookingAddress() string {
	return fmt.Sprintf("%s:%s", server.config.BookingHost, server.config.BookingPort)
}

func (server *Server) getSearchAddress() string {
	return fmt.Sprintf("%s:%s", server.config.SearchHost, server.config.SearchPort)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.HotelDBUsername, server.config.HotelDBPassword, server.config.HotelDBHost, server.config.HotelDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initAccommodationStore(client *mongo.Client) domain.AccommodationStore {
	store := persistence.NewAccommodationMongoDBStore(client)
	store.DeleteAll()
	for _, accommodation := range accommodations {
		_ = store.InsertWithId(accommodation)
	}
	return store
}

func (server *Server) initAccommodationService(store domain.AccommodationStore, bookingClient booking.BookingServiceClient, searchClient search.SearchServiceClient) *application.AccommodationService {
	return application.NewAccommodationService(store, bookingClient, searchClient, server.loki)
}

func (server *Server) initAccommodationHandler(service *application.AccommodationService) *api.AccommodationHandler {
	return api.NewAccommodationHandler(service, server.traceProvider, server.loki)
}
