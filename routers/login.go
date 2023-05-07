package routers

import (
	"context"
	"encoding/json"

	"github.com/ptilotta/twittor/bd"
	"github.com/ptilotta/twittor/jwt"
	"github.com/ptilotta/twittor/models"
)

/*Login realiza el login */
func Login(ctx context.Context) (int, string) {

	var t models.Usuario

	err := json.Unmarshal([]byte(ctx.Value("body").(string)), &t)

	if err != nil {
		return 400, "Usuario y/o Contraseña inválidos " + err.Error()
	}
	if len(t.Email) == 0 {
		return 400, "El email del usuario es requerido"
	}

	userData, existe := bd.IntentoLogin(t.Email, t.Password)
	if !existe {
		return 400, "Usuario y/o Contraseña inválidos "
	}

	jwtKey, err := jwt.GeneroJWT(ctx, userData)
	if err != nil {
		return 400, "Ocurrió un error al intentar general el Token correspondiente " + err.Error()
	}

	resp := models.RespuestaLogin{
		Token: jwtKey,
	}

	token, err2 := json.Marshal(resp)
	if err2 != nil {
		return 400, "Error al intentar formatear la respuesta a JSON"
	}
	return 200, string(token)
}
