package routers

import (
	"context"
	"encoding/json"

	"github.com/ptilotta/twittor/bd"
	"github.com/ptilotta/twittor/models"
)

/*ModificarPerfil modifica el perfil de usuario */
func ModificarPerfil(ctx context.Context, claim models.Claim) models.RespApi {

	var r models.RespApi
	r.Status = 400

	var t models.Usuario

	err := json.Unmarshal([]byte(string(ctx.Value("body").(models.Key))), &t)
	if err != nil {
		r.Message = "Datos Incorrectos " + err.Error()
		return r
	}

	var status bool

	status, err = bd.ModificoRegistro(t, claim.ID.Hex())
	if err != nil {
		r.Message = "Ocurrión un error al intentar modificar el registro. Reintente nuevamente " + err.Error()
		return r
	}

	if !status {
		r.Message = "No se ha logrado modificar el registro del usuario "
		return r
	}

	r.Status = 200
	r.Message = "Modificar Perfil OK !"
	return r
}
