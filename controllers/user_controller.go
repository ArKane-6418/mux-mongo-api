package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ArKane-6418/mux-mongo-api/configs"
	"github.com/ArKane-6418/mux-mongo-api/models"
	"github.com/ArKane-6418/mux-mongo-api/responses"

	_ "github.com/ArKane-6418/mux-mongo-api/docs"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

// CreateUser godoc
// @Description Create a new user
// @Accept json
// @Produce json
// @Param user body models.User true "User Data"
// @Success 201 {object} responses.UserResponse{status=int,message=string,data=map[string]interface{}} "Successfully create a new user"
// @Failure 400 {object} responses.UserResponse{status=int,message=string,data=map[string]interface{}} "There is an issue with the request body"
// @Failure 500 {object} responses.UserResponse{status=int,message=string,data=map[string]interface{}} "Fail to create a new user"
// @Router /user [post]
func CreateUser() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, reader *http.Request) {
		// Set timeout to 10 seconds
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		defer cancel()

		// Validate request body
		// Grab the user information from the request body
		if err := json.NewDecoder(reader.Body).Decode(&user); err != nil {
			// Set the API status code
			responseWriter.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			// Convert the response struct to JSON
			json.NewEncoder(responseWriter).Encode(response)
			return
		}

		// Validate required fields
		// Check if user is a struct, has the correct fields, and the values are of the right type
		if validationErr := validate.Struct(&user); validationErr != nil {
			// Set the API status code
			responseWriter.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			// Convert the response struct to JSON
			json.NewEncoder(responseWriter).Encode(response)
			return
		}

		newUser := models.User{
			Id:       primitive.NewObjectID(),
			Name:     user.Name,
			Location: user.Location,
			Title:    user.Title,
		}

		// Insert the new user into the collection and handle errors
		result, err := userCollection.InsertOne(ctx, newUser)
		if err != nil {
			// Set the API status code
			responseWriter.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			// Convert the response struct to JSON
			json.NewEncoder(responseWriter).Encode(response)
			return
		}

		// Write successful creation status to the header
		responseWriter.WriteHeader(http.StatusCreated)
		response := responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}}
		// Convert the response struct to JSON
		json.NewEncoder(responseWriter).Encode(response)
	}
}

// GetUser godoc
// @Description Get a user by userID
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} responses.UserResponse{status=int,message=string,data=map[string]interface{}} "Successfully get a user by their ID"
// @Failure 400 {object} responses.UserResponse{status=int,message=string,data=map[string]interface{}} "user ID is not provided"
// @Failure 404 {object} responses.UserResponse{status=int,message=string,data=map[string]interface{}} "A user with the specified ID could not be found"
// @Router /user/{userId} [get]
func GetUser() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, reader *http.Request) {
		// Set timeout to 10 seconds
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(reader)
		userId := params["userId"]

		var user models.User
		defer cancel()

		if userId == " " {
			responseWriter.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "userid was not specified"}
			json.NewEncoder(responseWriter).Encode(response)
			return
		}

		objId, _ := primitive.ObjectIDFromHex(userId)

		// Retrieve the user with the associated id and handle any errors and convert the JSON into a Go struct
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
		if err != nil {
			// Set the API status code
			responseWriter.WriteHeader(http.StatusNotFound)
			response := responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(responseWriter).Encode(response)
			return
		}

		responseWriter.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": user}}
		json.NewEncoder(responseWriter).Encode(response)
	}
}

// UpdateUser godoc
// @Description Update a user's information by userID
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} responses.UserResponse{status=int,message=string,data=map[string]interface{}} "Successfully update a user"
// @Failure 400 {object} responses.UserResponse{status=int,message=string,data=map[string]interface{}} "user ID is not provided or request body is improperly formatted"
// @Failure 404 {object} responses.UserResponse{status=int,message=string,data=map[string]interface{}} "A user with the specified ID could not be found"
// @Failure 500 {object} responses.UserResponse{status=int,message=string,data=map[string]interface{}} "Fail to update a user's information"
// @Router /user/{userId} [put]
func UpdateUser() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, reader *http.Request) {
		// Set timeout to 10 seconds
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(reader)
		userId := params["userId"]

		var user models.User
		defer cancel()

		if userId == " " {
			responseWriter.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "userid was not specified"}
			json.NewEncoder(responseWriter).Encode(response)
			return
		}

		objId, _ := primitive.ObjectIDFromHex(userId)

		// Validate request body
		// Grab the user information from the request body
		if err := json.NewDecoder(reader.Body).Decode(&user); err != nil {
			// Set the API status code
			responseWriter.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			// Convert the response struct to JSON
			json.NewEncoder(responseWriter).Encode(response)
			return
		}

		// Validate required fields
		// Check if user is a struct, has the correct fields, and the values are of the right type
		if validationErr := validate.Struct(&user); validationErr != nil {
			// Set the API status code
			responseWriter.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			// Convert the response struct to JSON
			json.NewEncoder(responseWriter).Encode(response)
			return
		}

		update := bson.M{"name": user.Name, "location": user.Location, "title": user.Title}

		// Find and update the appropriate user based on objId and handle any errors
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
		if err != nil {
			// Set the API status code
			responseWriter.WriteHeader(http.StatusNotFound)
			response := responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(responseWriter).Encode(response)
			return
		}

		result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
		if err != nil {
			// Set the API status code
			responseWriter.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(responseWriter).Encode(response)
			return
		}

		// Retrieve updated user's details
		var updatedUser models.User

		// Check that only one user was matched
		if result.MatchedCount == 1 {
			err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
			if err != nil {
				// Set the API status code
				responseWriter.WriteHeader(http.StatusInternalServerError)
				response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(responseWriter).Encode(response)
				return
			}
		}

		responseWriter.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedUser}}
		json.NewEncoder(responseWriter).Encode(response)
	}
}

