package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	api "weather-api/server/pkg/weather-api-client"
)

func handler(w http.ResponseWriter, r *http.Request) {
	weather := api.GetWeather(r.URL.Path[1:])

	fmt.Fprint(w, string(weather))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
