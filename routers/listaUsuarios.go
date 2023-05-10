package routers

import (
	"encoding/json"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ptilotta/twittor/bd"
	"github.com/ptilotta/twittor/models"
)

func ListaUsuarios(request events.APIGatewayV2HTTPRequest, claim models.Claim) models.RespApi {

	var r models.RespApi
	r.Status = 400

	page := request.QueryStringParameters["page"]
	typeUser := request.QueryStringParameters["type"]
	search := request.QueryStringParameters["search"]

	if len(page) == 0 {
		page = "1"
	}
	pagTemp, err := strconv.Atoi(page)
	if err != nil {
		r.Message = "Debe enviar el parámetro página como entero mayor a 0"
		return r
	}

	usuarios, status := bd.LeoUsuariosTodos(IDUsuario, int64(pagTemp), search, typeUser)
	if !status {
		r.Message = "Error al leer los usuarios"
		return r
	}

	respJson, err := json.Marshal(usuarios)
	if err != nil {
		r.Status = 500
		r.Message = "Error al formatear los datos de los usuarios como JSON"
		return r
	}

	r.Status = 200
	r.Message = string(respJson)
	return r
}
