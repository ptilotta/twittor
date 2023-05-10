package routers

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/ptilotta/twittor/bd"
	"github.com/ptilotta/twittor/models"
)

func BajaRelacion(request events.APIGatewayV2HTTPRequest, claim models.Claim) models.RespApi {

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

	status, err := bd.BorroRelacion(t)
	if err != nil {
		r.Message = "Ocurrió un error al intentar borrar relación " + err.Error()
		return r
	}
	if !status {
		r.Message = "No se ha logrado borrar la relación " + err.Error()
		return r
	}

	r.Status = 200
	r.Message = "Baja Relación OK!"
	return r
}
