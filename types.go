package main

//validate this
type User struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

func NewUser(email,password string)*User{
return &User{Email: email,Password: password}
}

//assign unique id
type Product struct{
	ID   int  `json:"id"`
	Name string `json:"name"`
	Price string `json:"price"`
	Desc string `json:"desc"`
}

func NewProduct(name,price,desc string)*Product{
return &Product{Name: name,Price: price,Desc: desc}
}


type CreateProductReq struct{
	Name string `json:"name" validate:"required"`
	Price string `json:"price" validate:"required,number"`
	Desc string `json:"desc" validate:"required"`
}


type CreateUserReq struct{
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,max=8,min=6"`
  RepeatPassword string `json:"repeatPassword" validate:"required,eqfield=Password"`
}