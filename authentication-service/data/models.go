package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// create database timeout
const dbTimeout = 10 * time.Minute

// create model object that hold user
/**
model akan digunakan sebagai receiver pada function main
uintuk mengakses user beserta function dari user tersebut
dimana, user digunakan untuk authentication yang berisi methoid unbtu mendapatkan data user
dari database dan melakukan authentication
*/
type Models struct {
	User User
}

// database object
var dbObj *sql.DB

// create function to init database an create model object
func Init(db *sql.DB) Models {
	dbObj = db

	return Models{
		User: User{},
	}
}

// create user model
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"string"`
	Active    int    `json:"active"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// create function to get user from database by id
func (user *User) GetUserById(id int) (*User, error) {
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	// defer cancel
	defer cancel()

	// create variable to hold data object
	var getId int
	var getFirstName string
	var getLastName string
	var getEmail string
	var getPassword string
	var getActive int
	var getCreatedAt time.Time
	var getUpdatedAt time.Time

	// create query
	queryTxt := `
		SELECT id, first_name, last_name, email, password, active, created_at, updated_at FROM public.users
			WHERE id=$1
	`

	// run query
	query := dbObj.QueryRowContext(
		ctx,
		queryTxt,
		id,
	)

	// scan query
	err := query.Scan(
		&getFirstName,
		&getLastName,
		&getEmail,
		&getPassword,
		&getActive,
		&getCreatedAt,
		&getUpdatedAt,
	)

	// check for an error
	if err != nil {
		return nil, err
	}

	err = query.Err()

	// check for an error
	if err != nil {
		return nil, err
	}

	// create object
	returnObj := User{
		ID:        getId,
		FirstName: getFirstName,
		LastName:  getLastName,
		Email:     getEmail,
		Password:  getPassword,
		Active:    getActive,
		CreatedAt: getCreatedAt,
		UpdatedAt: getUpdatedAt,
	}

	// success return obj
	return &returnObj, nil
}

// create function to get user by email
func (user *User) GetUserByEmail(email string) (*User, error) {
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	// defer cancel
	defer cancel()

	// create variable to hold data object
	var getId int
	var getFirstName string
	var getLastName string
	var getEmail string
	var getPassword string
	var getActive int
	var getCreatedAt time.Time
	var getUpdatedAt time.Time

	// create query
	queryTxt := `select id, email, first_name, last_name, password, access_level, created_at, updated_at from public.users where email = $1`

	// run query
	query := dbObj.QueryRowContext(
		ctx,
		queryTxt,
		email,
	)

	// scan query
	err := query.Scan(
		&getId,
		&getFirstName,
		&getLastName,
		&getEmail,
		&getPassword,
		&getActive,
		&getCreatedAt,
		&getUpdatedAt,
	)

	// check for an error
	if err != nil {
		log.Println("error when scan query")
		return nil, err
	}

	err = query.Err()

	// check for an error
	if err != nil {
		log.Println("error when scan query second time")
		return nil, err
	}

	// create object
	returnObj := User{
		ID:        getId,
		FirstName: getFirstName,
		LastName:  getLastName,
		Email:     getEmail,
		Password:  getPassword,
		Active:    getActive,
		CreatedAt: getCreatedAt,
		UpdatedAt: getUpdatedAt,
	}

	// success return obj
	return &returnObj, nil
}

