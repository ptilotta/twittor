package routers

import (
	"context"
	"encoding/json"

	"github.com/ptilotta/twittor/bd"
	"github.com/ptilotta/twittor/models"
)

/*Registro es la funcion para crear en la BD el registro de usuario */
func Registro(ctx context.Context) models.RespApi {

	var t models.Usuario
	var r models.RespApi

	err := json.Unmarshal([]byte(ctx.Value("body").(string)), &t)
	if err != nil {
		r.Message = err.Error()
		return r
	}

	if len(t.Email) == 0 {
		r.Message = "El email de usuario es requerido"
		return r
	}
	if len(t.Password) < 6 {
		r.Message = "Debe especificar una contraseña de al menos 6 caracteres"
		return r
	}

	_, encontrado, _ := bd.ChequeoYaExisteUsuario(t.Email)
	if encontrado {
		r.Message = "Ya existe un usuario registrado con ese email"
		return r
	}

	_, status, err := bd.InsertoRegistro(t)
	if err != nil {
		r.Message = "Ocurrió un error al intentar realizar el registro de usuario " + err.Error()
		return r
	}

	if !status {
		r.Message = "No se ha logrado insertar el registro del usuario"
		return r
	}

	r.Status = 200
	r.Message = "Registro OK"
	return r
}
