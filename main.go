package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/meals/{calorie}", handleMeals).Methods("GET")
	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), router))
}

func handleMeals(w http.ResponseWriter, r *http.Request) {
	calorie, _ := strconv.ParseFloat(mux.Vars(r)["calorie"], 64)
	response := GetMeals(context.Background(), calorie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
