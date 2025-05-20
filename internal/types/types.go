package types

type Student struct{
	Id int `json:"id"`
	Name string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
	Age int `json:"age" validate:"required"`
}

type Response struct{
	Status string `json:"status"`
	Error string `json:"error"`
}