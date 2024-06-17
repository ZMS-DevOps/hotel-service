package api

import (
	"encoding/json"
	"fmt"
	"github.com/ZMS-DevOps/hotel-service/application"
	"github.com/ZMS-DevOps/hotel-service/infrastructure/dto"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	router.HandleFunc("/accommodation/images", handler.GetImagesForAccommodations).Methods("POST")
}

func (handler *AccommodationHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accommodationId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid accommodation ID")
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Failed to parse form data")
		return
	}
	var createAccommodationDto dto.AccommodationDto
	jsonData := r.FormValue("json")
	if err := json.Unmarshal([]byte(jsonData), &createAccommodationDto); err != nil {
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := dto.ValidateAccommodationDto(&createAccommodationDto); err != nil {
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	photos, err := handlePhotoUploads(r, w, createAccommodationDto.HostId)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Failed to upload photos")
		return
	}

	updatedAccommodation := dto.MapAccommodation(&createAccommodationDto)
	updatedAccommodation.Photos = photos

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
		if err.Error() == "accommodation could not be deleted" {
			handleError(w, http.StatusPreconditionFailed, "accommodation could not be deleted")
		} else {
			handleError(w, http.StatusInternalServerError, "Failed to delete accommodation")
		}

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

func (handler *AccommodationHandler) GetImagesForAccommodations(w http.ResponseWriter, r *http.Request) {
	var accommodationIds []dto.GetImagesRequest
	if err := json.NewDecoder(r.Body).Decode(&accommodationIds); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	images, err := handler.service.GetImages(accommodationIds)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse, err := json.Marshal(images)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (handler *AccommodationHandler) Add(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Failed to parse form data")
		return
	}
	var createAccommodationDto dto.AccommodationDto
	jsonData := r.FormValue("json")
	if err := json.Unmarshal([]byte(jsonData), &createAccommodationDto); err != nil {
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := dto.ValidateAccommodationDto(&createAccommodationDto); err != nil {
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	photos, err := handlePhotoUploads(r, w, createAccommodationDto.HostId)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Failed to upload photos")
		return
	}

	newAccommodation := dto.MapAccommodation(&createAccommodationDto)
	newAccommodation.Photos = photos

	if err := handler.service.Add(newAccommodation); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handlePhotoUploads(r *http.Request, w http.ResponseWriter, hostId string) ([]string, error) {
	var photos []string
	for _, fileHeader := range r.MultipartForm.File["photos"] {
		file, err := fileHeader.Open()
		if err != nil {
			handleError(w, http.StatusInternalServerError, "Failed to open file")
			return nil, err
		}
		defer file.Close()

		uploadDir := "./uploads"
		if err := os.MkdirAll(uploadDir, 0777); err != nil {
			fmt.Println("Error creating directory:", err)
		}
		dst, err := os.Create(filepath.Join(uploadDir, hostId+"-"+fileHeader.Filename))
		if err != nil {
			fmt.Println("Error creating file:", err)
			handleError(w, http.StatusInternalServerError, "Failed to create file")
			return nil, err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			fmt.Println("Error Copy file:", err)
			handleError(w, http.StatusInternalServerError, "Failed to save file")
			return nil, err
		}

		photos = append(photos, uploadDir+"/"+hostId+"-"+fileHeader.Filename)
	}
	return photos, nil
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
	fmt.Fprintf(w, message)
}
