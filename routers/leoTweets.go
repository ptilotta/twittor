package routers

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ptilotta/twittor/bd"
	"github.com/ptilotta/twittor/models"
)

func LeoTweets(ctx context.Context, request events.APIGatewayV2HTTPRequest) models.RespApi {

	var r models.RespApi

	ID := request.QueryStringParameters["id"]
	pagina := request.QueryStringParameters["pagina"]

	if len(ID) < 1 {
		r.Message = "El parámetro ID es obligatorio"
		return r
	}

	if len(pagina) < 1 {
		pagina = "1"
	}

	pag, err := strconv.Atoi(pagina)
	if err != nil {
		r.Message = "Debe enviar el parámetro página con un valor mayor a 0"
		return r
	}

	tweets, correcto := bd.LeoTweets(ID, int64(pag))
	if !correcto {
		r.Message = "Error al leer los tweets"
		return r
	}

	respJson, err := json.Marshal(tweets)
	if err != nil {
		r.Status = 500
		r.Message = "Error al formatear los datos de los usuarios como JSON"
		return r
	}

	r.Status = 200
	r.Message = string(respJson)
	return r
}
