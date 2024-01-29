package services

func GetKelvin(celcius float64) float64 {
	kelvin := celcius + 273
	return kelvin
}
