package entity

type User struct {
	Id       uint64
	Login    string `json:"login"`
	Password string `json:"password"`
}
