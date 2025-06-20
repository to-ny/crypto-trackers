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
	CooldownMinutes       int
}

func New() *Config {
	return &Config{
		KafkaBootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS", "kafka-service:9092"),
		KafkaGroupID:          getEnv("KAFKA_GROUP_ID", "alert-service"),
		Port:                  getEnv("PORT", "8080"),
		LogLevel:              getEnv("LOG_LEVEL", "INFO"),
		CooldownMinutes:       getEnvInt("COOLDOWN_MINUTES", 5),
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
