package client

import (
    "encoding/json"
    "image/png"
    "log"
    "net/http"

    "github.com/gorilla/mux"

    "dicom/api/model"
    "dicom/api/service/processor"
    "dicom/api/service/parser"
    "dicom/api/service/fetcher"
)

type Handler struct {
    dicomParser    *parser.DicomParser
    dicomProcessor *processor.DicomProcessor
    dicomFetcher   *fetcher.DicomFetcher
    logger         *log.Logger
}

func NewHandler(dicomParser *parser.DicomParser, dicomProcessor *processor.DicomProcessor, dicomFetcher *fetcher.DicomFetcher, logger *log.Logger) *Handler {
    return &Handler{
        dicomParser:    dicomParser,
        dicomProcessor: dicomProcessor,
        dicomFetcher:   dicomFetcher,
        logger:         logger,
    }
}

// HandleDicomUpload handles the upload of a DICOM file
func (h *Handler) HandleDicomUpload(w http.ResponseWriter, r *http.Request) {
    // Parse the request body
    var requestBody struct {
        FilePath string `json:"path"`
    }

    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&requestBody); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if requestBody.FilePath == "" {
        http.Error(w, "File path is empty", http.StatusBadRequest)
        return
    }

    // dicomparse
    dataset, uuid, err := h.dicomParser.GetDicomDatasetByPath(requestBody.FilePath)
    if err != nil {
        http.Error(w, "Error handling DICOM file", http.StatusInternalServerError)
        return
    }

    // Convert the DICOM file to PNG
    err = h.dicomProcessor.ExtractDicomImage(uuid, dataset)
    if err != nil {
        http.Error(w, "Error extracting DICOM image", http.StatusInternalServerError)
        return
    }

    // Extract tags from the dicom file
    err = h.dicomProcessor.ExtractDicomHeaders(uuid, dataset)
    if err != nil {
        http.Error(w, "Error extracting DICOM tags", http.StatusInternalServerError)
        return
    }

    // Respond with success message
    response := map[string]string{"id": uuid}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)

    h.logger.Printf("Successfully uploaded dicom file at: %s", requestBody.FilePath)
}

func (h *Handler) HandleGetImage(w http.ResponseWriter, r *http.Request) {
    uuid := r.URL.Query().Get("id")
    if uuid == "" {
        http.Error(w, "UUID parameter is required", http.StatusBadRequest)
        return
    }

    // Get the DICOM image using the UUID
    img, err := h.dicomFetcher.GetImage(uuid)
    if err != nil {
        http.Error(w, "Failed to fetch DICOM image", http.StatusInternalServerError)
        return
    }

    // Set the content type header
    w.Header().Set("Content-Type", "image/png")

    // Encode the image to the response writer
    if err := png.Encode(w, img); err != nil {
        http.Error(w, "Failed to encode image", http.StatusInternalServerError)
        return
    }

    h.logger.Printf("Successfully retrieved dicom file for: %s", uuid)
}

func (h *Handler) HandleGetTags(w http.ResponseWriter, r *http.Request) {
    // Extract UUID from the query parameter
    uuid := r.URL.Query().Get("id")
    if uuid == "" {
        http.Error(w, "ID parameter is required", http.StatusBadRequest)
        return
    }

    // Get tags associated with the UUID
    tags, err := h.dicomFetcher.GetTags(uuid)
    if err != nil {
        http.Error(w, "Failed to fetch DICOM tags", http.StatusInternalServerError)
        return
    }

    // Prepare response JSON
    response := struct {
        UUID string       `json:"uuid"`
        Tags []model.Tag  `json:"tags"`
    }{
        UUID: uuid,
        Tags: tags,
    }

    // Set content type and encode response as JSON
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        return
    }

    h.logger.Printf("Successfully retrieved tags for: %s", uuid)
}


func HealthCheck(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

func Heartbeat(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Alive"))
}

func Setup(router *mux.Router, dicomParser *parser.DicomParser, dicomProcessor *processor.DicomProcessor, dicomFetcher *fetcher.DicomFetcher, logger *log.Logger) {
    handler := NewHandler(dicomParser, dicomProcessor, dicomFetcher, logger)

    router.HandleFunc("/dicom", handler.HandleDicomUpload).Methods("POST")
    router.HandleFunc("/tags", handler.HandleGetTags).Methods("GET")
    router.HandleFunc("/image", handler.HandleGetImage).Methods("GET")
    router.HandleFunc("/health", HealthCheck).Methods("GET")
    router.HandleFunc("/heartbeat", Heartbeat).Methods("GET")
}
