package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)



func WriteJSON(w http.ResponseWriter,status int, v any)error{
w.Header().Add("Content-Type","application/json")
w.WriteHeader(status)
if err:=json.NewEncoder(w).Encode(v);err!=nil{
	return err
}
return nil
}


func InitEnv()error{
	err := godotenv.Load()
  if err != nil {
    return err
  }
	return nil
}


var Validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())	

func ValidateInput(input any)error{
 if err:=Validate.Struct(input);err!=nil{
	return err
 }
 return nil
}


func HashPassword( password string)(string,error){
hashByte,err:=bcrypt.GenerateFromPassword([]byte(password),10)
if err!=nil{
	return "",err
}
return string(hashByte),nil
}


func ComparePassword(hash,pswd string)error{
	return bcrypt.CompareHashAndPassword([]byte(hash),[]byte(pswd))
}


func ValidateJWT(tokenString  string)(bool,error){
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
 
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(os.Getenv("JWT_SAMPLE_KEY")), nil
})
if err!=nil{
	return false,err
}
 return token.Valid,nil

}



func SignJWT(data any)(string,error){
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	"data": data,
	"ExpiresAt":604800,
})
return token.SignedString([]byte(os.Getenv("JWT_SAMPLE_KEY")))
}

func setCookie(w http.ResponseWriter,jwt string){
    cookie := http.Cookie{
        Name:  "jwt",
        Value: jwt,
				Path: "/",
				HttpOnly: true,
				Secure: false,
				Expires:  time.Now().Add(24 * 7*time.Hour),
        MaxAge:   86400,

    }
    http.SetCookie(w, &cookie)

}


func VerifyAccess(f http.HandlerFunc)http.HandlerFunc{
return func(w http.ResponseWriter, r *http.Request) {
	cookie,err:=r.Cookie("jwt")
	if err!=nil{
		WriteJSON(w,http.StatusUnauthorized,APIError{Error:"access dennied"})
		return 
	}
	token:=cookie.Value
	if ok,err:=ValidateJWT(token);err!=nil || !ok{
		WriteJSON(w,http.StatusUnauthorized,APIError{Error:"access dennied"})
		return 
	}
	
	 f(w,r)
}
}