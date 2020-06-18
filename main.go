package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/get", getTodo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on localhost:%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("Healthy")
	w.WriteHeader(http.StatusOK)
	target := os.Getenv("TARGET")
	if target == "" {
		target = "Togoist"
	}
	fmt.Fprintf(w, "Welcome to the to-do list manager %s!\n Usage: \n /get -> Get a sample todo list item \n", target)
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	data := struct {
		ID            int       `todo:"id"`
		Title         string    `todo:"title"`
		Description   string    `todo:"description,omitempty"`
		Completed     bool      `todo:"completed"`
		TimeCompleted time.Time `todo:"timeCompleted,omitempty"`
		TimeCreated   time.Time `todo:"timeCreated"`
	}{
		ID:          1,
		Title:       "Build Phase 4 of todo app",
		Description: "Setup CloudSQL Postgres connection",
		Completed:   false,
		TimeCreated: time.Now(),
	}

	b, err := json.Marshal(data)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	w.Header().Set("Last-Modified", time.Now().Format(http.TimeFormat))

	w.Write(b)
}
