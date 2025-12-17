package config

import "os"

// Config holds application configuration loaded from environment variables.
type Config struct {
	MongoURI string
	MongoDB  string
	Port     string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() Config {
	return Config{
		MongoURI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:  getEnv("MONGO_DB", "todos"),
		Port:     getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
