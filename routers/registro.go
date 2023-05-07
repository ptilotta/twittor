package routers

import (
	"context"
	"encoding/json"

	"github.com/ptilotta/twittor/bd"
	"github.com/ptilotta/twittor/models"
)

/*Registro es la funcion para crear en la BD el registro de usuario */
func Registro(ctx context.Context) (int, string) {

	var t models.Usuario
	err := json.Unmarshal([]byte(ctx.Value("body").(string)), &t)
	if err != nil {
		return 400, err.Error()
	}

	if len(t.Email) == 0 {
		return 400, "El email de usuario es requerido"
	}
	if len(t.Password) < 6 {
		return 400, "Debe especificar una contraseña de al menos 6 caracteres"
	}

	_, encontrado, _ := bd.ChequeoYaExisteUsuario(t.Email)
	if encontrado {
		return 400, "Ya existe un usuario registrado con ese email"
	}

	_, status, err := bd.InsertoRegistro(t)
	if err != nil {
		return 400, "Ocurrió un error al intentar realizar el registro de usuario " + err.Error()
	}

	if !status {
		return 400, "No se ha logrado insertar el registro del usuario"
	}

	return 200, "Registro OK"
}
