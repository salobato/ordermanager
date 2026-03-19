package config

import (
	"log"
	"os"
)

type Config struct {
	MongoDatabase string
	MongoURI      string
	RabbitMQURI   string
	Port          string
}

func Load() *Config {
	cfg := &Config{
		MongoDatabase: getEnv("MONGO_DATABASE", "orderenv"),
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		RabbitMQURI:   getEnv("RABBITMQ_URI", "amqp://guest:guest@localhost:5672/"),
		Port:          getEnv("PORT", "8080"),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	log.Printf("Variável de ambiente %s não definida, utilizando valor padrão", key)
	return fallback
}
