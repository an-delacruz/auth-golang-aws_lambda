package main

import (
	"fmt"
	"lambda-func/app"
	"lambda-func/middleware"

	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Username string `json:"username"`
}

func HandleRequest(event MyEvent)(string,error){
	if event.Username == "" {
		return "",fmt.Errorf("username cannot be empty")
	}

	return fmt.Sprintf("Hello %s!", event.Username),nil
}

func ProtectedHandler(request events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse,error){
	return events.APIGatewayProxyResponse{
		Body: "This is a protected route",
		StatusCode: http.StatusOK,
	},nil
}

func main()  {
	// fmt.Println("Starting lambda")
	myApp := app.NewApp()
	lambda.Start(func(request events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse,error){
		switch request.Path{
		case "/register":
			return myApp.AppHandler.RegisterUserHandler(request)
		case "/login":
			return myApp.AppHandler.LoginUserHandler(request)
		case "/protected":
			return middleware.ValidateJWTMiddleware(ProtectedHandler)(request)
		default:
			return events.APIGatewayProxyResponse{
				Body: "Not Found",
				StatusCode: http.StatusNotFound,
			},nil
		}
	})
}