package model


type User struct {

	Id int64 `json:"id"`

	Login string `json:"login"`

	Email string `json:"email"`

	Password string `json:"password"`
}
