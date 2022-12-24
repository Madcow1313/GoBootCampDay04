package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/lizrice/secure-connections/utils"
)

type requestFlags struct {
	Money      int    "json:money"
	CandyType  string "json:candyType"
	CandyCount int    "json:candyCount"
}

type respJson struct {
	Change int    "json:change"
	Thanks string "jsong:thanks"
}

func getClient() (*http.Client, error) {
	cp := x509.NewCertPool()
	data, err := os.ReadFile("../minica/minica.pem")
	if err != nil {
		return nil, err
	}
	cp.AppendCertsFromPEM(data)
	config := &tls.Config{
		RootCAs:               cp,
		GetClientCertificate:  utils.ClientCertReqFunc("./client-cert/cert.pem", "./client-cert/key.pem"),
		VerifyPeerCertificate: utils.CertificateChains,
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: config,
		},
	}
	return client, nil
}

func main() {
	var rFlags requestFlags
	flag.StringVar(&rFlags.CandyType, "k", "", "type of candy")
	flag.IntVar(&rFlags.CandyCount, "c", 0, "amount of candy to buy")
	flag.IntVar(&rFlags.Money, "m", 0, "amount of money")
	flag.Parse()
	client, err := getClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	js, err := json.Marshal(rFlags)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	resp, err := client.Post("https://localhost:3333/buy_candy", "application/json", bytes.NewBuffer(js))
	if err != nil {
		fmt.Println("Get Error")
		fmt.Println(err)
		os.Exit(-1)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	var respJs respJson
	err = json.Unmarshal(body, &respJs)
	if err != nil {
		fmt.Printf("Status: %s: %s\n", resp.Status, string(body))
	} else {
		fmt.Println(respJs.Thanks, "Your change is", respJs.Change)
	}
}
