package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
		return
	}
	log.Println(".env file loaded")
}

func getString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}

func getInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return i
}

type DatabaseConfig struct {
	Url string
}

func NewDataBaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Url: getString("DB_URL", ""),
	}
}

type LogConfig struct {
	Level  int
	Format string
}

func NewLogConfig() *LogConfig {
	return &LogConfig{
		Level:  getInt("LOG_LEVEL", 0),
		Format: getString("LOG_FORMAT", "json"),
	}
}

type ServerConfig struct {
	Port string
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Port: getString("HTTP_PORT", "8081"),
	}
}
