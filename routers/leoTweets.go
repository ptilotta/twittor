package routers

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ptilotta/twittor/bd"
)

func LeoTweets(ctx context.Context, request events.APIGatewayV2HTTPRequest) (int, string) {

	ID := request.QueryStringParameters["id"]
	pagina := request.QueryStringParameters["pagina"]

	if len(ID) < 1 {
		return 400, "El parámetro ID es obligatorio"
	}

	if len(pagina) < 1 {
		pagina = "1"
	}

	pag, err := strconv.Atoi(pagina)
	if err != nil {
		return 400, "Debe enviar el parámetro página con un valor mayor a 0"
	}

	tweets, correcto := bd.LeoTweets(ID, int64(pag))
	if !correcto {
		return 400, "Error al leer los tweets"
	}

	respJson, err := json.Marshal(tweets)
	if err != nil {
		return 500, "Error al formatear los datos de los usuarios como JSON"
	}

	return 200, string(respJson)
}
