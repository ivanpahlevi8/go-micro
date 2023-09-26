package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// create handler function for serve broker
func (app *AppConfig) ServeBroker(w http.ResponseWriter, r *http.Request) {
	// create json object
	resObj := JsonRes{
		Error:   false,
		Message: "Testing serve broker",
	}

	_ = app.WriteJsonBody(w, http.StatusOK, resObj)
}

// create request payload to receive payload from broker body request
type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth"`
}

// create auth payload, to hold authentication payload like email and password
/**
auth payload akan digunakan sebagai body untuk melakukan request ke authentication-service
sehingga akan digunakan sebagai body melalui service ini
*/
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// create handler function to authentication
func (app *AppConfig) HandleAuthentication(w http.ResponseWriter, r *http.Request) {
	// create request payload
	var requestPayload RequestPayload

	// get response from request body
	err := app.ReadJsonBody(w, r, &requestPayload)

	// check for an error
	if err != nil {
		log.Println("errro when reading from json body request : ", err)
		app.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}

	// check action
	switch requestPayload.Action {
	case "auth":
		// user need to authentication
		app.Authenticated(w, requestPayload.Auth)
	default:
		app.ErrorJson(w, errors.New("error invalid action"), http.StatusInternalServerError)
		return
	}
}

// create function to authenticate user
func (app *AppConfig) Authenticated(w http.ResponseWriter, payload AuthPayload) {
	// create json from payload
	jsonPayload, err := json.MarshalIndent(payload, "", "\t")

	// check fro an error
	if err != nil {
		log.Println("errro when marshalling json")
		app.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}

	// create request
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonPayload))

	// check fro an error
	if err != nil {
		log.Println("errro when createing request")
		app.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}

	// create clien
	client := &http.Client{}

	// create response
	resp, err := client.Do(request)

	// check fro an error
	if err != nil {
		log.Println("errro when getting response")
		app.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}

	// defer close response body
	defer resp.Body.Close()

	// check for response status
	if resp.StatusCode == http.StatusUnauthorized {
		// if user is not authorized yet ot failed authorized
		log.Println("authorization fail")
		app.ErrorJson(w, errors.New("authorization fails"), http.StatusBadRequest)
		return
	} else if resp.StatusCode != http.StatusAccepted {
		// if user is not authorized yet ot failed authorized
		log.Println("authorization fail")
		app.ErrorJson(w, errors.New("authorization fails"), http.StatusBadRequest)
		return
	}

	// create payload
	var getPayload JsonRes

	// decode response body
	err = json.NewDecoder(resp.Body).Decode(&getPayload)
	// check fro an error
	if err != nil {
		log.Println("errro when decode response body")
		app.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}

	// check respons body
	if !getPayload.Error {
		log.Println("errro from payload body")
		app.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}

	// create response
	response := JsonRes{
		Error:   false,
		Message: "user success authenticated",
		Data:    getPayload.Data,
	}

	// write response
	app.WriteJsonBody(w, http.StatusOK, response)
}
