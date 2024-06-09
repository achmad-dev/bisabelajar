package handler

import (
	middlewarev1 "bisabelajar/api/v1/middleware"
	"bisabelajar/api/v1/response"
	"bisabelajar/dto"
	"bisabelajar/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type SeriesHandler struct {
	seriesService service.SeriesService
}

func NewSeriesHandler(seriesService service.SeriesService) SeriesHandler {
	return SeriesHandler{seriesService: seriesService}
}

func (b *SeriesHandler) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/create", b.InsertSeries)
	r.Get("/{id}", b.GetSeriesByID)
	return r
}

func (b *SeriesHandler) InsertSeries(w http.ResponseWriter, r *http.Request) {
	var seriesDto dto.SeriesDto
	if err := json.NewDecoder(r.Body).Decode(&seriesDto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.NewErrorResponse("Invalid request payload", err))
		return
	}
	fmt.Println("hello")
	requestContext := middlewarev1.GetRequestDetails(r)
	fmt.Println("hello2")
	_, err := b.seriesService.InsertSeries(seriesDto, requestContext)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewErrorResponse("Failed to insert series", err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.NewSuccessResponse(nil))
}

func (b *SeriesHandler) GetSeriesByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.NewErrorResponse("Invalid series ID", err))
		return
	}
	requestContext := middlewarev1.GetRequestDetails(r)
	series, err := b.seriesService.GetSeriesByID(id, requestContext)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response.NewErrorResponse("Series not found", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.NewSuccessResponse(series))
}
