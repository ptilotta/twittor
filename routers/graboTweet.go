package routers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ptilotta/twittor/bd"
	"github.com/ptilotta/twittor/models"
)

func GraboTweet(ctx context.Context) models.RespApi {
	var mensaje models.Tweet
	var r models.RespApi
	r.Status = 400

	err := json.Unmarshal([]byte(string(ctx.Value("body").(models.Key))), &mensaje)
	if err != nil {
		r.Message = "Ocurrió un error al intentar decodificar el body " + err.Error()
		return r
	}

	registro := models.GraboTweet{
		UserID:  IDUsuario,
		Mensaje: mensaje.Mensaje,
		Fecha:   time.Now(),
	}

	_, status, err := bd.InsertoTweet(registro)
	if err != nil {
		r.Message = "Ocurrió un error al intentar insertar el registro, reintente nuevamente" + err.Error()
		return r
	}

	if !status {
		r.Message = "No se ha logrado insertar el Tweet"
		return r
	}

	r.Status = 200
	r.Message = "Tweet creado correctamente"
	return r
}
