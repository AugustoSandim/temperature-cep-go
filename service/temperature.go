package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"temperature-cep-go/model"
)

var (
	ViaCEPURL     = "https://viacep.com.br/ws/%s/json/"
	WeatherAPIURL = "http://api.weatherapi.com/v1/current.json?key=%s&q=%s"
	weatherAPIKey = "e93e91918e5d4069af1124522252207"
)

func GetTemperatureByCEP(cep string) (*model.Temperature, error) {
	location, err := getLocationByCEP(cep)
	if err != nil {
		return nil, err
	}

	weather, err := getTemperatureByLocation(location)
	if err != nil {
		return nil, err
	}

	tempC := weather.Current.TempC
	tempF := tempC*1.8 + 32
	tempK := tempC + 273

	return &model.Temperature{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}, nil
}

func getLocationByCEP(cep string) (string, error) {
	resp, err := http.Get(fmt.Sprintf(ViaCEPURL, cep))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var viaCEPErrorResponse model.ViaCEPErrorResponse
	if err := json.Unmarshal(body, &viaCEPErrorResponse); err == nil && viaCEPErrorResponse.Erro == "true" {
		return "", fmt.Errorf("can not find zipcode")
	}

	var viaCEPResponse model.ViaCEPResponse
	if err := json.Unmarshal(body, &viaCEPResponse); err != nil {
		return "", err
	}

	return viaCEPResponse.Localidade, nil
}

func getTemperatureByLocation(location string) (*model.WeatherAPIResponse, error) {
	escapedLocation := url.QueryEscape(location)
	url := fmt.Sprintf(WeatherAPIURL, weatherAPIKey, escapedLocation)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weatherAPIResponse model.WeatherAPIResponse
	if err := json.Unmarshal(body, &weatherAPIResponse); err != nil {
		return nil, err
	}

	return &weatherAPIResponse, nil
}
