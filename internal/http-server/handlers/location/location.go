package location

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	redisClient "weather-api/server/internal/cache/redis"
	sl "weather-api/server/pkg/logger"
	api "weather-api/server/pkg/weather-api-client"
)

func Handler(logger *slog.Logger, cache *redisClient.Client) http.HandlerFunc {
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
			logger.Error("Redis cache write failed", sl.Err(err))
		}

		fmt.Fprint(w, string(weatherData))
	}
}
