package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
)

//extracting request headers
//extracting claims
//validating everything

func ValidateJWTMiddleware(next func (request events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse, error)) func (request events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse, error){
	return func(request events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse, error){
		//Extract the token from the header
		token := extractTokenFromHeader(request.Headers)
		//Validate the token
		if token == "" {
			return events.APIGatewayProxyResponse{
				Body: "Missing token",
				StatusCode: http.StatusUnauthorized,
			},nil
		}
		//If token is invalid, return unauthorized
		claims, err := parseToken(token)
		if err != nil {
			return events.APIGatewayProxyResponse{
				Body: "User unauthorized",
				StatusCode: http.StatusUnauthorized,
				},nil
			}
		
		expires := int64(claims["expires"].(float64))

		if time.Now().Unix() > expires {
			return events.APIGatewayProxyResponse{
				Body: "Token expired",
				StatusCode: http.StatusUnauthorized,
			},nil
		}

		//If token is valid, call the next function
		return next(request)
	}
}

func extractTokenFromHeader(headers map[string]string)string{
	authHeader,ok := headers["Authorization"]
	if !ok {
		return ""
	}
	splitToken := strings.Split(authHeader,"Bearer ")
	if len(splitToken) != 2 {
		return ""
	}
	return splitToken[1] 
}

func parseToken(tokenString string)(jwt.MapClaims,error){
	token, err := jwt.Parse(tokenString, func(token *jwt.Token)(interface{},error){
		secret := "secret"
		return []byte(secret),nil
	})

	if err != nil {
		return nil,fmt.Errorf("unauthorized")
	}

	if !token.Valid{
		return nil,fmt.Errorf("invalid token - unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil,fmt.Errorf("invalid claims - unauthorized")
	}

	return claims,nil
}