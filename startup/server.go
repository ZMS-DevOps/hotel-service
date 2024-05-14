package startup

import (
	"fmt"
	"github.com/gorilla/mux"
	//"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/mmmajder/zms-devops-hotel-service/application"
	"github.com/mmmajder/zms-devops-hotel-service/domain"
	"github.com/mmmajder/zms-devops-hotel-service/infrastructure/api"
	"github.com/mmmajder/zms-devops-hotel-service/infrastructure/persistence"
	hotelPb "github.com/mmmajder/zms-devops-hotel-service/proto"
	"github.com/mmmajder/zms-devops-hotel-service/startup/config"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

type Server struct {
	config *config.Config
	//mux    *runtime.ServeMux
	router *mux.Router
}

func NewServer(config *config.Config) *Server {
	server := &Server{
		config: config,
		router: mux.NewRouter(),
		//mux:    runtime.NewServeMux(),
	}
	return server
}

func (server *Server) Start() {
	mongoClient := server.initMongoClient()
	accommodationStore := server.initAccommodationStore(mongoClient)
	accommodationService := server.initAccommodationService(accommodationStore)
	accommodationHandler := server.initAccommodationHandler(accommodationService)
	accommodationHandler.Init(server.router)
	bookingHandler := server.initBookingHandler(accommodationService)
	go server.startGrpcServer(bookingHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), server.router))
}

func (server *Server) startGrpcServer(bookingHandler *api.BookingHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	hotelPb.RegisterBookingServiceServer(grpcServer, bookingHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
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
		err := store.Insert(accommodation)
		if err != nil {
			log.Fatal(err)
		}
	}
	return store
}

func (server *Server) initAccommodationService(store domain.AccommodationStore) *application.AccommodationService {
	return application.NewAccommodationService(store)
}

func (server *Server) initAccommodationHandler(service *application.AccommodationService) *api.AccommodationHandler {
	return api.NewAccommodationHandler(service)
}

func (server *Server) initBookingHandler(service *application.AccommodationService) *api.BookingHandler {
	return api.NewBookingHandler(service)
}