// DeleteUser godoc
// @Description Update a user's information by userID
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} responses.UserResponse{status=int,message=string,data=map[string]interface{}} "Successfully delete a user"
// @Failure 400 {object} responses.UserResponse{status=int,message=string,data=map[string]interface{}} "user ID is not provided"
// @Failure 404 {object} responses.UserResponse{status=int,message=string,data=map[string]interface{}} "A user with the specified ID could not be found"
// @Failure 500 {object} responses.UserResponse{status=int,message=string,data=map[string]interface{}} "Fail to delete a user"
// @Router /user/{userId}/ [delete]
func DeleteUser() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, reader *http.Request) {
		// Set timeout to 10 seconds
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(reader)
		userId := params["userId"]

		defer cancel()

		if userId == " " {
			responseWriter.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "userid was not specified"}
			json.NewEncoder(responseWriter).Encode(response)
			return
		}

		objId, _ := primitive.ObjectIDFromHex(userId)

		result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
		if err != nil {
			// Set the API status code
			responseWriter.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(responseWriter).Encode(response)
			return
		}

		// Could not find the user with the specified id
		if result.DeletedCount < 1 {
			// Set the API status code
			responseWriter.WriteHeader(http.StatusNotFound)
			response := responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User with specified ID was not found"}}
			json.NewEncoder(responseWriter).Encode(response)
			return
		}
		responseWriter.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "User was successfully deleted"}}
		json.NewEncoder(responseWriter).Encode(response)
	}
}

// GetAllUsers godoc
// @Description Retrieve all users from the database
// @Produce json
// @Success 200 {object} responses.UserResponse{status=int,message=string,data=map[string]interface{}} "Successfully retrieve all users"
// @Failure 500 {object} responses.UserResponse{status=int,message=string,data=map[string]interface{}} "Fail to retrieve all users"
// @Router /users/ [get]
func GetAllUsers() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, reader *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var users []models.User
		defer cancel()

		// Find all users
		results, err := userCollection.Find(ctx, bson.M{})

		if err != nil {
			responseWriter.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(responseWriter).Encode(response)
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleUser models.User
			// Decode the result into the user struct
			if err := results.Decode(&singleUser); err != nil {
				responseWriter.WriteHeader(http.StatusInternalServerError)
				response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(responseWriter).Encode(response)
			}

			users = append(users, singleUser)
		}
		responseWriter.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": users}}
		json.NewEncoder(responseWriter).Encode(response)
	}
}

// DeleteAllUsers godoc
// @Description Delete all users from the database
// @Produce json
// @Success 200 {object} responses.UserResponse{status=int,message=string,data=map[string]interface{}} "Successfully delete all users"
// @Failure 500 {object} responses.UserResponse{status=int,message=string,data=map[string]interface{}} "Fail to delete all users"
// @Router /users/ [delete]
func DeleteAllUsers() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, reader *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Find all users
		results, err := userCollection.Find(ctx, bson.M{})

		if err != nil {
			responseWriter.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(responseWriter).Encode(response)
			return
		}

		var count int = 0
		defer results.Close(ctx)
		for results.Next(ctx) {
			count += 1
		}

		// Delete all users
		deleteResults, deleteErr := userCollection.DeleteMany(ctx, bson.M{})

		if deleteErr != nil {
			responseWriter.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(responseWriter).Encode(response)
			return
		}

		if count != int(deleteResults.DeletedCount) {
			// Set the API status code
			responseWriter.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Did not delete all users"}}
			json.NewEncoder(responseWriter).Encode(response)
			return
		}

		responseWriter.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "All users were successfully deleted"}}
		json.NewEncoder(responseWriter).Encode(response)
	}
}
