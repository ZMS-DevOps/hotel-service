package application_test

import (
	booking "github.com/ZMS-DevOps/booking-service/proto"
	application2 "github.com/ZMS-DevOps/hotel-service/application"
	"github.com/ZMS-DevOps/hotel-service/application/test"
	"github.com/ZMS-DevOps/hotel-service/domain"
	"github.com/ZMS-DevOps/hotel-service/infrastructure/dto"
	search "github.com/ZMS-DevOps/search-service/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestAccommodationService_Get(t *testing.T) {
	mockStore := new(application.MockAccommodationStore)
	mockBookingClient := new(application.MockBookingServiceClient)
	mockSearchClient := new(application.MockSearchServiceClient)
	lokiMock := new(application.LokiMock)
	spanMock := new(application.SpanMock)
	service := application2.NewAccommodationService(mockStore, mockBookingClient, mockSearchClient, lokiMock)

	accommodationID := primitive.NewObjectID()
	expectedAccommodation := &domain.Accommodation{Id: accommodationID}

	lokiMock.On("Infof", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	spanMock.On("AddEvent", mock.Anything, mock.Anything).Return(nil)
	mockStore.On("Get", accommodationID).Return(expectedAccommodation, nil)

	result, err := service.Get(accommodationID, spanMock, lokiMock)

	assert.NoError(t, err)
	assert.Equal(t, expectedAccommodation, result)
	mockStore.AssertExpectations(t)
}

func TestAccommodationService_GetAll(t *testing.T) {
	mockStore := new(application.MockAccommodationStore)
	mockBookingClient := new(application.MockBookingServiceClient)
	mockSearchClient := new(application.MockSearchServiceClient)
	lokiMock := new(application.LokiMock)
	spanMock := new(application.SpanMock)
	service := application2.NewAccommodationService(mockStore, mockBookingClient, mockSearchClient, lokiMock)

	expectedAccommodations := []*domain.Accommodation{
		{Id: primitive.NewObjectID()},
		{Id: primitive.NewObjectID()},
	}
	lokiMock.On("Infof", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	spanMock.On("AddEvent", mock.Anything, mock.Anything).Return(nil)
	mockStore.On("GetAll").Return(expectedAccommodations, nil)

	result, err := service.GetAll(spanMock, lokiMock)

	assert.NoError(t, err)
	assert.Equal(t, expectedAccommodations, result)
	mockStore.AssertExpectations(t)
}

func TestAccommodationService_Add(t *testing.T) {
	mockStore := new(application.MockAccommodationStore)
	mockBookingClient := new(application.MockBookingServiceClient)
	mockSearchClient := new(application.MockSearchServiceClient)
	lokiMock := new(application.LokiMock)
	spanMock := new(application.SpanMock)
	service := application2.NewAccommodationService(mockStore, mockBookingClient, mockSearchClient, lokiMock)

	accommodationID := primitive.NewObjectID()
	newAccommodation := &domain.Accommodation{Id: accommodationID, Name: "Test Accommodation"}

	lokiMock.On("Infof", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	spanMock.On("AddEvent", mock.Anything, mock.Anything).Return(nil)
	mockStore.On("Insert", newAccommodation).Return(nil)
	mockBookingClient.On("AddUnavailability", mock.Anything, mock.Anything, mock.Anything).Return(&booking.AddUnavailabilityResponse{}, nil)
	mockSearchClient.On("AddAccommodation", mock.Anything, mock.Anything, mock.Anything).Return(&search.AddAccommodationResponse{}, nil)

	err := service.Add(newAccommodation, spanMock, lokiMock)

	assert.NoError(t, err)
	mockStore.AssertExpectations(t)
	mockSearchClient.AssertExpectations(t)
}

func TestAccommodationService_Update(t *testing.T) {
	mockStore := new(application.MockAccommodationStore)
	mockBookingClient := new(application.MockBookingServiceClient)
	mockSearchClient := new(application.MockSearchServiceClient)
	lokiMock := new(application.LokiMock)
	spanMock := new(application.SpanMock)
	service := application2.NewAccommodationService(mockStore, mockBookingClient, mockSearchClient, lokiMock)

	accommodationID := primitive.NewObjectID()
	updatedAccommodation := &domain.Accommodation{Id: accommodationID, Name: "Updated Accommodation"}

	lokiMock.On("Infof", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	spanMock.On("AddEvent", mock.Anything, mock.Anything).Return(nil)
	mockStore.On("Get", accommodationID).Return(updatedAccommodation, nil)
	mockStore.On("Update", accommodationID, updatedAccommodation).Return(nil)
	mockBookingClient.On("EditAccommodation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&booking.EditAccommodationResponse{}, nil)
	mockSearchClient.On("EditAccommodation", mock.Anything, mock.Anything, mock.Anything).Return(&search.EditAccommodationResponse{}, nil)

	err := service.Update(accommodationID, updatedAccommodation, spanMock, lokiMock)

	assert.NoError(t, err)
	mockStore.AssertExpectations(t)
	mockSearchClient.AssertExpectations(t)
}

func TestAccommodationService_Delete(t *testing.T) {
	mockStore := new(application.MockAccommodationStore)
	mockBookingClient := new(application.MockBookingServiceClient)
	mockSearchClient := new(application.MockSearchServiceClient)
	lokiMock := new(application.LokiMock)
	spanMock := new(application.SpanMock)
	service := application2.NewAccommodationService(mockStore, mockBookingClient, mockSearchClient, lokiMock)

	accommodationID := primitive.NewObjectID()

	lokiMock.On("Infof", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	spanMock.On("AddEvent", mock.Anything, mock.Anything).Return(nil)
	mockStore.On("Delete", accommodationID).Return(nil)
	mockBookingClient.On("CheckAccommodationHasReservation", mock.Anything, mock.Anything, mock.Anything).Return(&booking.CheckAccommodationHasReservationResponse{Success: true}, nil)
	mockSearchClient.On("DeleteAccommodation", mock.Anything, mock.Anything, mock.Anything).Return(&search.DeleteAccommodationResponse{}, nil)

	err := service.Delete(accommodationID, spanMock, lokiMock)

	assert.NoError(t, err)
	mockStore.AssertExpectations(t)
	mockSearchClient.AssertExpectations(t)
}

func TestAccommodationService_Delete_Fail(t *testing.T) {
	mockStore := new(application.MockAccommodationStore)
	mockBookingClient := new(application.MockBookingServiceClient)
	mockSearchClient := new(application.MockSearchServiceClient)
	lokiMock := new(application.LokiMock)
	spanMock := new(application.SpanMock)
	service := application2.NewAccommodationService(mockStore, mockBookingClient, mockSearchClient, lokiMock)

	accommodationID := primitive.NewObjectID()

	lokiMock.On("Infof", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	spanMock.On("AddEvent", mock.Anything, mock.Anything).Return(nil)
	mockStore.On("Delete", accommodationID).Return(nil)
	mockBookingClient.On("CheckAccommodationHasReservation", mock.Anything, mock.Anything, mock.Anything).Return(&booking.CheckAccommodationHasReservationResponse{Success: false}, nil)
	mockSearchClient.On("DeleteAccommodation", mock.Anything, mock.Anything, mock.Anything).Return(&search.DeleteAccommodationResponse{}, nil)

	err := service.Delete(accommodationID, spanMock, lokiMock)

	assert.Error(t, err)
}

func TestAccommodationService_UpdateDefaultPrice(t *testing.T) {
	mockStore := new(application.MockAccommodationStore)
	mockBookingClient := new(application.MockBookingServiceClient)
	mockSearchClient := new(application.MockSearchServiceClient)
	lokiMock := new(application.LokiMock)
	spanMock := new(application.SpanMock)
	service := application2.NewAccommodationService(mockStore, mockBookingClient, mockSearchClient, lokiMock)

	accommodationID := primitive.NewObjectID()

	lokiMock.On("Infof", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	spanMock.On("AddEvent", mock.Anything, mock.Anything).Return(nil)
	mockStore.On("Get", accommodationID).Return(&domain.Accommodation{}, nil)
	mockStore.On("UpdateDefaultPrice", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockSearchClient.On("EditAccommodation", mock.Anything, mock.Anything, mock.Anything).Return(&search.EditAccommodationResponse{}, nil)

	price := float32(500)
	service.UpdatePrice(accommodationID, dto.UpdatePriceDto{Price: &price}, spanMock, lokiMock)

	mockStore.AssertCalled(t, "UpdateDefaultPrice", mock.Anything, mock.Anything, mock.Anything)
}
