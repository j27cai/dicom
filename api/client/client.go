package client

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"dicom/api/service/processor"
	"dicom/api/service/parser"
	"dicom/api/service/fetcher"
)

type Handler struct {
	dicomParser    parser.DicomParser
	dicomProcessor processor.DicomProcessor
	dicomFetcher   fetcher.DicomFetcher
}

// HandleDicomUpload handles the upload of a DICOM file
func (h *Handler) HandleDicomUpload(w http.ResponseWriter, r *http.Request) {
    // Parse the request body
    var requestBody struct {
        FilePath string `json:"file_path"`
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
    uuid, err := h.dicomParser.GetDicomDatasetByPath(requestBody.FilePath)
    if err != nil {
        // http.Error(w, "Error parsing DICOM file", http.StatusInternalServerError)
        return
    }

    // Convert the DICOM file to PNG
    err = h.dicomProcessor.ExtractDicomImage(h.dicomParser.GetDicomDataset())
    if err != nil {
        // http.Error(w, "Error extracting dicom image", http.StatusInternalServerError)
        return
    }

    //Respond with something
    response := map[string]string{"message": "your thing is %s", uuid}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *Handler) HandleGetImage(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
	if uuid == "" {
		http.Error(w, "UUID parameter is required", http.StatusBadRequest)
		return
	}

	// Get the DICOM image using the UUID
	img, err := h.dicomProvider.GetImage(id)
	if err != nil {
		http.Error(w, "Failed to fetch DICOM image", http.StatusInternalServerError)
		return
	}

	// Set the content type header
    w.Header().Set("Content-Type", "image/png") // Assuming the image format is PNG, adjust if different

    // Encode the image to the response writer
    if err := png.Encode(w, img); err != nil {
        http.Error(w, "Failed to encode image", http.StatusInternalServerError)
        return
    }
}

func Setup(router *mux.Router, dicomParser parser.DicomParser, dicomConverter converter.DicomConverter) {
	handler := Handler{dicomParser, dicomConverter}

	router.HandleFunc("/dicom", handler.HandleDicomUpload).Methods("POST")
	router.HandleFunc("/tags", handler.HandleGetTags).Methods("GET")
	router.HandleFunc("/image", handler.HandleGetImage).Methods("GET")
	router.HandleFunc("/health", ).Methods("GET")
	router.HandleFunc("/heartbeat", ).Methods("GET")

}
