package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type reqBody struct {
	Money      int
	CandyType  string
	CandyCount int
}

type response struct {
	Change int    `json:"change"`
	Thanks string `json:"thanks"`
}

func simpleHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/buy_candy" {
		http.Error(writer, "404 not found", http.StatusNotFound)
		return
	}
	var requestBody reqBody
	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	candyMap := map[string]int{
		"CE": 10,
		"AA": 15,
		"NA": 17,
		"DE": 21,
		"YR": 23,
	}
	switch request.Method {
	case "POST":
		if err := request.ParseForm(); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if value, key := candyMap[requestBody.CandyType]; key && requestBody.Money > 0 && requestBody.CandyCount > 0 {
			if requestBody.Money >= value*requestBody.CandyCount {
				writer.WriteHeader(http.StatusCreated)
				writer.Header().Set("Content-Type", "application/json")
				resp := response{requestBody.Money - value*requestBody.CandyCount, "Thank you"}
				jsonResp, err := json.Marshal(resp)
				if err != nil {
					http.Error(writer, "400 Bad request", http.StatusBadRequest)
					return
				}
				writer.Write(jsonResp)
			} else {
				writer.WriteHeader(http.StatusPaymentRequired)
				writer.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(writer, "You need {%d} more money!\n", candyMap[requestBody.CandyType]-requestBody.Money)
			}
		} else {
			http.Error(writer, "400 Bad request", http.StatusBadRequest)
		}

	default:
		http.Error(writer, "405 Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", simpleHandler)
	http.ListenAndServe(":3333", nil)
}
