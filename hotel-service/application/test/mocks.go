package application

import (
	"context"
	booking "github.com/ZMS-DevOps/booking-service/proto"
	"github.com/ZMS-DevOps/hotel-service/domain"
	search "github.com/ZMS-DevOps/search-service/proto"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

// Mock dependencies
type MockAccommodationStore struct {
	mock.Mock
}

func (m *MockAccommodationStore) Get(id primitive.ObjectID) (*domain.Accommodation, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Accommodation), args.Error(1)
}

func (m *MockAccommodationStore) GetAll() ([]*domain.Accommodation, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Accommodation), args.Error(1)
}

func (m *MockAccommodationStore) GetByHostId(hostId string) ([]*domain.Accommodation, error) {
	args := m.Called(hostId)
	return args.Get(0).([]*domain.Accommodation), args.Error(1)
}

func (m *MockAccommodationStore) Insert(accommodation *domain.Accommodation) error {
	args := m.Called(accommodation)
	return args.Error(0)
}

func (m *MockAccommodationStore) InsertWithId(accommodation *domain.Accommodation) error {
	args := m.Called(accommodation)
	return args.Error(0)
}

func (m *MockAccommodationStore) DeleteAll() {
	m.Called()
}

func (m *MockAccommodationStore) Delete(id primitive.ObjectID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAccommodationStore) Update(id primitive.ObjectID, accommodation *domain.Accommodation) error {
	args := m.Called(id, accommodation)
	return args.Error(0)
}

func (m *MockAccommodationStore) UpdateDefaultPrice(id primitive.ObjectID, price *float32) error {
	args := m.Called(id, price)
	return args.Error(0)
}

func (m *MockAccommodationStore) UpdateSpecialPrice(id primitive.ObjectID, newSpecialPrices []domain.SpecialPrice) error {
	args := m.Called(id, newSpecialPrices)
	return args.Error(0)
}

func (m *MockAccommodationStore) UpdateTypeOfPayment(id primitive.ObjectID, typeOfPayment *string) error {
	args := m.Called(id, typeOfPayment)
	return args.Error(0)
}

func (m *MockAccommodationStore) GetSpecialPrices(id primitive.ObjectID) ([]domain.SpecialPrice, error) {
	args := m.Called(id)
	return args.Get(0).([]domain.SpecialPrice), args.Error(1)
}

func (m *MockAccommodationStore) DeleteByHostId(hostId string) error {
	args := m.Called(hostId)
	return args.Error(0)
}

type MockBookingServiceClient struct {
	mock.Mock
}

func (m *MockBookingServiceClient) AddUnavailability(ctx context.Context, in *booking.AddUnavailabilityRequest, opts ...grpc.CallOption) (*booking.AddUnavailabilityResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*booking.AddUnavailabilityResponse), args.Error(1)
}

func (m *MockBookingServiceClient) EditAccommodation(ctx context.Context, in *booking.EditAccommodationRequest, opts ...grpc.CallOption) (*booking.EditAccommodationResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*booking.EditAccommodationResponse), args.Error(1)
}

func (m *MockBookingServiceClient) FilterAvailableAccommodation(ctx context.Context, in *booking.FilterAvailableAccommodationRequest, opts ...grpc.CallOption) (*booking.FilterAvailableAccommodationResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*booking.FilterAvailableAccommodationResponse), args.Error(1)
}

func (m *MockBookingServiceClient) CheckDeleteHost(ctx context.Context, in *booking.CheckDeleteHostRequest, opts ...grpc.CallOption) (*booking.CheckDeleteHostResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*booking.CheckDeleteHostResponse), args.Error(1)
}

func (m *MockBookingServiceClient) CheckDeleteClient(ctx context.Context, in *booking.CheckDeleteClientRequest, opts ...grpc.CallOption) (*booking.CheckDeleteClientResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*booking.CheckDeleteClientResponse), args.Error(1)
}

func (m *MockBookingServiceClient) CheckGuestHasReservationForHost(ctx context.Context, in *booking.CheckGuestHasReservationForHostRequest, opts ...grpc.CallOption) (*booking.CheckGuestHasReservationForHostResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*booking.CheckGuestHasReservationForHostResponse), args.Error(1)
}

func (m *MockBookingServiceClient) CheckGuestHasReservationForAccommodation(ctx context.Context, in *booking.CheckGuestHasReservationForAccommodationRequest, opts ...grpc.CallOption) (*booking.CheckGuestHasReservationForAccommodationResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*booking.CheckGuestHasReservationForAccommodationResponse), args.Error(1)
}

func (m *MockBookingServiceClient) CheckAccommodationHasReservation(ctx context.Context, in *booking.CheckAccommodationHasReservationRequest, opts ...grpc.CallOption) (*booking.CheckAccommodationHasReservationResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*booking.CheckAccommodationHasReservationResponse), args.Error(1)
}

type MockSearchServiceClient struct {
	mock.Mock
}

func (m *MockSearchServiceClient) AddAccommodation(ctx context.Context, in *search.AddAccommodationRequest, opts ...grpc.CallOption) (*search.AddAccommodationResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*search.AddAccommodationResponse), args.Error(1)
}

func (m *MockSearchServiceClient) EditAccommodation(ctx context.Context, in *search.EditAccommodationRequest, opts ...grpc.CallOption) (*search.EditAccommodationResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*search.EditAccommodationResponse), args.Error(1)
}

func (m *MockSearchServiceClient) DeleteAccommodation(ctx context.Context, in *search.DeleteAccommodationRequest, opts ...grpc.CallOption) (*search.DeleteAccommodationResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*search.DeleteAccommodationResponse), args.Error(1)
}
