package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type WeatherResponse struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

const apiKey string = "Your api key"

func fetchWeatherDetails(city string) (WeatherResponse, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, apiKey)
	response, err := http.Get(url)
	if err != nil {
		return WeatherResponse{}, err
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return WeatherResponse{}, err
	}

	var weatherData WeatherResponse
	err = json.Unmarshal(responseData, &weatherData)
	if err != nil {
		return WeatherResponse{}, err
	}

	return weatherData, nil

}

func main() {

	cities := []string{"Hyderabad", "London", "Tokyo", "Paris", "Guntur"}

	var wg sync.WaitGroup
	var mu sync.Mutex

	var allResponses []WeatherResponse

	for _, city := range cities {
		wg.Add(1)
		go func(city string) {
			defer wg.Done()
			weatherData, err := fetchWeatherDetails(city)
			if err != nil {
				log.Printf("Error fetching for %s: %v", city, err)
				return
			}
			mu.Lock()
			allResponses = append(allResponses, weatherData)
			mu.Unlock()
		}(city)
	}

	wg.Wait()

    for _, res := range allResponses {
        fmt.Printf("City: %s\n", res.Name)
        fmt.Printf("Temperature: %.2f°C\n", res.Main.Temp-273.15)
        fmt.Printf("Humidity: %d%%\n", res.Main.Humidity)
        fmt.Printf("Weather: %s\n", res.Weather[0].Description)
        fmt.Println()
    }

}
