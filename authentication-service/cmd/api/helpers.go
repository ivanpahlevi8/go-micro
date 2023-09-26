package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

// create object to hold json
type JsonRes struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

/**
helpers digunakan untuk memaminupalasi json
*/

// create function to read json from body
func (a *Config) ReadJsonBody(w http.ResponseWriter, r *http.Request, item interface{}) error {
	// set maximum data that can be hold by body
	maximumData := 1104857

	// set maximum body
	r.Body = http.MaxBytesReader(w, r.Body, int64(maximumData))

	// get data from body
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&item)

	// check for an error
	if err != nil {
		log.Println("error when decoding body from json")
		return err
	}

	// check if body only contain 1 json object
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		err = errors.New("json object must only contain one json in body")
		log.Println("json object must only contain one json in body")
		return err
	}

	return nil
}

// create function to write json
func (a *Config) WriteJsonBody(w http.ResponseWriter, status int, item interface{}, header ...http.Header) error {
	// set header as json response
	w.Header().Set("Content-Type", "application/json")

	// check if there is header or not
	if len(header) > 0 {
		for k, v := range header[0] {
			w.Header()[k] = v
		}
	}

	// set header status
	w.WriteHeader(status)

	// create json object
	jsonObject, err := json.MarshalIndent(item, "", "\t")

	// check for an error
	if err != nil {
		log.Println("error when converting object to json")
		return err
	}

	// write to output
	_, err = w.Write(jsonObject)

	// check for an error
	if err != nil {
		log.Println("error when write to http output")
		return err
	}

	return nil
}

// create function to create json error object
func (a *Config) ErrorJson(w http.ResponseWriter, errMsg error, status ...int) error {
	// create status code as error
	statusCode := http.StatusBadRequest

	// check if tghere is status passed
	if len(status) > 0 {
		statusCode = status[0]
	}

	// create json response
	var jsonResponse JsonRes

	// assign value
	jsonResponse.Error = false
	jsonResponse.Message = fmt.Sprintf("error happen : %s, with status code : %d\n", errMsg.Error(), statusCode)

	// write json
	return a.WriteJsonBody(w, statusCode, jsonResponse)
}
