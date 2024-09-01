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

func worker(city string, results chan<- WeatherResponse, errors chan<- error) {
	weatherData, err := fetchWeatherDetails(city)
	if err != nil {
		errors <- fmt.Errorf("\nerror fetching for %s: %v\n", city, err)
		return
	}
	results <- weatherData
}

func main() {
	cities := []string{"Hyderabad", "London", "Tokyo", "Paris", "Guntur", "testCity"}
	results := make(chan WeatherResponse, len(cities))
	errors := make(chan error, len(cities))

	var wg sync.WaitGroup
	for _, city := range cities {
		wg.Add(1)
		go func(city string) {
			defer wg.Done()
			worker(city, results, errors)
		}(city)
	}

	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	var allResponses []WeatherResponse
	for range cities {
		select {
		case result := <-results:
			allResponses = append(allResponses, result)
		case err := <-errors:
			log.Println(err)
		}
	}

	for _, res := range allResponses {
		fmt.Printf("City: %s\n", res.Name)
		fmt.Printf("Temperature: %.2f°C\n", res.Main.Temp-273.15)
		fmt.Printf("Humidity: %d%%\n", res.Main.Humidity)
		fmt.Printf("Weather: %s\n", res.Weather[0].Description)
		fmt.Println()
	}
}
