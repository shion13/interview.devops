package server

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

// HandlePhotoUpload handles the photo upload request.
// It reads the cloudProvider query parameter and the request body.
func (s *Server) HandlePhotoUpload(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()

	cloudProvider := queryParams.Get("cloudProvider")

	if cloudProvider == "" {
		http.Error(w, "please provide cloud provider", http.StatusBadRequest)
		return
	}

	if cloudProvider != "aws" && cloudProvider != "gcp" && cloudProvider != "azure" {
		http.Error(w, "invalid cloud provider, must be one of: aws, gcp, azure", http.StatusBadRequest)
		return
	}

	if cloudProvider == "gcp" || cloudProvider == "azure" {
		http.Error(w, "unable to process request for this cloud provider", http.StatusBadRequest)
		return
	}

	s.logger.Sugar().Debugf("received request with cloudProvider: \"%s\"", cloudProvider)

	bucketName := queryParams.Get("bucketName")
	if bucketName == "" {
		http.Error(w, "invalid bucket name", http.StatusBadRequest)
		return
	}

	prefix := queryParams.Get("prefix")

	body, err := io.ReadAll(r.Body)

	if !strings.HasPrefix(string(body), "\x89PNG\r\n\x1a\n") {
		http.Error(w, "invalid file type signature, must be png file", http.StatusBadRequest)
		return
	}

	if err != nil {
		s.logger.Sugar().Errorf("unable to read request body: %w", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	s.logger.Sugar().Debugf("received request body of size %d", len(body))

	s.processing <- true
	defer func() {
		<-s.processing
	}()

	err = s.bucketClients[cloudProvider].PushFileToBucket(bucketName, prefix, bytes.NewReader(body))

	if err != nil {
		if strings.Contains(err.Error(), "specified bucket does not exist") {
			http.Error(w, "bucket does not exist", http.StatusBadRequest)
			return
		}
		s.logger.Sugar().Errorf("unable to push file to bucket: %w", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}
