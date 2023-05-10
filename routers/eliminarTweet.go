package routers

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/ptilotta/twittor/bd"
	"github.com/ptilotta/twittor/models"
)

func EliminarTweet(request events.APIGatewayV2HTTPRequest, claim models.Claim) models.RespApi {

	var r models.RespApi
	r.Status = 400

	ID := request.QueryStringParameters["id"]
	if len(ID) < 1 {
		r.Message = "El parámetro ID es obligatorio"
		return r
	}

	err := bd.BorroTweet(ID, claim.ID.Hex())
	if err != nil {
		r.Message = "Ocurrió un error al intentar borrar el tweet " + err.Error()
		return r
	}

	r.Status = 200
	r.Message = "Eliminar Tweet OK !"
	return r
}
