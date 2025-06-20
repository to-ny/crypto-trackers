package config

import (
	"os"
	"strconv"
)

type Config struct {
	KafkaBootstrapServers string
	KafkaGroupID          string
	Port                  string
	LogLevel              string
	SpikeThreshold        float64
}

func New() *Config {
	return &Config{
		KafkaBootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS", "kafka-service:9092"),
		KafkaGroupID:          getEnv("KAFKA_GROUP_ID", "volume-spike-detector"),
		Port:                  getEnv("PORT", "8080"),
		LogLevel:              getEnv("LOG_LEVEL", "INFO"),
		SpikeThreshold:        getEnvFloat("SPIKE_THRESHOLD", 1.3),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}
