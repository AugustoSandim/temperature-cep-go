package handler

import (
	"encoding/json"
	"net/http"
	"regexp"
	"temperature-cep-go/service"
)

func GetTemperature(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if !isValidCEP(cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	temp, err := service.GetTemperatureByCEP(cep)
	if err != nil {
		if err.Error() == "can not find zipcode" {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(temp)
}

func isValidCEP(cep string) bool {
	re := regexp.MustCompile(`^\d{8}$`)
	return re.MatchString(cep)
}
