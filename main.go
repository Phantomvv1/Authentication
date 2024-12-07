package main

import (
	"context"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
)

type Response struct {
	Response string `json:"response"`
}

func SHA512(text string) string {
	algorithm := sha512.New()
	algorithm.Write([]byte(text))
	result := algorithm.Sum(nil)
	return fmt.Sprintf("%x", result)
}

func signUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	err = conn.QueryRow(context.Background(), "select email from authentication where email = $1", email).Scan()
	emailExists := true
	if err != nil {
		if err == pgx.ErrNoRows {
			emailExists = false
		} else {
			log.Fatal(err)
		}
	}

	if emailExists {
		var response Response
		response.Response = "Email is already in use."
		json.NewEncoder(w).Encode(response)
		return
	}

	hashedPassword := SHA512(password)
	_, err = conn.Exec(context.Background(), "insert into authentication (email, password) values ($1, $2)", email, hashedPassword)
	if err != nil {
		log.Fatal(err)
	}
}

func logIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	err = conn.QueryRow(context.Background(), "select * from authentication where email = $1, password = $2", email, SHA512(password)).Scan()
	var response Response
	if err != nil {
		if err == pgx.ErrNoRows {
			response.Response = "Wrong email or password"
			json.NewEncoder(w).Encode(response)
		} else {
			log.Fatal(err)
		}
	}
	response.Response = "OK"
	json.NewEncoder(w).Encode(response)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/signUp", signUp)
	mux.HandleFunc("/logIn", logIn)

	http.ListenAndServe(":42069", mux)
}
