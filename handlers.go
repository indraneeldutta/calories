package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func handleMeals(w http.ResponseWriter, r *http.Request) {
	// calorie, _ := strconv.ParseFloat(mux.Vars(r)["calorie"], 64)
	response := GetMeals(context.Background(), mux.Vars(r)["calorie"])
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleStoreMeals(w http.ResponseWriter, r *http.Request) {
	var reqBody RequestStoreMeals
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Invalid request")
		return
	}
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Invalid request")
		return
	}

	response := StoreMeals(context.Background(), reqBody)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	json.NewEncoder(w).Encode(response)
}

func handleGetUserMeals(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseInt(mux.Vars(r)["userid"], 2, 64)
	response := GetUserMeals(context.Background(), userID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	json.NewEncoder(w).Encode(response)
}
