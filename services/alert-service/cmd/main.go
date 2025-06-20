package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"alert-service/internal/alerts"
	"alert-service/internal/config"
	"alert-service/internal/kafka"

	"github.com/gorilla/mux"
)

type HealthResponse struct {
	Status string `json:"status"`
}

type ReadyResponse struct {
	Status string `json:"status"`
}

type Server struct {
	config    *config.Config
	ready     bool
	consumer  *kafka.Consumer
	processor *alerts.AlertProcessor
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
		ready:  false,
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{Status: "healthy"})
}

func (s *Server) readyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := "ready"
	if !s.ready {
		status = "not ready"
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	json.NewEncoder(w).Encode(ReadyResponse{Status: status})
}

func (s *Server) setReady(ready bool) {
	s.ready = ready
}

func (s *Server) initializeKafka() error {
	brokers := strings.Split(s.config.KafkaBootstrapServers, ",")

	processor := alerts.NewAlertProcessor(s.config.CooldownMinutes)
	s.processor = processor

	consumer, err := kafka.NewConsumer(
		brokers,
		s.config.KafkaGroupID,
		[]string{"trading-signals"},
		processor.ProcessSignal,
	)
	if err != nil {
		return err
	}
	s.consumer = consumer

	return nil
}

func main() {
	cfg := config.New()
	server := NewServer(cfg)

	router := mux.NewRouter()
	router.HandleFunc("/health", server.healthHandler).Methods("GET")
	router.HandleFunc("/ready", server.readyHandler).Methods("GET")

	httpServer := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Printf("Starting Alert Service on port %s", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	if err := server.initializeKafka(); err != nil {
		log.Fatalf("Failed to initialize Kafka: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := server.consumer.Start(ctx); err != nil {
			log.Printf("Consumer error: %v", err)
		}
	}()

	server.setReady(true)
	log.Printf("Alert Service is ready and consuming trading signals (cooldown: %d minutes)", cfg.CooldownMinutes)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down Alert Service")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if server.consumer != nil {
		server.consumer.Close()
	}

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
}
