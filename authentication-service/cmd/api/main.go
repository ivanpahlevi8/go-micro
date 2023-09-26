package main

import (
	"authentication/data"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var loopIterate = 0

// set web port
const webPort = "80"

// create variable to hold data of database
const maximumOpenConnection = 10
const maximumLifetime = 5 * time.Minute
const maxIdleDbConn = 5

type Config struct {
	DB    *sql.DB
	Model data.Models
}

func main() {
	log.Println("starting authentication server")

	// create database connection
	conn := ConnectDatabase()

	// check there is connection
	if conn == nil {
		log.Fatal("connectio failed, returne empty connection to database")
		return
	}

	// create config
	myConfig := Config{
		DB:    conn,
		Model: data.Init(conn),
	}

	// create server
	serv := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: myConfig.Routes(),
	}

	err := serv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}

// create function to open connection with database
func openDb(dsn string) (*sql.DB, error) {
	// get sql
	conn, err := sql.Open("pgx", dsn)

	// check for an error
	if err != nil {
		log.Println("error when connecting to database : ", err)
		return nil, err
	}

	// check connection
	err = conn.Ping()

	// check for an error
	if err != nil {
		log.Println("error when ping to database : ", err)
		return nil, err
	}

	// if success
	return conn, nil
}

// create function to get connection from datavase
func ConnectDatabase() *sql.DB {
	// parse flag
	flag.Parse()

	// get dsn from os variable
	getDsn := os.Getenv("DSN")

	// create connection with database
	// host=postgres port=5432 user=postgres password=03052001ivan dbname=users connect_timeout=5
	//fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%s", "localhost", "users", "postgres", "03052001ivan", "5432")
	conn, err := openDb(getDsn)

	// loop if not yet connected
	for {
		if err != nil {
			log.Println("trying to connecting to database ...")
			loopIterate++
		} else {
			log.Println("success connected to datavase ...")
			// set db characteristic
			conn.SetMaxOpenConns(maximumOpenConnection)
			conn.SetConnMaxLifetime(maximumLifetime)
			conn.SetMaxIdleConns(maxIdleDbConn)
			return conn
		}

		if loopIterate > 10 {
			log.Println("cannot connect to database")
			return nil
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(1 * time.Second)
		continue
	}
}
