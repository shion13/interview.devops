package handlers

import (
	"io"
	"log"
	"net/http"
)

// HandlePhotoUpload handles the photo upload request.
// It reads the cloudProvider query parameter and the request body.
func HandlePhotoUpload(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	cloudProvider := queryParams.Get("cloudProvider")
	log.Printf("received request with cloudProvider: \"%s\"", cloudProvider)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	log.Printf("received request body of size %d", len(body))

	w.WriteHeader(http.StatusOK)
}
