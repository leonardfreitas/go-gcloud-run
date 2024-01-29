package services_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/leonardfreitas/go-gcloud-run/models"
	"github.com/leonardfreitas/go-gcloud-run/services"
)

type ClimateMockHTTPClient struct {
	mock.Mock
}

func (m *ClimateMockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestGetWeather(t *testing.T) {
	tests := []struct {
		name            string
		city            string
		httpStatusCode  int
		responseBody    string
		expectedClimate *models.Climate
		expectedError   error
	}{
		{
			name:           "Valid City",
			city:           "London",
			httpStatusCode: http.StatusOK,
			responseBody: `{
				"current": {
					"temp_c": 20.5,
					"temp_f": 68.9
				}
			}`,
			expectedClimate: &models.Climate{
				Celsius:    20.5,
				Fahrenheit: 68.9,
			},
			expectedError: nil,
		},
		{
			name:            "Invalid City",
			city:            "NonexistentCity",
			httpStatusCode:  http.StatusNotFound,
			responseBody:    "{}",
			expectedClimate: nil,
			expectedError:   errors.New("climate not found"),
		},
		{
			name:            "Error on HTTP request",
			city:            "London",
			httpStatusCode:  http.StatusNotFound,
			responseBody:    "",
			expectedClimate: nil,
			expectedError:   errors.New("climate not found"),
		},
		{
			name:            "Invalid JSON in response",
			city:            "London",
			httpStatusCode:  http.StatusOK,
			responseBody:    "{invalid json}",
			expectedClimate: nil,
			expectedError:   errors.New("error decoding JSON response"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHTTPClient := new(ClimateMockHTTPClient)

			// Configurando a resposta do cliente HTTP mockado
			mockResponse := &http.Response{
				StatusCode: tt.httpStatusCode,
				Body:       io.NopCloser(strings.NewReader(tt.responseBody)),
			}

			mockHTTPClient.On("Do", mock.Anything).Return(mockResponse, nil)

			// Injetando o cliente HTTP mockado
			result, err := services.GetWeather(tt.city, mockHTTPClient)

			assert.Equal(t, tt.expectedClimate, result)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
