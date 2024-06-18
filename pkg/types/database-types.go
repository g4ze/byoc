package types

type User struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Login struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
