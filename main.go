package main

import (
	"fmt"
	"log"
)



func main(){

if err:=InitEnv();err!=nil{
	log.Fatal("Error loading .env file")
}
store,err:=NewMySqlStore()
if err!=nil{
	log.Fatal(err)
}

if err:=store.init();err!=nil{
log.Fatal(err)
}

server:=NewAPIServer(":3000",store)
fmt.Println("server running...")
if err:=server.Run();err!=nil{
	log.Fatal(err)
}
}

