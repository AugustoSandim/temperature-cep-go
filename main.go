package main

import (
	"fmt"
	"net/http"
	"temperature-cep-go/handler"
)

func main() {
	http.HandleFunc("/temperature", handler.GetTemperature)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
