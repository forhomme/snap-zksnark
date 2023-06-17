package models

var (
	SUCCESS = "Success"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ProofResponse struct {
	Hash string `json:"hash"`
}

type LoginResponse struct {
	AccessToken     string `json:"accessToken"`
	AccessExpireAt  string `json:"accessExpireAt"`
	RefreshToken    string `json:"refreshToken"`
	RefreshExpireAt string `json:"refreshExpireAt"`
}
