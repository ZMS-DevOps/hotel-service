package api

import (
	"context"
	"github.com/ZMS-DevOps/hotel-service/application"
	pb "github.com/ZMS-DevOps/hotel-service/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	pb.UnimplementedHotelServiceServer
	service *application.UserService
}

func NewHotelHandler(service *application.UserService) *HotelHandler {
	return &HotelHandler{
		service: service,
	}
}

func (handler *HotelHandler) AddUser(ctx context.Context, request *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	id := request.Id
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	if err := handler.service.Add(userId); err != nil {
		return nil, err
	}
	return &pb.AddUserResponse{}, nil
}
