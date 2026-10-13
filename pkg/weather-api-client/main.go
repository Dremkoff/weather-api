package weatherapiclient

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

const urlBase = "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/"

func buildURL(location string) string {
	params := url.Values{}
	params.Add("contentType", "json")
	params.Add("include", "days,hours,current,alerts")
	params.Add("key", os.Getenv("API_KEY"))
	params.Add("unitGroup", "metric")

	return urlBase + url.QueryEscape(location) + "?" + params.Encode()
}

func GetWeather(location string) []byte {
	url := buildURL(location)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
	return body
}
