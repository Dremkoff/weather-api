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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redis, err := redisClient.New(os.Getenv("ADDRESS"))
	if err != nil {
		log.Fatal(err)
	}
	defer redis.Close()

	http.HandleFunc("/", handler(redis))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
