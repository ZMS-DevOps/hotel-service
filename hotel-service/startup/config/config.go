package config

import "os"

type Config struct {
	Port              string
	GrpcPort          string
	HotelDBHost       string
	HotelDBPort       string
	HotelDBUsername   string
	HotelDBPassword   string
	BookingHost       string
	BookingPort       string
	SearchHost        string
	SearchPort        string
	BootstrapServers  string
	KafkaAuthPassword string
	JaegerHost        string
	LokiHost          string
}

func NewConfig() *Config {
	return &Config{
		Port:              os.Getenv("SERVICE_PORT"),
		HotelDBHost:       os.Getenv("DB_HOST"),
		HotelDBPort:       os.Getenv("DB_PORT"),
		HotelDBUsername:   os.Getenv("MONGO_INITDB_ROOT_USERNAME"),
		HotelDBPassword:   os.Getenv("MONGO_INITDB_ROOT_PASSWORD"),
		BookingHost:       os.Getenv("BOOKING_HOST"),
		BookingPort:       os.Getenv("BOOKING_PORT"),
		SearchHost:        os.Getenv("SEARCH_HOST"),
		SearchPort:        os.Getenv("SEARCH_PORT"),
		BootstrapServers:  os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		KafkaAuthPassword: os.Getenv("KAFKA_AUTH_PASSWORD"),
		JaegerHost:        os.Getenv("JAEGER_ENDPOINT"),
		LokiHost:          os.Getenv("LOKI_ENDPOINT"),
	}
}
