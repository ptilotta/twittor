package models

/*RespuestaLogin tiene el token que se devuelve con el login */
type RespuestaLogin struct {
	Token string `json:"token,omitempty"`
}
