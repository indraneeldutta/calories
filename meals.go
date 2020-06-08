package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Meal defines the structure of meals object
type Meal struct {
	Name string `json:"title"`
}

// ResponseGetMeals defines the structure of response received from API
type ResponseGetMeals struct {
	Meals     []Meal `json:"meals"`
	Nutrients struct {
		Calories float64 `json:"calories"`
	} `json:"nutrients"`
}

// RequestStoreMeals defined structure of request sent for storing meals
type RequestStoreMeals struct {
	Meals  []Meal `json:"meals"`
	Date   string `json:"date"`
	UserID int64  `json:"userid"`
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

// GetMeals returns set of meals based on the calories
func GetMeals(ctx context.Context, calories string) ResponseGetMeals {

	url := "https://api.spoonacular.com/mealplanner/generate?" + "apiKey=797b9554c0bd4543be9a2d1d0b165563" + "&targetCalories=" + calories + "&timeFrame=day"

	req, _ := http.NewRequest("GET", url, nil)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var response ResponseGetMeals
	if err = json.Unmarshal(body, &response); err != nil {
		return response
	}

	return response
}

// StoreMeals stores the meals selected for the particular day
func StoreMeals(ctx context.Context, request RequestStoreMeals) ResponseStoreMeals {
	client := GetClient()

	var check bson.M
	collection := client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("COLL_NAME")).FindOne(ctx, primitive.M{"userid": request.UserID, "date": request.Date})

	err := collection.Decode(&check)

	if err != nil {
		_, err = client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("COLL_NAME")).InsertOne(ctx, request)
	} else {
		collection := client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("COLL_NAME"))
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

	collection, err := client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("COLL_NAME")).Find(ctx, primitive.M{"userid": userID})

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
