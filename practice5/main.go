package main

import (
	"log"
	"net/http"
)

func main() {
	InitDB("postgres://adilkanatov:15432@localhost:5432/practice5")

	http.HandleFunc("/jobs", GetJobsHandler)

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
