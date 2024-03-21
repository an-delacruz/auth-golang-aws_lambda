package app

import (
	"lambda-func/api"
	"lambda-func/database"
)

type App struct {
	AppHandler api.ApiHandler
}

func NewApp() App{
	//initialize the DB store
	//gets passed down into the api handler
	db := database.NewDynamoDBClient()
	apiHandler := api.NewApiHandler(db)

	return App{
		AppHandler: apiHandler,
	}
}