package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/popsa-platform/interview.devops/internal/http/handlers"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/photo", handlers.HandlePhotoUpload).Methods(http.MethodPost)

	srv := &http.Server{
		Addr:         ":8081",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      router,
	}

	log.Printf("listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("listening on %s: %v", srv.Addr, err)
	}
}
