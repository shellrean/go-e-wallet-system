package dto

type AuthRes struct {
	UserId int64 `json:"-"`
	Token string `json:"token"`
}
