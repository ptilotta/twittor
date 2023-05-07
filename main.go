package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/ptilotta/twittor/awsgo"
	"github.com/ptilotta/twittor/bd"
	"github.com/ptilotta/twittor/handlers"
	"github.com/ptilotta/twittor/models"
	"github.com/ptilotta/twittor/secretmanager"
)

func EjecutoLambda(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {

	var res *events.APIGatewayProxyResponse

	awsgo.InicializoAWS()

	if !ValidoParametros() {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error en las variables de entorno. deben incluir 'SecretName'",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	}

	SecretModel, err := secretmanager.GetSecret(os.Getenv("SecretName"))
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error en la lectura de Secret " + err.Error(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	}

	awsgo.Ctx = context.WithValue(awsgo.Ctx, "path", request.RawPath)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, "method", request.RequestContext.HTTP.Method)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, "user", SecretModel.Username)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, "password", SecretModel.Password)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, "host", SecretModel.Host)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, "jwtSign", SecretModel.JWTSign)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, "body", request.Body)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, "header", request.Headers)

	// Chequeo Conexión a la BD o Conecto la BD

	if !bd.BaseConectada() {
		bd.ConectarBD(awsgo.Ctx)
	}
	status, message := handlers.Manejadores(awsgo.Ctx, request)

	if request.RawPath == "login" && status == 200 {
		var t models.RespuestaLogin
		cookie := &http.Cookie{
			Name:    "token",
			Value:   t.Token,
			Expires: time.Now().Add(24 * time.Hour),
		}
		cookieString := cookie.String()
		res = &events.APIGatewayProxyResponse{
			StatusCode: status,
			Body:       string(message),
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
				"Set-Cookie":                  cookieString,
			},
		}
	} else {
		res = &events.APIGatewayProxyResponse{
			StatusCode: status,
			Body:       string(message),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
	}

	return res, nil
}

func main() {
	lambda.Start(EjecutoLambda)
}

func ValidoParametros() bool {
	_, traeParametro := os.LookupEnv("SecretName")
	if !traeParametro {
		return traeParametro
	}
	return traeParametro
}
