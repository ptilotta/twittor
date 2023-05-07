package routers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ptilotta/twittor/bd"
	"github.com/ptilotta/twittor/models"
)

/*GraboTweet permite grabar el tweet en la base de datos */
func GraboTweet(ctx context.Context) (int, string) {
	var mensaje models.Tweet
	err := json.Unmarshal([]byte(ctx.Value("body").(string)), &mensaje)
	if err != nil {
		return 400, "Ocurrió un error al intentar decodificar el body " + err.Error()
	}

	registro := models.GraboTweet{
		UserID:  IDUsuario,
		Mensaje: mensaje.Mensaje,
		Fecha:   time.Now(),
	}

	_, status, err := bd.InsertoTweet(registro)
	if err != nil {
		return 400, "Ocurrió un error al intentar insertar el registro, reintente nuevamente" + err.Error()
	}

	if !status {
		return 400, "No se ha logrado insertar el Tweet"
	}

	return 200, "Tweet creado correctamente"
}
