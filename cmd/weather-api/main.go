package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	redisClient "weather-api/server/internal/cache/redis"
	api "weather-api/server/pkg/weather-api-client"
)

func handler(cache *redisClient.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		city := r.URL.Path[1:]

		weather, err := cache.Get(r.Context(), city)
		if err == nil {
			fmt.Fprint(w, weather)
			return
		}

		weatherData := api.GetWeather(city)

		if err := cache.Set(
			r.Context(),
			city,
			string(weatherData),
			10*time.Minute,
		); err != nil {
			log.Printf("Redis cache write failed: %v", err)
		}

		fmt.Fprint(w, string(weatherData))
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Warning: could not load .env: %v", err)
		}
	}

	redis, err := redisClient.New(os.Getenv("REDIS_ADDR"))
	if err != nil {
		log.Fatal(err)
	}
	defer redis.Close()

	http.HandleFunc("/", handler(redis))
	http.HandleFunc("/health", healthCheck)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
