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

	"ma-signal-detector/internal/config"
	"ma-signal-detector/internal/kafka"
	"ma-signal-detector/internal/signals"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HealthResponse struct {
	Status string `json:"status"`
}

type ReadyResponse struct {
	Status string `json:"status"`
}

var (
	priceEventsProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "price_events_processed_total",
			Help: "Total number of price events processed",
		},
		[]string{"symbol"},
	)
	signalsGenerated = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "trading_signals_generated_total",
			Help: "Total number of trading signals generated",
		},
		[]string{"symbol", "type"},
	)
	processingTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "price_event_processing_seconds",
			Help: "Time spent processing price events",
		},
		[]string{"symbol"},
	)
)

func init() {
	prometheus.MustRegister(priceEventsProcessed)
	prometheus.MustRegister(signalsGenerated)
	prometheus.MustRegister(processingTime)
}

type Server struct {
	config    *config.Config
	ready     bool
	consumer  *kafka.Consumer
	producer  kafka.SignalProducer
	detector  *signals.MADetector
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

	producer, err := kafka.NewProducer(brokers)
	if err != nil {
		return err
	}
	s.producer = producer

	detector := signals.NewMADetector(producer, *priceEventsProcessed, *signalsGenerated, *processingTime)
	s.detector = detector

	consumer, err := kafka.NewConsumer(
		brokers,
		s.config.KafkaGroupID,
		[]string{"crypto-prices"},
		detector.ProcessPriceEvent,
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
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	httpServer := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Printf("Starting MA Signal Detector on port %s", cfg.Port)
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
	log.Println("MA Signal Detector is ready and consuming messages")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down MA Signal Detector")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if server.consumer != nil {
		server.consumer.Close()
	}
	if server.producer != nil {
		server.producer.Close()
	}

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
}