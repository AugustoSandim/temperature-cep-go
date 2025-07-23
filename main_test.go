package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"temperature-cep-go/handler"
	"temperature-cep-go/model"
	"temperature-cep-go/service"
	"testing"
)

func TestGetTemperature(t *testing.T) {
	// Mock ViaCEP API
	viaCEPServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cep := r.URL.Path[len("/ws/") : len(r.URL.Path)-len("/json/")]
		if cep == "99999999" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"erro": "true"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"localidade": "São Paulo"}`))
	}))
	defer viaCEPServer.Close()

	// Mock WeatherAPI
	weatherAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		if q == "São Paulo" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"current": {"temp_c": 25.0}}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": {"message": "city not found"}}`))
		}
	}))
	defer weatherAPIServer.Close()

	// Override API URLs to use mock servers
	originalViaCEPURL := service.ViaCEPURL
	originalWeatherAPIURL := service.WeatherAPIURL
	defer func() {
		service.ViaCEPURL = originalViaCEPURL
		service.WeatherAPIURL = originalWeatherAPIURL
	}()
	service.ViaCEPURL = viaCEPServer.URL + "/ws/%s/json/"
	service.WeatherAPIURL = weatherAPIServer.URL + "?key=%s&q=%s"

	// Test cases
	tests := []struct {
		name           string
		cep            string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "valid cep",
			cep:            "01001000",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"temp_C":25,"temp_F":77,"temp_K":298}`,
		},
		{
			name:           "invalid cep",
			cep:            "123",
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   "invalid zipcode\n",
		},
		{
			name:           "not found cep",
			cep:            "99999999",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "can not find zipcode\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/temperature?cep="+tt.cep, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handler.GetTemperature)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var target model.Temperature
				json.Unmarshal(rr.Body.Bytes(), &target)
				if target.TempC != 25 || target.TempF != 77 || target.TempK != 298 {
					t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
				}
			} else {
				if rr.Body.String() != tt.expectedBody {
					t.Errorf("handler returned unexpected body: got %v want %v",
						rr.Body.String(), tt.expectedBody)
				}
			}
		})
	}
}
