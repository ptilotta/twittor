package models

type Secret struct {
	Host     string `json:"host"`
	Username string `json:"username"` // Alt izq + 96
	Password string `json:"password"`
	JWTSign  string `json:"jwtsign"`
	Database string `json:"database"`
}
