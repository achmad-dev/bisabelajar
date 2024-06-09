package handler

import (
	middlewarev1 "bisabelajar/api/v1/middleware"
	"bisabelajar/dto"
	"bisabelajar/service"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type ShortHandler struct {
	ShortService service.ShortService
	TempFilePath string
	Log          *logrus.Logger
}

func NewShortHandler(ShortService service.ShortService, tempFilePath string, log *logrus.Logger) ShortHandler {
	return ShortHandler{ShortService: ShortService, TempFilePath: tempFilePath, Log: log}
}

func (sh *ShortHandler) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/create", sh.InsertShort)
	return r
}

func (sh *ShortHandler) InsertShort(w http.ResponseWriter, r *http.Request) {
	// Parse form values
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max file size
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Extract form values
	seriesIDStr := r.FormValue("series_id")
	seriesID, err := strconv.Atoi(seriesIDStr)
	if err != nil {
		http.Error(w, "Invalid series_id", http.StatusBadRequest)
		return
	}
	title := r.FormValue("title")
	file, handler, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Failed to get video from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Ensure the temporary directory exists
	if err := os.MkdirAll(sh.TempFilePath, os.ModePerm); err != nil {
		http.Error(w, "Failed to create temporary directory", http.StatusInternalServerError)
		return
	}

	// Create a temporary file to store the uploaded video
	tempFilePath := filepath.Join(sh.TempFilePath, handler.Filename)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		sh.Log.Errorf("Failed to create temporary file: %v", err)
		http.Error(w, "Failed to create temporary file", http.StatusInternalServerError)
		return
	}

	// Copy the uploaded video to the temporary file
	if _, err := io.Copy(tempFile, file); err != nil {
		tempFile.Close()
		sh.Log.Errorf("Failed to copy video to temporary file: %v", err)
		http.Error(w, "Failed to copy video to temporary file", http.StatusInternalServerError)
		return
	}
	tempFile.Close() // Ensure the file is closed before removing it

	// Prepare DTO for short service
	shortDTO := dto.ShortDTO{
		SerieID: seriesID,
		Title:   title,
	}

	// Upload the video to Firestore
	requestContext := middlewarev1.GetRequestDetails(r)
	if err := sh.ShortService.UploadShort(tempFilePath, shortDTO, requestContext); err != nil {
		sh.Log.Errorf("Failed to upload video: %v", err)
		http.Error(w, "Failed to upload video", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
}
