package models

type Errors struct {
	Message string `json:"message"`
}

type Success struct {
	Message string `json:"message"`
}

type Token struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Request struct {
	RefreshToken string `json:"refresh_token"`
}