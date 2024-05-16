package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mmmajder/zms-devops-hotel-service/application"
	"github.com/mmmajder/zms-devops-hotel-service/domain"
	"github.com/mmmajder/zms-devops-hotel-service/infrastructure/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type AccommodationHandler struct {
	service *application.AccommodationService
}

type AccommodationsResponse struct {
	Accommodations []*domain.Accommodation `json:"accommodations"`
}

type AccommodationResponse struct {
	Accommodation *domain.Accommodation `json:"accommodation"`
}

type HealthCheckResponse struct {
	Size string `json:"size"`
}

func NewAccommodationHandler(service *application.AccommodationService) *AccommodationHandler {
	server := &AccommodationHandler{
		service: service,
	}
	return server
}

func (handler *AccommodationHandler) Init(router *mux.Router) {
	router.HandleFunc(`/hotel/accommodation`, handler.GetAll).Methods("GET")
	router.HandleFunc("/hotel/accommodation/{id}", handler.GetById).Methods("GET")
	router.HandleFunc("/hotel/accommodation", handler.Add).Methods("POST")
	router.HandleFunc("/hotel/accommodation/{id}", handler.Update).Methods("PUT")
	router.HandleFunc("/hotel/accommodation/{id}", handler.Delete).Methods("DELETE")
	router.HandleFunc("/hotel/accommodation/price/{id}", handler.UpdatePrice).Methods("PUT")
	router.HandleFunc("/hotel/health", handler.GetHealthCheck).Methods("GET")
}

func (handler *AccommodationHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accommodationId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid accommodation ID")
		return
	}

	var updatedAccommodationDto dto.AccommodationDto
	if err := json.NewDecoder(r.Body).Decode(&updatedAccommodationDto); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate the updated accommodation DTO
	if err := dto.ValidateAccommodationDto(&updatedAccommodationDto); err != nil {
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}
	updatedAccommodation := dto.MapAccommodation(&updatedAccommodationDto)
	updatedAccommodation.Id = accommodationId

	if err := handler.service.Update(accommodationId, updatedAccommodation); err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to update accommodation")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *AccommodationHandler) UpdatePrice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accommodationId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid accommodation ID")
		return
	}

	var updatePriceDto dto.UpdatePriceDto
	if err := json.NewDecoder(r.Body).Decode(&updatePriceDto); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate update price dto
	if err := dto.ValidateUpdatePriceDto(&updatePriceDto); err != nil {
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.service.UpdatePrice(accommodationId, updatePriceDto); err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to update accommodation")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *AccommodationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accommodationId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		handleError(w, http.StatusInternalServerError, "Invalid accommodation ID")
		return
	}

	if err := handler.service.Delete(accommodationId); err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to delete accommodation")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *AccommodationHandler) GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	response := HealthCheckResponse{
		Size: "Hotel SERVICE OK",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (handler *AccommodationHandler) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accommodationId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	accommodation, err := handler.service.Get(accommodationId)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response := AccommodationResponse{
		Accommodation: accommodation,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (handler *AccommodationHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	accommodation, err := handler.service.GetAll()

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response := AccommodationsResponse{
		Accommodations: accommodation,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (handler *AccommodationHandler) Add(w http.ResponseWriter, r *http.Request) {
	var createAccommodationDto dto.AccommodationDto
	if err := json.NewDecoder(r.Body).Decode(&createAccommodationDto); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate the updated accommodation DTO
	if err := dto.ValidateAccommodationDto(&createAccommodationDto); err != nil {
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	newAccommodation := dto.MapAccommodation(&createAccommodationDto)

	if err := handler.service.Add(newAccommodation); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handleError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, message)
}