package main

import (
	"context"
	"crypto/sha512"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
)

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

	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "create table if not exists authentication (email text, password text);")
	if err != nil {
		log.Fatal(err)
	}

	var check string
	err = conn.QueryRow(context.Background(), "select email from authentication where email = $1;", email).Scan(&check)
	emailExists := true
	if err != nil {
		if err == pgx.ErrNoRows {
			emailExists = false
		} else {
			log.Fatal(err)
		}
	}

	if emailExists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	hashedPassword := SHA512(password)
	_, err = conn.Exec(context.Background(), "insert into authentication (email, password) values ($1, $2);", email, hashedPassword)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusCreated)
}

func logIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	var passwordCheck string
	err = conn.QueryRow(context.Background(), "select password from authentication a where a.email = $1;", email).Scan(&passwordCheck)
	if err != nil {
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			log.Fatal(err)
		}
	}

	if SHA512(password) != passwordCheck {
		w.WriteHeader(http.StatusUnauthorized)
	}
	w.WriteHeader(http.StatusAccepted)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/signUp", signUp)
	mux.HandleFunc("/logIn", logIn)

	http.ListenAndServe(":42069", mux)
}
