package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/ptilotta/twittor/awsgo"
	"github.com/ptilotta/twittor/bd"
	"github.com/ptilotta/twittor/handlers"
	"github.com/ptilotta/twittor/secretmanager"
)

func EjecutoLambda(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {

	var res *events.APIGatewayProxyResponse

	awsgo.InicializoAWS()

	if !ValidoParametros() {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error en las variables de entorno. deben incluir 'SecretName', 'BucketName",
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

	path := strings.Replace(request.RawPath, os.Getenv("UrlPrefix"), "", -1)

	awsgo.Ctx = context.WithValue(awsgo.Ctx, "path", path)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, "method", request.RequestContext.HTTP.Method)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, "user", SecretModel.Username)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, "password", SecretModel.Password)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, "host", SecretModel.Host)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, "jwtSign", SecretModel.JWTSign)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, "body", request.Body)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, "header", request.Headers)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, "bucketName", os.Getenv("BucketName"))

	// Chequeo Conexión a la BD o Conecto la BD

	bd.ConectarBD(awsgo.Ctx)

	respAPI := handlers.Manejadores(awsgo.Ctx, request)

	fmt.Println("Sali de Manejadores")
	if respAPI.CustomResp == nil {
		var res *events.APIGatewayProxyResponse
		headersResp := map[string]string{
			"Content-Type": "application/json",
		}
		res = &events.APIGatewayProxyResponse{
			StatusCode: respAPI.Status,
			Body:       string(respAPI.Message),
			Headers:    headersResp,
		}
		return res, nil
	} else {
		return respAPI.CustomResp, nil
	}
}

func main() {
	lambda.Start(EjecutoLambda)
}

func ValidoParametros() bool {
	_, traeParametro := os.LookupEnv("SecretName")
	if !traeParametro {
		return traeParametro
	}
	_, traeParametro = os.LookupEnv("BucketName")
	if !traeParametro {
		return traeParametro
	}
	_, traeParametro = os.LookupEnv("UrlPrefix")
	if !traeParametro {
		return traeParametro
	}

	return traeParametro
}
