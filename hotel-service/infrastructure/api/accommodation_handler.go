package api

import (
	"encoding/json"
	"github.com/ZMS-DevOps/hotel-service/application"
	"github.com/ZMS-DevOps/hotel-service/infrastructure/dto"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

type AccommodationHandler struct {
	service *application.AccommodationService
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
	router.HandleFunc(`/accommodation`, handler.GetAll).Methods("GET")
	router.HandleFunc("/accommodation/{id}", handler.GetById).Methods("GET")
	router.HandleFunc("/accommodation/host/{id}", handler.GetByHostId).Methods("GET")
	router.HandleFunc("/accommodation", handler.Add).Methods("POST")
	router.HandleFunc("/accommodation/{id}", handler.Update).Methods("PUT")
	router.HandleFunc("/accommodation/{id}", handler.Delete).Methods("DELETE")
	router.HandleFunc("/accommodation/price/{id}", handler.UpdatePrice).Methods("PUT")
	router.HandleFunc("/accommodation/health", handler.GetHealthCheck).Methods("GET")
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

func (handler *AccommodationHandler) GetByHostId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]
	if userId == "" {
		handleError(w, http.StatusBadRequest, "Invalid owner ID")
		return
	}

	accommodations, err := handler.service.GetByHostId(userId)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var response []*dto.AccommodationResponse
	for _, acc := range accommodations {
		accommodationResponse := dto.MapAccommodationResponse(*acc)
		response = append(response, accommodationResponse)
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

	response := dto.MapAccommodationResponse(*accommodation)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (handler *AccommodationHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	accommodations, err := handler.service.GetAll()

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var responses []*dto.AccommodationResponse
	for _, acc := range accommodations {
		response := dto.MapAccommodationResponse(*acc)
		responses = append(responses, response)
	}

	jsonResponse, err := json.Marshal(responses)
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

func (handler *AccommodationHandler) OnDeleteAccommodations(message *kafka.Message) {
	var deleteAccommodationRequest dto.DeleteAccommodationsRequest
	if err := json.Unmarshal(message.Value, &deleteAccommodationRequest); err != nil {
		log.Printf("Error unmarshalling rating change request: %v", err)
	}

	handler.service.OnDeleteAccommodations(deleteAccommodationRequest.HostId)
}

func handleError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}
