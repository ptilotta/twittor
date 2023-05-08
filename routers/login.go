package routers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ptilotta/twittor/bd"
	"github.com/ptilotta/twittor/jwt"
	"github.com/ptilotta/twittor/models"
)

/*Login realiza el login */
func Login(ctx context.Context) models.RespApi {

	var t models.Usuario
	var r models.RespApi
	r.Status = 400

	err := json.Unmarshal([]byte(ctx.Value("body").(string)), &t)

	if err != nil {
		r.Message = "Usuario y/o Contraseña inválidos " + err.Error()
		return r
	}
	if len(t.Email) == 0 {
		r.Message = "El email del usuario es requerido"
		return r
	}

	userData, existe := bd.IntentoLogin(t.Email, t.Password)
	if !existe {
		r.Message = "Usuario y/o Contraseña inválidos "
		return r
	}

	jwtKey, err := jwt.GeneroJWT(ctx, userData)
	if err != nil {
		r.Message = "Ocurrió un error al intentar general el Token correspondiente " + err.Error()
		return r
	}

	resp := models.RespuestaLogin{
		Token: jwtKey,
	}

	token, err2 := json.Marshal(resp)
	if err2 != nil {
		r.Message = "Error al intentar formatear la respuesta a JSON"
		return r
	}

	cookie := &http.Cookie{
		Name:    "token",
		Value:   jwtKey,
		Expires: time.Now().Add(24 * time.Hour),
	}
	cookieString := cookie.String()

	res := &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(token),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
			"Set-Cookie":                  cookieString,
		},
	}

	r.Status = 200
	r.Message = string(token)
	r.CustomResp = res

	return r
}
