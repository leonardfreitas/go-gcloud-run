package main

import (
	"encoding/json"
	"net/http"

	"github.com/leonardfreitas/go-gcloud-run/models"
	"github.com/leonardfreitas/go-gcloud-run/services"
)

func GetClimate(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")

	if cep == "" {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	if len(cep) < 8 {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	cepInfo, err := services.GetCep(cep, services.DefaultHTTPClient)

	if err != nil {
		http.Error(w, "can not found zipcode", http.StatusNotFound)
		return
	}

	climate, err := services.GetWeather(cepInfo.City, services.DefaultHTTPClient)

	if err != nil {
		http.Error(w, "can not found zipcode", http.StatusNotFound)
		return
	}

	kelvin := services.GetKelvin(float64(climate.Celsius))

	payload := models.Report{
		TempC: climate.Celsius,
		TempF: climate.Fahrenheit,
		TempK: kelvin,
	}

	responseAPI, _ := json.Marshal(payload)
	w.Write(responseAPI)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", GetClimate)
	http.ListenAndServe(":8080", mux)
}
