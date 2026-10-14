package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"

	redisClient "weather-api/server/internal/cache/redis"
	"weather-api/server/internal/http-server/handlers/healthcheck"
	"weather-api/server/internal/http-server/handlers/location"
	mwLogger "weather-api/server/internal/http-server/middleware/logger"
	sl "weather-api/server/pkg/logger"
)

func main() {
	logger := sl.SetupLogger()

	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			logger.Error("Failed to load .env", sl.Err(err))
		}
	}

	logger.Info("Starting the application")
	logger.Debug(os.Getenv("API_KEY"))

	redis, err := redisClient.New(os.Getenv("REDIS_ADDR"))
	if err != nil {
		logger.Error("Failed to connect to Redis", sl.Err(err))
	}
	defer redis.Close()

	logger.Info("Connected to Redis")
	logger.Info("Starting server")

	mux := http.NewServeMux()
	mw := mwLogger.New(logger)

	mux.HandleFunc("/{location}", location.Handler(logger, redis))
	mux.HandleFunc("/health", healthcheck.Handler)

	if err := http.ListenAndServe(":8000", mw(mux)); err != nil {
		logger.Error("Failed to start server", sl.Err(err))
	}
}
