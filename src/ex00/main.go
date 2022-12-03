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
		if value, key := candyMap[req.CandyType]; key && req.Money > 0 && req.CandyCount > 0 {
			if req.Money >= value*req.CandyCount {
				w.WriteHeader(http.StatusCreated)
				w.Header().Set("Content-Type", "application/json")
				resp := response{req.Money - value*req.CandyCount, "Thank you"}
				jsonResp, err := json.Marshal(resp)
				if err != nil {
					http.Error(w, "400 Bad request", http.StatusBadRequest)
					return
				}
				w.Write(jsonResp)
			} else {
				w.WriteHeader(http.StatusPaymentRequired)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, "You need {%d} more money!\n", candyMap[req.CandyType]-req.Money)
			}
		} else {
			http.Error(w, "400 Bad request", http.StatusBadRequest)
		}

	default:
		http.Error(w, "405, Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", simpleHandler)
	http.ListenAndServe(":3333", nil)
}
