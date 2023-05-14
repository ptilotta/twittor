package handlers

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ptilotta/twittor/jwt"
	"github.com/ptilotta/twittor/models"
	"github.com/ptilotta/twittor/routers"
)

func Manejadores(ctx context.Context, request events.APIGatewayProxyRequest) models.RespApi {

	fmt.Println("Voy a procesar " + ctx.Value(models.Key("path")).(string) + " > " + ctx.Value(models.Key("method")).(string))

	var r models.RespApi
	r.Status = 400

	isOk, statusCode, msg, claim := validoAuthorization(ctx, request)
	if !isOk {
		r.Status = statusCode
		r.Message = msg
		return r
	}

	switch ctx.Value(models.Key("method")).(string) {
	case "POST":
		switch ctx.Value(models.Key("path")).(string) {
		case "registro": // listo
			return routers.Registro(ctx)
		case "login": // listo
			return routers.Login(ctx)
		case "tweet": // listo
			return routers.GraboTweet(ctx, claim)
		case "altaRelacion": // listo
			return routers.AltaRelacion(ctx, request, claim)
		case "subirAvatar": // listo
			return routers.UploadImage(ctx, "A", request, claim)
		case "subirBanner": // listo
			return routers.UploadImage(ctx, "B", request, claim)
		}
	case "GET":
		switch ctx.Value(models.Key("path")).(string) {
		case "verperfil": // listo
			return routers.VerPerfil(request)
		case "leoTweets": // listo
			return routers.LeoTweets(request)
		case "consultaRelacion": // listo
			return routers.ConsultaRelacion(request, claim)
		case "listaUsuarios": // listo
			return routers.ListaUsuarios(request, claim)
		case "leoTweetsSeguidores": // listo
			return routers.LeoTweetsSeguidores(request, claim)
		case "obtenerAvatar": // listo
			return routers.ObtenerImagen(ctx, "A", request, claim)
		case "obtenerBanner": // listo
			return routers.ObtenerImagen(ctx, "B", request, claim)
		}
	case "PUT":
		switch ctx.Value(models.Key("path")).(string) {
		case "modificarPerfil": // listo
			return routers.ModificarPerfil(ctx, claim)
		}
	case "DELETE":
		switch ctx.Value(models.Key("path")).(string) {
		case "eliminarTweet": // listo
			return routers.EliminarTweet(request, claim)
		case "bajaRelacion": // listo
			return routers.BajaRelacion(request, claim)
		}
	}

	r.Status = 400
	r.Message = "Method Invalid"
	return r
}

func validoAuthorization(ctx context.Context, request events.APIGatewayProxyRequest) (bool, int, string, models.Claim) {

	path := ctx.Value(models.Key("path")).(string)
	if path == "registro" || path == "login" || path == "obtenerAvatar" || path == "obtenerBanner" {
		return true, 200, "", models.Claim{}
	}

	token := request.Headers["Authorization"]
	if len(token) == 0 {
		return false, 401, "Token requerido", models.Claim{}
	}

	claim, todoOK, msg, err := jwt.ProcesoToken(token, ctx.Value(models.Key("jwtSign")).(string))
	if !todoOK {
		if err != nil {
			fmt.Println("Error en el token " + err.Error())
			return false, 401, err.Error(), models.Claim{}
		} else {
			fmt.Println("Error en el token " + msg)
			return false, 401, msg, models.Claim{}
		}
	}

	fmt.Println("Token OK")
	return true, 200, msg, *claim
}
