package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/leonardfreitas/go-gcloud-run/models"
)

type CEPApiResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func GetCep(cep string, client HTTPClient) (*models.CEP, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var responseAPI CEPApiResponse
	json.Unmarshal(responseData, &responseAPI)

	cepData := models.CEP{
		CEP:          responseAPI.Cep,
		Street:       responseAPI.Logradouro,
		Complement:   responseAPI.Complemento,
		Neighborhood: responseAPI.Bairro,
		City:         responseAPI.Localidade,
		State:        responseAPI.Uf,
	}

	if cepData.City == "" {
		return nil, errors.New("can not found zipcode")
	}

	return &cepData, nil
}
