package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/meals", handleMeals).Methods("GET")
	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), router))
}

func handleMeals(w http.ResponseWriter, r *http.Request) {
	response := GetMeals(context.Background())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if response.Body == nil {
		json.NewEncoder(w).Encode("No movies found")
	} else {
		json.NewEncoder(w).Encode(response)
	}
}
