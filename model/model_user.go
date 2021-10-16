package model


type User struct {

	Id int64 `json:"id,omitempty"`

	Login string `json:"login,omitempty"`

	Email string `json:"email,omitempty"`

	Password string `json:"password,omitempty"`
}
