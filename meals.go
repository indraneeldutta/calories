package main

import (
	"context"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

// Meals defines the structure of meals object
type Meals struct {
	Name     string  `json:"name"`
	Calories float64 `json:"calories"`
}

// ResponseMeals defines the structure for Response sent for GetMeals
type ResponseMeals struct {
	Status int     `json:"status"`
	Body   []Meals `json:"body"`
}

// GetMeals returns all the meals from DB
func GetMeals(ctx context.Context) ResponseMeals {
	client := GetClient()
	cursor, err := client.Database("calories").Collection("meals").Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var meals []Meals
	if err = cursor.All(ctx, &meals); err != nil {
		log.Fatal(err)
	}

	response := ResponseMeals{
		Status: http.StatusOK,
		Body:   meals,
	}

	return response
}
