package main

import (
	"context"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// RequestStoreMeals defined structure of request sent for storing meals
type RequestStoreMeals struct {
	Meals  []Meals `json:"meals"`
	Date   string  `json:"date"`
	UserID int64   `json:"userid"`
}

// ResponseStoreMeals defines the structure for response sent for storing meals
type ResponseStoreMeals struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
}

// RequestUserMeals defines the structure of request for getting user details
type RequestUserMeals struct {
	UserID int64 `json:"userid"`
}

// ResponseUserMeals defines the structure of response sent for user meals
type ResponseUserMeals struct {
	Status int                 `json:"status"`
	Body   []RequestStoreMeals `json:"Body"`
}

// GetMeals returns all the meals from DB
func GetMeals(ctx context.Context, calories float64) ResponseMeals {
	client := GetClient()

	pipeline := []bson.D{primitive.D{{"$sample", primitive.D{{"size", 3}}}}}
	cursor, err := client.Database("calories").Collection("meals").Aggregate(ctx, pipeline)

	var meals []Meals
	if err = cursor.All(ctx, &meals); err != nil {
		return ResponseMeals{
			Status: http.StatusInternalServerError,
			Body:   meals,
		}
	}

	return ResponseMeals{
		Status: http.StatusOK,
		Body:   meals,
	}
}

// StoreMeals stores the meals selected for the particular day
func StoreMeals(ctx context.Context, request RequestStoreMeals) ResponseStoreMeals {
	client := GetClient()

	var check bson.M
	collection := client.Database("calories").Collection("userMeals").FindOne(ctx, primitive.M{"userid": request.UserID, "date": request.Date})

	err := collection.Decode(&check)

	if err != nil {
		_, err = client.Database("calories").Collection("userMeals").InsertOne(ctx, request)
	} else {
		collection := client.Database("calories").Collection("userMeals")
		_, err = collection.UpdateOne(
			context.TODO(),
			primitive.M{
				"userid": request.UserID,
				"date":   request.Date,
			},
			primitive.D{
				{"$set", primitive.D{{"meals", request.Meals}}},
			},
		)
	}

	if err != nil {
		return ResponseStoreMeals{
			Status: http.StatusInternalServerError,
			Body:   "Failed to store meals",
		}
	}

	return ResponseStoreMeals{
		Status: http.StatusOK,
		Body:   "Successfully stored meals",
	}
}

// GetUserMeals returns the meals with details per user
func GetUserMeals(ctx context.Context, userID int64) ResponseUserMeals {
	client := GetClient()

	collection, err := client.Database("calories").Collection("userMeals").Find(ctx, primitive.M{"userid": userID})

	if err != nil {
		log.Fatal("somethings")
	}

	var result []RequestStoreMeals

	for collection.Next(context.TODO()) {
		var single RequestStoreMeals
		err := collection.Decode(&single)
		if err != nil {
			return ResponseUserMeals{
				Status: http.StatusInternalServerError,
				Body:   result,
			}
		}
		result = append(result, single)
	}

	return ResponseUserMeals{
		Status: http.StatusOK,
		Body:   result,
	}
}
