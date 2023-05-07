package routers

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ptilotta/twittor/bd"
	"github.com/ptilotta/twittor/models"
)

func AltaRelacion(ctx context.Context, request events.APIGatewayV2HTTPRequest, claim models.Claim) (int, string) {

	ID := request.QueryStringParameters["id"]
	if len(ID) < 1 {
		return 400, "El parámetro ID es obligatorio"
	}

	var t models.Relacion
	t.UsuarioID = claim.ID.Hex()
	t.UsuarioRelacionID = ID

	status, err := bd.InsertoRelacion(t)
	if err != nil {
		return 400, "Ocurrió un error al intentar insertar relación " + err.Error()
	}
	if !status {
		return 400, "No se ha logrado insertar la relación " + err.Error()
	}
	return 200, "Alta Relación OK"
}
