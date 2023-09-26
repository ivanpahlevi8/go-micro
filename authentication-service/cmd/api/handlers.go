package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

// create handlers function to authenticate user
func (app *Config) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	// create payload to hold data from body
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// read body request
	err := app.ReadJsonBody(w, r, &payload)

	// check for an error
	if err != nil {
		log.Println("error when getting json from request body : ", err)
		app.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}

	// get user by email
	getUser, err := app.Model.User.GetUserByEmail(payload.Email)

	// check for an error
	if err != nil {
		log.Println("error when getting user by email from database : ", err)
		app.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}

	// authenticate user with payload password
	auth, err := getUser.AuthenticateUser(payload.Password)

	// check for an erorr
	if err != nil || !auth {
		// there is an error
		// or user entering wrong credentials
		err := errors.New("invalid credentials user input")
		app.ErrorJson(w, err, http.StatusUnauthorized)
		return
	}

	// if user success do authentication, create some feedback
	feedBack := JsonRes{
		Error:   false,
		Message: fmt.Sprintf("success do authentication to user : %s\n", payload.Email),
		Data:    getUser,
	}

	// write feedback to user
	app.WriteJsonBody(w, http.StatusAccepted, feedBack)
}

// for testing only
func (app *Config) test(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	getData, _ := app.Model.User.GetUserByEmail(query.Get("email"))

	app.WriteJsonBody(w, http.StatusOK, getData)
}
