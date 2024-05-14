package api

import (
	"context"
	"github.com/mmmajder/zms-devops-hotel-service/application"
	pb "github.com/mmmajder/zms-devops-hotel-service/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingHandler struct {
	pb.UnimplementedBookingServiceServer
	service *application.AccommodationService
}

func NewBookingHandler(service *application.AccommodationService) *BookingHandler {
	return &BookingHandler{
		service: service,
	}
}

func (handler *BookingHandler) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	Hotel, err := handler.service.Get(objectId)
	if err != nil {
		return nil, err
	}
	HotelPb := mapHotel(Hotel)
	response := &pb.GetResponse{
		Hotel: HotelPb,
	}
	return response, nil
}

func (handler *BookingHandler) GetAll(ctx context.Context, request *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	Hotels, err := handler.service.GetAll()
	if err != nil {
		return nil, err
	}
	response := &pb.GetAllResponse{
		Hotels: []*pb.Hotel{},
	}
	for _, Hotel := range Hotels {
		current := mapHotel(Hotel)
		response.Hotels = append(response.Hotels, current)
	}
	return response, nil
}
