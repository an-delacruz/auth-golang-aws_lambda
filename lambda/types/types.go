package types

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUser struct{
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct{
	Username string `json:"username"`
	PasswordHash string `json:"password"`
}

func NewUser(registeredUser RegisterUser)(User,error){
	hashedPassword,err := bcrypt.GenerateFromPassword([]byte(registeredUser.Password),10)
	if err != nil {
		return User{},err
	}

	return User{
		Username: registeredUser.Username,
		PasswordHash: string(hashedPassword),
	}, nil
}

func ValidatePassword(hashedPassword, plainTextPassword string)bool{
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainTextPassword))
	return err == nil
}

func CreateToken(user User)string{
	now := time.Now()

	validaUntil := now.Add(time.Hour * 1).Unix()

	claims := jwt.MapClaims{
		"username": user.Username,
		"expires": validaUntil,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims,nil)
	secret := "secret"

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return ""
	}

	return tokenString
}
