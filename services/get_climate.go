package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/leonardfreitas/go-gcloud-run/models"
)

const key = "037bb38eb1914c638d1183342242901"

type WeatherClimateApiResponse struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
}

type WeatherApiResponse struct {
	Current WeatherClimateApiResponse `json:"current"`
}

func GetWeather(city string, client HTTPClient) (*models.Climate, error) {
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", key, city)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("climate not found")
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var weatherApiResponse WeatherApiResponse
	err = json.Unmarshal(responseData, &weatherApiResponse)
	if err != nil {
		return nil, errors.New("error decoding JSON response")
	}

	climate := models.Climate{
		Celsius:    weatherApiResponse.Current.TempC,
		Fahrenheit: weatherApiResponse.Current.TempF,
	}

	if climate.Celsius == 0 && climate.Fahrenheit == 0 {
		return nil, errors.New("internal server error")
	}

	return &climate, nil
}
