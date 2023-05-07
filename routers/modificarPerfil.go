package routers

import (
	"context"
	"encoding/json"

	"github.com/ptilotta/twittor/bd"
	"github.com/ptilotta/twittor/models"
)

/*ModificarPerfil modifica el perfil de usuario */
func ModificarPerfil(ctx context.Context, claim models.Claim) (int, string) {

	var t models.Usuario

	err := json.Unmarshal([]byte(ctx.Value("body").(string)), &t)
	if err != nil {
		return 400, "Datos Incorrectos " + err.Error()
	}

	var status bool

	status, err = bd.ModificoRegistro(t, claim.ID.Hex())
	if err != nil {
		return 400, "Ocurrión un error al intentar modificar el registro. Reintente nuevamente " + err.Error()
	}

	if !status {
		return 400, "No se ha logrado modificar el registro del usuario "
	}

	return 200, "Modificar Perfil OK !"
}
