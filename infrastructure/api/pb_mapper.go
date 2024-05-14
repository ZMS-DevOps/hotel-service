package api

import (
	"github.com/mmmajder/zms-devops-hotel-service/domain"
	pb "github.com/mmmajder/zms-devops-hotel-service/proto"
)

func mapHotel(hotel *domain.Accommodation) *pb.Hotel {
	hotelPb := &pb.Hotel{
		Id:   hotel.Id.Hex(),
		Name: hotel.Name,
	}
	return hotelPb
}
