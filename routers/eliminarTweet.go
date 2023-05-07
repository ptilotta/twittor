package routers

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ptilotta/twittor/bd"
	"github.com/ptilotta/twittor/models"
)

func EliminarTweet(ctx context.Context, request events.APIGatewayV2HTTPRequest, claim models.Claim) (int, string) {

	ID := request.QueryStringParameters["id"]
	if len(ID) < 1 {
		return 400, "El parámetro ID es obligatorio"
	}

	err := bd.BorroTweet(ID, claim.ID.Hex())
	if err != nil {
		return 400, "Ocurrió un error al intentar borrar el tweet " + err.Error()
	}

	return 200, "Eliminar Tweet OK !"
}
