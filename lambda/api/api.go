package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct{
	dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(request events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse,error){
	var registerUser types.RegisterUser

	err := json.Unmarshal([]byte(request.Body), &registerUser)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Invalid request",
			StatusCode: http.StatusBadRequest,
		},err
	}

	if(registerUser.Username == "" || registerUser.Password == ""){
		return events.APIGatewayProxyResponse{
			Body: "Username and password cannot be empty",
			StatusCode: http.StatusBadRequest,
		},err
	}	

	//Validate if username already exists

	userExists,err := api.dbStore.DoesUserExist(registerUser.Username)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		},err
	}

	if userExists {
		return events.APIGatewayProxyResponse{
			Body: "User already exists",
			StatusCode: http.StatusConflict,
		},nil
	}

	user,err := types.NewUser(registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		},fmt.Errorf("Error creating user %w",err)
	}

	err = api.dbStore.InsertUser(user)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		},fmt.Errorf("Error inserting user into the database %w",err)
	}

	return events.APIGatewayProxyResponse{
		Body: "User registered successfully",
		StatusCode: http.StatusOK,
	},nil
}

func (api ApiHandler) LoginUserHandler(request events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse,error){
	type LoginRequest struct{
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var loginRequest LoginRequest

	err := json.Unmarshal([]byte(request.Body), &loginRequest)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Invalid request",
			StatusCode: http.StatusBadRequest,
		},err
	}

	user,err := api.dbStore.GetUser(loginRequest.Username)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		},nil
	}

	if !types.ValidatePassword(user.PasswordHash, loginRequest.Password){
		return events.APIGatewayProxyResponse{
			Body: "Invalid user credentials",
			StatusCode: http.StatusUnauthorized,
		},nil
	}

	token := types.CreateToken(user)
	successMsg := fmt.Sprintf(`{"access_token":"%s"}`,token)

	return events.APIGatewayProxyResponse{
		Body: successMsg,
		StatusCode: http.StatusOK,
	},nil
}