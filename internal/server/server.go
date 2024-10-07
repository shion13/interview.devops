package server

import (
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gorilla/mux"
	"github.com/shion13/interview.devops/internal/bucket"
	"go.uber.org/zap"
)

type Server struct {
	bucketClients map[string]bucket.BucketUser
	router        *mux.Router
	processing    chan bool
	logger        *zap.Logger
}

func (s *Server) Setup(awsConfig aws.Config, logger *zap.Logger) error {
	s.setupChannel()
	s.setupLogger(logger)
	s.setupRoutes()
	s.bucketClients = make(map[string]bucket.BucketUser)
	s.bucketClients["aws"] = bucket.SetupS3User(awsConfig)
	return nil
}

func (s *Server) setupChannel() {
	channel := make(chan bool, 1)
	s.processing = channel
}

func (s *Server) setupRoutes() {
	router := mux.NewRouter()
	router.HandleFunc("/photo", s.HandlePhotoUpload).Methods(http.MethodPost)
	s.router = router
}

func (s *Server) setupLogger(logger *zap.Logger) {
	s.logger = logger
}

func (s *Server) Serve() error {
	srv := &http.Server{
		Addr:         ":8081",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      s.router,
	}
	s.logger.Sugar().Infof("listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		s.logger.Sugar().Infof("listening on %s: %v", srv.Addr, err)
	}
	return nil
}
