package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonardfreitas/go-gcloud-run/services"
)

func TestGetKelvin(t *testing.T) {
	result := services.GetKelvin(25.0)
	assert.Equal(t, 298.0, result)
}
