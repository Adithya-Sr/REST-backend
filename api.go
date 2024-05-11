package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)



type APIError struct{
	Error string `json:"error"`
}


type apiFunc func(http.ResponseWriter, *http.Request)error


func makeHttpHandler(f apiFunc)http.HandlerFunc{
return func(w http.ResponseWriter, r *http.Request){
if err:=f(w,r);err!=nil{
	if err:=WriteJSON(w,http.StatusBadRequest,APIError{Error: err.Error()});err!=nil{
		log.Fatal(err)
	}
} 
}
}

type APIServer struct{
	listenAddr string
  store Store
}


func NewAPIServer(listenAddr string,store Store)*APIServer{
return &APIServer{
	listenAddr: listenAddr,
	store: store,
}
}


func (s *APIServer) Run()error{
router:=mux.NewRouter()
router.Handle("/healthCheck",makeHttpHandler(s.handleHealthCheck))
router.Handle("/product",makeHttpHandler(s.handleProduct))
router.Handle("/user",makeHttpHandler(s.handleUser))
router.Handle("/signup",makeHttpHandler(s.signup))
router.Handle("/login",makeHttpHandler(s.login))
err:=http.ListenAndServe(s.listenAddr,router);if err!=nil{
	return err
}
return nil
}


func (s *APIServer) handleHealthCheck(w http.ResponseWriter,r *http.Request)error{
return WriteJSON(w,http.StatusOK,map[string]string{"status":"ok"})
}



func (s *APIServer) signup(w http.ResponseWriter,r *http.Request)error{
if r.Method=="POST"{
userReq:=&CreateUserReq{}
if err:=json.NewDecoder(r.Body).Decode(&userReq);err!=nil{
	return err
}
if err:=ValidateInput(userReq);err!=nil{
	return err
}
hash,err:=HashPassword(userReq.Password)
if err!=nil{
	return err
}
user:=NewUser(userReq.Email,hash)
if err:=s.store.CreateUser(user);err!=nil{
	return err
}
return WriteJSON(w,http.StatusOK,user)	
}
return fmt.Errorf("method not allowed!:%s",r.Method)
}








func (s *APIServer) login(w http.ResponseWriter,r *http.Request)error{

if r.Method=="POST"{
userReq:=&CreateUserReq{}
if err:=json.NewDecoder(r.Body).Decode(&userReq);err!=nil{
	return err
}
if err:=ValidateInput(userReq);err!=nil{
	return err
}
user,err:=s.store.GetUser(userReq.Email)
if err!=nil{
	return err
}
if err:=ComparePassword(user.Password,userReq.Password);err!=nil{
	return err
}
jwtStr,err:=SignJWT(user.Email)
if err!=nil{
	return err
}
setCookie(w,jwtStr) 
return nil


}
return fmt.Errorf("method not allowed!:%s",r.Method)
}






func (s *APIServer) handleUser(w http.ResponseWriter,r *http.Request)error{
switch r.Method{
case "GET":
  return  s.handleGetUsers(w,r)
case "DELETE":
	return s.handleDeleteUser(w,r)

default:
	return fmt.Errorf("method not allowed!:%s",r.Method)
}

}





func (s*APIServer) handleGetUsers(w http.ResponseWriter,r *http.Request)error{
users,err:=s.store.GetUsers()
if err!=nil{
	return err
}
return WriteJSON(w,http.StatusOK,users)
}






func (s*APIServer) handleDeleteUser(w http.ResponseWriter,r *http.Request)error{
vars:=mux.Vars(r)
email:=vars["email"]
if err:=s.store.DeleteUser(email);err!=nil{
	return err
}
return WriteJSON(w,http.StatusOK,map[string]string{"deleted":email})
}







func (s*APIServer) handleProduct(w http.ResponseWriter,r *http.Request)error{
switch r.Method{
case "GET":
     id := r.URL.Query().Get("id")
    if id==""{
     return s.handleGetProducts(w,r)
		}
	  return s.handleGetProduct(w,r)
case "POST":
	return s.handleCreateProduct(w,r)

case "DELETE":
	return s.handleDeleteProduct(w,r)

case "PATCH":
	return s.handleUpdateProduct(w,r) 
default:
	return fmt.Errorf("method not allowed!:%s",r.Method)
}
 
}


func (s*APIServer) handleCreateProduct(w http.ResponseWriter,r *http.Request)error{
prodReq:=&CreateProductReq{}
if err:=json.NewDecoder(r.Body).Decode(&prodReq);err!=nil{
	return err
}
if err:=ValidateInput(prodReq);err!=nil{
	return err
}
product:=NewProduct(prodReq.Name,prodReq.Price,prodReq.Desc)
if err:=s.store.CreateProduct(product);err!=nil{
	return err
}
return WriteJSON(w,http.StatusOK,product)
}





func (s*APIServer) handleGetProduct(w http.ResponseWriter,r *http.Request)error{
idStr := r.URL.Query().Get("id")
id,err:=strconv.Atoi(idStr)
if err!=nil{
	return err
}
product,err:=s.store.GetProduct(id)
if err!=nil{
	return err
}
return WriteJSON(w,http.StatusOK,product)
}





func (s*APIServer) handleGetProducts(w http.ResponseWriter,r *http.Request)error{
products,err:=s.store.GetProducts()
if err!=nil{
	return err
}
return WriteJSON(w,http.StatusOK,products)
}





func (s*APIServer) handleUpdateProduct(w http.ResponseWriter,r *http.Request)error{
prodReq:=&CreateProductReq{}
if err:=json.NewDecoder(r.Body).Decode(prodReq);err!=nil{
	return err
}
if err:=ValidateInput(prodReq);err!=nil{
	return err
}
Product:=NewProduct(prodReq.Name,prodReq.Price,prodReq.Desc)
if err:=s.store.UpdateProduct(Product);err!=nil{
	return err
}
return WriteJSON(w,http.StatusOK,map[string]int{"updated":Product.ID})
}




func (s*APIServer) handleDeleteProduct(w http.ResponseWriter, r *http.Request)error{
idStr := r.URL.Query().Get("id")
id,err:=strconv.Atoi(idStr)
if err!=nil{
	return err
}
if err:=s.store.DeleteProduct(id);err!=nil{
	return err
}
return WriteJSON(w,http.StatusOK,map[string]int{"deleted":id})
}

