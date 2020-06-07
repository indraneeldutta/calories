# Calories

Below are the REST APIs created for calories

1. `/meals/{calories}` - shows suggestions on meals with their calorie count.
2. `/meals/store` - stores the meals for the date and user.
3. `/user/meals` - gets the meals for the user and their date

## Setup

Clone this repo

The program uses MongoDB as its database. The DB Dumps can be found in `dbDump` folder in JSON format. Consists of two collections `meals` and `userMeals` with prepopulated data. Feel free to edit the data before import.

1. Change the connection string in `.env` to your Mongo instance.
2. Make sure to install dependencies to run the program. Run `make deps`
OR manually install the deps by running
`go get -u go.mongodb.org/mongo-driver/mongo` and `go get -u github.com/gorilla/mux` and `go get -u github.com/joho/godotenv`
3. run `make build` OR `go build -o bin/main`
4. run `make run` OR `./bin/main`


## Usage

The program serves default to port 8080. This can be changed in `.env`

Import `PostmanCollection.json` to your Postman

Execute the desired API by changing the values in Body for `POST` requests


## Assumptions

1. User details are not generated. 
2. The `userid` is hypothetical in this scenario as there is no user data connected to that `userid`
