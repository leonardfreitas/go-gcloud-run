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

type CepMockHTTPClient struct {
	mock.Mock
}

func (m *CepMockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestGetCep(t *testing.T) {
	tests := []struct {
		name             string
		cep              string
		httpStatusCode   int
		responseBody     string
		expectedCEP      *models.CEP
		expectedError    error
		expectedErrorMsg string
	}{
		{
			name:           "Valid CEP",
			cep:            "12345678",
			httpStatusCode: http.StatusOK,
			responseBody: `{
				"cep": "12345678",
				"logradouro": "Rua Teste",
				"complemento": "Complemento Teste",
				"bairro": "Bairro Teste",
				"localidade": "Cidade Teste",
				"uf": "TS",
				"ibge": "1234567",
				"gia": "7890",
				"ddd": "11",
				"siafi": "9876"
			}`,
			expectedCEP: &models.CEP{
				CEP:          "12345678",
				Street:       "Rua Teste",
				Complement:   "Complemento Teste",
				Neighborhood: "Bairro Teste",
				City:         "Cidade Teste",
				State:        "TS",
			},
			expectedError:    nil,
			expectedErrorMsg: "",
		},
		{
			name:             "Invalid CEP",
			cep:              "00000000",
			httpStatusCode:   http.StatusNotFound,
			responseBody:     "{}",
			expectedCEP:      nil,
			expectedError:    errors.New("can not found zipcode"),
			expectedErrorMsg: "can not found zipcode",
		},
		{
			name:             "Error on HTTP request",
			cep:              "12345678",
			httpStatusCode:   http.StatusInternalServerError,
			responseBody:     "",
			expectedCEP:      nil,
			expectedError:    errors.New("can not found zipcode"),
			expectedErrorMsg: "can not found zipcode",
		},
		{
			name:             "Invalid JSON in response",
			cep:              "12345678",
			httpStatusCode:   http.StatusOK,
			responseBody:     "{invalid json}",
			expectedCEP:      nil,
			expectedError:    errors.New("can not found zipcode"),
			expectedErrorMsg: "can not found zipcode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHTTPClient := new(CepMockHTTPClient)

			mockResponse := &http.Response{
				StatusCode: tt.httpStatusCode,
				Body:       io.NopCloser(strings.NewReader(tt.responseBody)),
			}

			mockHTTPClient.On("Do", mock.Anything).Return(mockResponse, nil)

			result, err := services.GetCep(tt.cep, mockHTTPClient)

			assert.Equal(t, tt.expectedCEP, result)
			assert.Equal(t, tt.expectedError, err)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedErrorMsg)
			}
		})
	}
}
