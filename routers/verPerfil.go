package routers

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ptilotta/twittor/bd"
)

func VerPerfil(ctx context.Context, request events.APIGatewayV2HTTPRequest) (int, string) {

	ID := request.QueryStringParameters["id"]
	if len(ID) < 1 {
		return 400, "El parámetro ID es obligatorio"
	}

	perfil, err := bd.BuscoPerfil(ID)
	if err != nil {
		return 400, "Ocurrió un error al intentar buscar el registro " + err.Error()
	}

	respJson, err := json.Marshal(perfil)
	if err != nil {
		return 500, "Error al formatear los datos de los usuarios como JSON"
	}

	return 200, string(respJson)
}
