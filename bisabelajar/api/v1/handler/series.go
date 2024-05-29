package handler

import (
	"bisabelajar/dto"
	"bisabelajar/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type SeriesHandler struct {
	seriesService service.SeriesService
}

func NewSeriesHandler(seriesService service.SeriesService) *SeriesHandler {
	return &SeriesHandler{seriesService: seriesService}
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
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	id, err := b.seriesService.InsertSeries(seriesDto)
	if err != nil {
		http.Error(w, "Failed to insert series", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
}

func (b *SeriesHandler) GetSeriesByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid series ID", http.StatusBadRequest)
		return
	}

	series, err := b.seriesService.GetSeriesByID(id)
	if err != nil {
		http.Error(w, "Series not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(series)
}
