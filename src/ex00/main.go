package main

import (
	"encoding/json"
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

func simpleHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/buy_candy" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	var req reqBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	candyMap := map[string]int{
		"CE": 10,
		"AA": 15,
		"NA": 17,
		"DE": 21,
		"YR": 23,
	}
	switch r.Method {
	case "POST":
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if value, key := candyMap[req.CandyType]; key {
			if req.Money >= value*req.CandyCount {
				w.WriteHeader(http.StatusCreated)
				resp := response{req.Money - value*req.CandyCount, "Thank you"}
				json.NewEncoder(w).Encode(resp)
			} else {
				w.WriteHeader(http.StatusPaymentRequired)
				resp := response{req.Money, "Thank you"}
				json.NewEncoder(w).Encode(resp)
			}

		}

	default:
		http.Error(w, "405, Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", simpleHandler)
	http.ListenAndServe(":3333", nil)
}
