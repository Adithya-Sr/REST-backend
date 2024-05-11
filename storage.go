package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)


type Store interface{
	//user data
	CreateUser(*User)error
	GetUsers()([]*User,error)
	DeleteUser(string)error
	GetUser(string)(*User,error)
	//product data
	CreateProduct(*Product)error
  GetProduct(int)(*Product,error)
  GetProducts()([]*Product,error)
  UpdateProduct(*Product)error
  DeleteProduct(int)error
}

type MySqlStore struct{
db *sql.DB
}

func NewMySqlStore()(*MySqlStore,error){
mysqlPass:=os.Getenv("MYSQL_PASSWORD")
connStr:=fmt.Sprintf("root:%s@tcp(172.17.0.2:3306)/",mysqlPass)
db, err := sql.Open("mysql",connStr )
if err != nil {
	return nil,err
}
if err:=db.Ping();err!=nil{
return nil,err
}


log.Println("database created")
return &MySqlStore{db:db},nil
}



func (s *MySqlStore)init()error{
if err:=s.CreateTables();err!=nil{
	return err
}

log.Println("tables created")
return nil
}






func (s *MySqlStore) CreateTables()error{
query:=`CREATE DATABASE IF NOT EXISTS ecommerce;`
_, err := s.db.Exec(query)  
	if err != nil {
		return err
	}

query=`USE ecommerce;`
_, err = s.db.Exec(query)  
	if err != nil {
		return err
	}

query=`CREATE TABLE IF NOT EXISTS user (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(25) UNIQUE,
    password VARCHAR(100)
);`
_, err = s.db.Exec(query)  
	if err != nil {
		return err
	}
query=`CREATE TABLE IF NOT EXISTS product (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(25),
    price INT,
		description VARCHAR(100)
);`
_, err = s.db.Exec(query)  
	if err != nil {
		return err
	} 

return nil
}



 



 
 
 


func (s *MySqlStore)CreateUser(user *User)error{
query:=`insert into user
	(email,password)
	values (?,?)`
_, err:= s.db.Exec(query,user.Email,user.Password)  
	if err != nil {
		return err
	}
return nil
}


func (s *MySqlStore)GetUser(email string)(*User,error){
query:=`select * from user where email=?`
resp, err:= s.db.Query(query,email)  
	if err != nil {
		return nil,err
	}

for resp.Next(){
return scanIntoUsers(resp)
}
 return nil,fmt.Errorf("user %s not found",email)
}


func (s *MySqlStore)GetUsers()([]*User,error){
users:=[]*User{}
query:=`select * from user`
resp, err:= s.db.Query(query)  
	if err != nil {
		return nil,err
	}
for resp.Next(){
	var user *User
	user,err=scanIntoUsers(resp)
	if err!=nil{
		return nil,err
	}
	users=append(users,user )
}
return users,nil
}


func (s *MySqlStore)DeleteUser( email string)error{
query:=`delete from user where email=?`
_, err:= s.db.Query(query,email)  
	if err != nil {
		return err
	}
return nil
}






func (s *MySqlStore)CreateProduct(product *Product)error{
query:=`insert into product
	(name,price,description)
	values (?,?,?)`
_, err:= s.db.Exec(query,product.Name,product.Price,product.Desc)  
	if err != nil {
		return err
	}
return nil
}


func (s *MySqlStore)GetProduct(id int)(*Product,error){
query:=`select * from product where id=?`
resp, err:= s.db.Query(query,id)  
	if err != nil {
		return nil,err
	}
for resp.Next(){
	return scanIntoProds(resp)
}
return nil,fmt.Errorf("product %d not found",id)
}




func (s *MySqlStore)GetProducts()([]*Product,error){
products:=[]*Product{}
query:=`select * from product`
resp, err:= s.db.Query(query)  
	if err != nil {
		return nil,err
	}
for resp.Next(){
	var product *Product
	product,err=scanIntoProds(resp)
	if err!=nil{
		return nil,err
	}
	products=append(products,product) 
}
return products,nil
}




func (s *MySqlStore)UpdateProduct( product *Product)error{
query:=`UPDATE product
SET name=?,price=?,description=?
WHERE name=?;
`
_, err:= s.db.Query(query,product.Name,product.Price,product.Desc,product.ID)  
	if err != nil {
		return err
	}
return nil
}




func (s *MySqlStore)DeleteProduct(id int)error{
query:=`delete from product where id=?`
_, err:= s.db.Query(query,id)  
	if err != nil {
		return err
	}
return nil
}




func scanIntoProds(rows *sql.Rows)(*Product,error){
product:=&Product{}
if err:=rows.Scan(&product.ID,&product.Name,&product.Price,&product.Desc);err!=nil{
	return nil,err
}
return product,nil
}


func scanIntoUsers(rows *sql.Rows)(*User,error){
user:=&User{}
//have to resolve this
var id int
if err:=rows.Scan(&id,&user.Email,&user.Password);err!=nil{
	return nil,err
}
return user,nil
}
