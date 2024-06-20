package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ZMS-DevOps/hotel-service/application"
	"github.com/ZMS-DevOps/hotel-service/domain"
	"github.com/ZMS-DevOps/hotel-service/infrastructure/dto"
	"github.com/ZMS-DevOps/hotel-service/util"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type AccommodationHandler struct {
	service       *application.AccommodationService
	traceProvider *sdktrace.TracerProvider
	loki          promtail.Client
}

type HealthCheckResponse struct {
	Size string `json:"size"`
}

func NewAccommodationHandler(service *application.AccommodationService, traceProvider *sdktrace.TracerProvider, loki promtail.Client) *AccommodationHandler {
	server := &AccommodationHandler{
		service:       service,
		traceProvider: traceProvider,
		loki:          loki,
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
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "update-put")
	defer func() { span.End() }()

	vars := mux.Vars(r)
	accommodationId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		util.HttpTraceError(err, "invalid accommodation id", span, handler.loki, "Update", "id: "+vars["id"])
		handleError(w, http.StatusBadRequest, "Invalid accommodation ID")
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		util.HttpTraceError(err, "failed to parse form data", span, handler.loki, "Update", "")
		handleError(w, http.StatusBadRequest, "failed to parse form data")
		return
	}
	var createAccommodationDto dto.AccommodationDto
	jsonData := r.FormValue("json")
	if err := json.Unmarshal([]byte(jsonData), &createAccommodationDto); err != nil {
		util.HttpTraceError(err, "failed to parse json data", span, handler.loki, "Update", "")
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := dto.ValidateAccommodationDto(&createAccommodationDto); err != nil {
		util.HttpTraceError(err, "failed to validate json data", span, handler.loki, "Update", "")
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	photos, err := handlePhotoUploads(r, w, createAccommodationDto.HostId, span, handler)
	if err != nil {
		util.HttpTraceError(err, "failed to upload images", span, handler.loki, "Update", "")
		handleError(w, http.StatusBadRequest, "failed to upload photos")
		return
	}

	updatedAccommodation := dto.MapAccommodation(&createAccommodationDto)
	updatedAccommodation.Photos = photos

	updatedAccommodation.Id = accommodationId

	if err := handler.service.Update(accommodationId, updatedAccommodation, span, handler.loki); err != nil {
		util.HttpTraceError(err, "failed to update accommodation", span, handler.loki, "Update", "")
		handleError(w, http.StatusInternalServerError, "failed to update accommodation")
		return
	}
	util.HttpTraceInfo("Accommodation updated successfully", span, handler.loki, "Update", accommodationId.Hex())
	w.WriteHeader(http.StatusOK)
}

func (handler *AccommodationHandler) UpdatePrice(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "update-price-put")
	defer func() { span.End() }()

	vars := mux.Vars(r)
	accommodationId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		util.HttpTraceError(err, "invalid accommodation id", span, handler.loki, "UpdatePrice", "")
		handleError(w, http.StatusBadRequest, "Invalid accommodation ID")
		return
	}

	var updatePriceDto dto.UpdatePriceDto
	if err := json.NewDecoder(r.Body).Decode(&updatePriceDto); err != nil {
		util.HttpTraceError(err, "invalid request payload", span, handler.loki, "UpdatePrice", accommodationId.Hex())
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := dto.ValidateUpdatePriceDto(&updatePriceDto); err != nil {
		util.HttpTraceError(err, "failed to validate payload", span, handler.loki, "UpdatePrice", accommodationId.Hex())
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.service.UpdatePrice(accommodationId, updatePriceDto, span, handler.loki); err != nil {
		util.HttpTraceError(err, "failed to update accommodation price", span, handler.loki, "UpdatePrice", accommodationId.Hex())
		handleError(w, http.StatusInternalServerError, "failed to update accommodation price")
		return
	}
	util.HttpTraceInfo("Accommodation price updated successfully", span, handler.loki, "Add", "")

	w.WriteHeader(http.StatusOK)
}

func (handler *AccommodationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "delete-delete")
	defer func() { span.End() }()
	vars := mux.Vars(r)
	accommodationId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		util.HttpTraceError(err, "Invalid accommodation id", span, handler.loki, "Delete", "")
		handleError(w, http.StatusInternalServerError, "Invalid accommodation id")
		return
	}

	if err := handler.service.Delete(accommodationId, span, handler.loki); err != nil {
		if err.Error() == "accommodation could not be deleted" {
			util.HttpTraceError(err, "Accommodation could not be deleted", span, handler.loki, "Delete", accommodationId.Hex())
			handleError(w, http.StatusPreconditionFailed, "accommodation could not be deleted")
		} else {
			util.HttpTraceError(err, "failed to delete accommodation", span, handler.loki, "Delete", accommodationId.Hex())
			handleError(w, http.StatusInternalServerError, "failed to delete accommodation")
		}

		return
	}
	util.HttpTraceInfo("Accommodations deleted successfully", span, handler.loki, "Delete", "")

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
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "get-by-hostId-get")
	defer func() { span.End() }()
	vars := mux.Vars(r)
	userId := vars["id"]
	accommodations, err := handler.service.GetByHostId(userId, span, handler.loki)

	if err != nil {
		util.HttpTraceError(err, "invalid host id", span, handler.loki, "GetByHostId", "")
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
		util.HttpTraceError(err, "failed to marshal data", span, handler.loki, "GetByHostId", userId)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	util.HttpTraceInfo("Accommodations retrieved successfully by host id", span, handler.loki, "GetByHostId", "")

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (handler *AccommodationHandler) GetById(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "get-by-id-get")
	defer func() { span.End() }()
	vars := mux.Vars(r)
	accommodationId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		util.HttpTraceError(err, "invalid accommodation id", span, handler.loki, "GetById", "")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	accommodation, err := handler.service.Get(accommodationId, span, handler.loki)

	if err != nil {
		util.HttpTraceError(err, "failed to get accommodation by id", span, handler.loki, "GetById", accommodationId.Hex())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response := dto.MapAccommodationResponse(*accommodation)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		util.HttpTraceError(err, "failed to marshal data", span, handler.loki, "GetById", accommodationId.Hex())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	util.HttpTraceInfo("Accommodation retrieved successfully by id", span, handler.loki, "GetById", "")

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (handler *AccommodationHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "get-all-get")
	defer func() { span.End() }()
	accommodations, err := handler.service.GetAll(span, handler.loki)

	if err != nil {
		util.HttpTraceError(err, "failed to get all accommodation", span, handler.loki, "GetAll", "")
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
		util.HttpTraceError(err, "failed to marshal data", span, handler.loki, "GetAll", "")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	util.HttpTraceInfo("All accommodations retrieved successfully", span, handler.loki, "GetAll", "")

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (handler *AccommodationHandler) GetImagesForAccommodations(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "get-images-for-accommodation-get")
	defer func() { span.End() }()
	var accommodationIds []dto.GetImagesRequest
	if err := json.NewDecoder(r.Body).Decode(&accommodationIds); err != nil {
		util.HttpTraceError(err, "invalid request payload", span, handler.loki, "GetImagesForAccommodations", "")
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	images, err := handler.service.GetImages(accommodationIds, span, handler.loki)
	if err != nil {
		util.HttpTraceError(err, "failed to get images", span, handler.loki, "GetImagesForAccommodations", "")
		handleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse, err := json.Marshal(images)
	if err != nil {
		util.HttpTraceError(err, "failed to marshal data", span, handler.loki, "GetImagesForAccommodations", "")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	util.HttpTraceInfo("Accommodation images successfully fetched", span, handler.loki, "GetImagesForAccommodations", "")

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (handler *AccommodationHandler) Add(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "add-post")
	defer func() { span.End() }()
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		util.HttpTraceError(err, "failed to parse form data", span, handler.loki, "Add", "")
		handleError(w, http.StatusBadRequest, "failed to parse form data")
		return
	}
	var createAccommodationDto dto.AccommodationDto
	jsonData := r.FormValue("json")
	if err := json.Unmarshal([]byte(jsonData), &createAccommodationDto); err != nil {
		util.HttpTraceError(err, "failed to parse json data", span, handler.loki, "Add", "")
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := dto.ValidateAccommodationDto(&createAccommodationDto); err != nil {
		util.HttpTraceError(err, "failed to validate payload", span, handler.loki, "Add", "")
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	photos, err := handlePhotoUploads(r, w, createAccommodationDto.HostId, span, nil)
	if err != nil {
		util.HttpTraceError(err, "failed to upload photos", span, handler.loki, "Add", "")
		handleError(w, http.StatusBadRequest, "failed to upload photos")
		return
	}

	newAccommodation := dto.MapAccommodation(&createAccommodationDto)
	newAccommodation.Photos = photos

	if err := handler.service.Add(newAccommodation, span, handler.loki); err != nil {
		util.HttpTraceError(err, "failed to add accommodation", span, handler.loki, "Add", "")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	util.HttpTraceInfo("Accommodations added successfully", span, handler.loki, "Add", "")
	w.WriteHeader(http.StatusCreated)
}

func (handler *AccommodationHandler) OnDeleteAccommodations(message *kafka.Message) {
	ctx := context.Background()
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(ctx, "on-delete-accommodations")
	defer func() { span.End() }()

	var deleteAccommodationRequest dto.DeleteAccommodationsRequest
	if err := json.Unmarshal(message.Value, &deleteAccommodationRequest); err != nil {
		util.HttpTraceError(err, "failed to marshal data", span, handler.loki, "OnDeleteAccommodations", "")
		log.Printf("Error unmarshalling rating change request: %v", err)
	}

	handler.service.OnDeleteAccommodations(deleteAccommodationRequest.HostId, span, handler.loki)
}

func handleError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
	fmt.Fprintf(w, message)
}

func handlePhotoUploads(r *http.Request, w http.ResponseWriter, hostId string, span trace.Span, handler *AccommodationHandler) ([]string, error) {
	var photos []string
	for _, fileHeader := range r.MultipartForm.File["photos"] {
		file, err := fileHeader.Open()
		if err != nil {
			util.HttpTraceError(err, "failed to open file", span, handler.loki, "Login", "")
			handleError(w, http.StatusInternalServerError, "failed to open file")
			return nil, err
		}
		defer file.Close()

		uploadDir := "./uploads"
		if err := os.MkdirAll(uploadDir, 0777); err != nil {
			fmt.Println("Error creating directory:", err)
		}
		dst, err := os.Create(filepath.Join(uploadDir, hostId+"-"+fileHeader.Filename))
		if err != nil {
			util.HttpTraceError(err, "failed to create file", span, handler.loki, "Login", "")
			fmt.Println("Error creating file:", err)
			handleError(w, http.StatusInternalServerError, "failed to create file")
			return nil, err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			util.HttpTraceError(err, "failed to save file", span, handler.loki, "Login", "")
			fmt.Println("Error Copy file:", err)
			handleError(w, http.StatusInternalServerError, "failed to save file")
			return nil, err
		}

		photos = append(photos, uploadDir+"/"+hostId+"-"+fileHeader.Filename)
	}
	return photos, nil
}
