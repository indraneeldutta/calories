package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/meals/{calorie}", handleMeals).Methods("GET")
	router.HandleFunc("/meals/store", handleStoreMeals).Methods("POST")
	router.HandleFunc("/user/meals/{userid}", handleGetUserMeals).Methods("GET")
	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), router))
}
