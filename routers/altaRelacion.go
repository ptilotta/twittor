package routers

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ptilotta/twittor/bd"
	"github.com/ptilotta/twittor/models"
)

func AltaRelacion(ctx context.Context, request events.APIGatewayV2HTTPRequest, claim models.Claim) models.RespApi {

	var r models.RespApi
	r.Status = 400

	ID := request.QueryStringParameters["id"]
	if len(ID) < 1 {
		r.Message = "El parámetro ID es obligatorio"
		return r
	}

	var t models.Relacion
	t.UsuarioID = claim.ID.Hex()
	t.UsuarioRelacionID = ID

	status, err := bd.InsertoRelacion(t)
	if err != nil {
		r.Message = "Ocurrió un error al intentar insertar relación " + err.Error()
		return r
	}
	if !status {
		r.Message = "No se ha logrado insertar la relación " + err.Error()
		return r
	}

	r.Status = 200
	r.Message = "Alta Relación OK"
	return r
}
