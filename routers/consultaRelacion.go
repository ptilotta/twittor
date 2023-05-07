package routers

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ptilotta/twittor/bd"
	"github.com/ptilotta/twittor/models"
)

func ConsultaRelacion(ctx context.Context, request events.APIGatewayV2HTTPRequest, claim models.Claim) (int, string) {

	ID := request.QueryStringParameters["id"]
	if len(ID) < 1 {
		return 400, "El parámetro ID es obligatorio"
	}

	var t models.Relacion
	t.UsuarioID = claim.ID.Hex()
	t.UsuarioRelacionID = ID

	var resp models.RespuestaConsultaRelacion

	status, err := bd.ConsultoRelacion(t)
	if err != nil || !status {
		resp.Status = false
	} else {
		resp.Status = true
	}

	respJson, err := json.Marshal(status)
	if err != nil {
		return 500, "Error al formatear los datos de los usuarios como JSON"
	}

	return 200, string(respJson)
}
