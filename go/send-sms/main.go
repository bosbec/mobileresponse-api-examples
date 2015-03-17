package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	senderName string = "Bosbec"
	baseUrl    string = "https://api.mobileresponse.se/"
	message    string = "Hello from Golang!"
)

var (
	recipients []string = []string{"+46705176608", "+46735090065"}
	username   string   = os.Getenv("MOBILERESPONSE_API_USERNAME")
	password   string   = os.Getenv("MOBILERESPONSE_API_PASSWORD")
)

func main() {
	if username == "" {
		log.Fatal("Missing MobileResponse API username. Did you forget to set the MOBILERESPONSE_API_USERNAME environment variable?")
	}

	if password == "" {
		log.Fatal("Missing MobileResponse API password. Did you forget to set the MOBILERESPONSE_API_PASSWORD environment variable?")
	}

	fmt.Printf("Sending \"%s\" to %s\n", message, strings.Join(recipients, ", "))

	sendSms(username, password, message, recipients)

	fmt.Println("Authenticating")

	authenticationToken := authenticate(username, password)

	fmt.Printf("Authenticated and received \"%s\" authentication token\n", authenticationToken)

	fmt.Println("Checking if the authentication token is still valid")

	if isAuthenticated(authenticationToken) {
		fmt.Printf("\"%s\" is valid\n", authenticationToken)
	} else {
		fmt.Printf("\"%s\" is not valid\n", authenticationToken)
	}
}

func authenticate(username, password string) string {
	data := &Request{
		Data: &AuthenticateRequest{
			username,
			password,
		},
	}

	requestJson, _ := json.Marshal(data)

	client := &http.Client{}

	request, err := http.NewRequest("POST", baseUrl+"authenticate", bytes.NewReader(requestJson))

	if err != nil {
		panic(err)
	}

	response, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	var r = AuthenticateResponse{}

	err = json.Unmarshal(body, &r)

	return r.Data.Id
}

func sendSms(username, password, message string, recipients []string) {
	data := &Request{
		Data: &SendSmsRequest{
			username,
			password,
			message,
			recipients,
			senderName,
		},
	}

	requestJson, _ := json.Marshal(data)

	client := &http.Client{}

	request, err := http.NewRequest("POST", baseUrl+"quickie/send-message", bytes.NewReader(requestJson))

	if err != nil {
		panic(err)
	}

	response, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	response.Body.Close()
}

func isAuthenticated(authenticationToken string) bool {
	data := &Request{
		AuthenticationToken: authenticationToken,
	}

	requestJson, _ := json.Marshal(data)

	client := &http.Client{}

	request, err := http.NewRequest("POST", baseUrl+"is-authenticated", bytes.NewReader(requestJson))

	if err != nil {
		panic(err)
	}

	response, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	var r = IsAuthenticatedResponse{}

	err = json.Unmarshal(body, &r)

	return r.Status == "Success"

	return false
}

type Request struct {
	Data                interface{} `json:"data"`
	AuthenticationToken string      `json:"authenticationToken"`
}

type AuthenticateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SendSmsRequest struct {
	Username   string   `json:"username"`
	Password   string   `json:"password"`
	Message    string   `json:"message"`
	Recipients []string `json:"recipients"`
	SenderName string   `json:"senderName"`
}

type AuthenticateResponse struct {
	Data struct {
		Id string
	}
}

type IsAuthenticatedResponse struct {
	Status string
}