// create function to get all user
func (user *User) GetAllUser() ([]*User, error) {
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	// defer cancel
	defer cancel()

	// make return slice
	var returnedSlice []*User

	// create variable to hold data object
	var getId int
	var getFirstName string
	var getLastName string
	var getEmail string
	var getPassword string
	var getActive int
	var getCreatedAt time.Time
	var getUpdatedAt time.Time

	// create query
	queryTxt := `
		SELECT id, first_name, last_name, email, password, active, created_at, updated_at FROM public.users
	`

	// do query
	query, err := dbObj.QueryContext(ctx, queryTxt)

	// check for an erorr
	if err != nil {
		return returnedSlice, err
	}

	// do query
	for query.Next() {
		// scan query
		err := query.Scan(
			&getId,
			&getFirstName,
			&getLastName,
			&getEmail,
			&getPassword,
			&getActive,
			&getCreatedAt,
			&getUpdatedAt,
		)

		// check for an erorr
		if err != nil {
			return returnedSlice, err
		}

		// create object
		returnObj := User{
			ID:        getId,
			FirstName: getFirstName,
			LastName:  getLastName,
			Email:     getEmail,
			Password:  getPassword,
			Active:    getActive,
			CreatedAt: getCreatedAt,
			UpdatedAt: getUpdatedAt,
		}

		// add object to sliice
		returnedSlice = append(returnedSlice, &returnObj)
	}

	err = query.Err()

	// check for an error
	if err != nil {
		return nil, err
	}

	// if success
	return returnedSlice, nil
}

// create function to add user to database
func (user *User) AddUserToDatabase(userNew User) (int, error) {
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	// defer cancel
	defer cancel()

	// create query
	queryTxt := `
		INSERT INTO public.users (first_name, last_name, email, password, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	// execute query
	res, err := dbObj.ExecContext(
		ctx,
		queryTxt,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.Active,
		user.CreatedAt,
		user.UpdatedAt,
	)

	// check for an error
	if err != nil {
		return -1, err
	}

	// get id
	getId, err := res.LastInsertId()

	// check for an error
	if err != nil {
		return -1, err
	}

	// if success
	return int(getId), nil
}

// create function to update user
func (user *User) UpdateUserById(userUpdated User) error {
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	// defer cancel
	defer cancel()

	// create query
	queryTxt := `
		UPDATE public.users SET 
			first_name = $1, 
			last_name = $2, 
			email = $3, 
			password = $4, 
			active = $5, 
			created_at = $6, 
			updated_at = $7 
		WHERE id=$8
	`

	// execute query
	_, err := dbObj.ExecContext(
		ctx,
		queryTxt,
		userUpdated.FirstName,
		userUpdated.LastName,
		userUpdated.Email,
		userUpdated.Password,
		userUpdated.Active,
		userUpdated.CreatedAt,
		userUpdated.UpdatedAt,
		userUpdated.ID,
	)

	// check for an error
	if err != nil {
		return err
	}

	// if success
	return nil
}

// create function to delete user base user receiver id
func (user *User) Delete() error {
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	// defer cancel
	defer cancel()

	// create query
	querytxt := `
		DELETE from public.users WHERE id=$1
	`

	// exec query
	_, err := dbObj.ExecContext(ctx, querytxt, user.ID)

	// check for an error
	if err != nil {
		return err
	}

	// if success
	return nil
}

// create function to delete user base user receiver id
func (user *User) DeleteById(id int) error {
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	// defer cancel
	defer cancel()

	// create query
	querytxt := `
		DELETE from public.users WHERE id=$1
	`

	// exec query
	_, err := dbObj.ExecContext(ctx, querytxt, id)

	// check for an error
	if err != nil {
		return err
	}

	// if success
	return nil
}

// create function to reset password
func (user *User) ResetPassword(newPass string) error {
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	// defer cancel
	defer cancel()

	// create query
	queryTxt := `
		UPDATE public.users SET password = $1 WHERE id = $2
	`

	// execute query
	_, err := dbObj.ExecContext(ctx, queryTxt, newPass, user.ID)

	// check for an error
	if err != nil {
		return err
	}

	// if success
	return nil
}

// create function to authentication password
func (user *User) AuthenticateUser(pass string) (bool, error) {
	test := bcrypt.CompareHashAndPassword([]byte(""), []byte(""))

	if test == bcrypt.ErrMismatchedHashAndPassword {
		// if user input not matched encrypted password
		err := errors.New("input password wrong, please try again")
		return false, err
	} else if test != nil {
		// if user input not matched encrypted password
		err := errors.New("error in bycrypt, please try again")
		return false, err
	}

	// if success
	return true, nil
}
