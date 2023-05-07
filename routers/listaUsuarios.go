package routers

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ptilotta/twittor/bd"
	"github.com/ptilotta/twittor/models"
)

func ListaUsuarios(ctx context.Context, request events.APIGatewayV2HTTPRequest, claim models.Claim) (int, string) {

	page := request.QueryStringParameters["pagina"]
	typeUser := request.QueryStringParameters["type"]
	search := request.QueryStringParameters["search"]

	if len(page) == 0 {
		page = "1"
	}
	pagTemp, err := strconv.Atoi(page)
	if err != nil {
		return 400, "Debe enviar el parámetro página como entero mayor a 0"
	}

	usuarios, status := bd.LeoUsuariosTodos(IDUsuario, int64(pagTemp), search, typeUser)
	if !status {
		return 400, "Error al leer los usuarios"
	}

	respJson, err := json.Marshal(usuarios)
	if err != nil {
		return 500, "Error al formatear los datos de los usuarios como JSON"
	}

	return 200, string(respJson)
}
