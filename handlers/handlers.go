package handlers

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ptilotta/twittor/models"
	"github.com/ptilotta/twittor/routers"
)

type key models.Key

func Manejadores(ctx context.Context, request events.APIGatewayV2HTTPRequest) models.RespApi {

	fmt.Println("Voy a procesar " + ctx.Value("path").(key) + " > " + ctx.Value("method").(key))

	var r models.RespApi
	r.Status = 400

	isOk, statusCode, msg, claim := validoAuthorization(ctx, request)
	if !isOk {
		r.Status = statusCode
		r.Message = msg
		return r
	}

	switch ctx.Value("method").(key) {
	case "POST":
		switch ctx.Value("path").(key) {
		case "registro":
			return routers.Registro(ctx)
		case "login":
			return routers.Login(ctx)
		case "tweet":
			return routers.GraboTweet(ctx)
		case "altaRelacion":
			return routers.AltaRelacion(ctx, request, claim)
		case "subirAvatar":
			return routers.UploadImage(ctx, "A", request, claim)
		case "subirBanner":
			return routers.UploadImage(ctx, "B", request, claim)
		}
	case "GET":
		switch ctx.Value("path").(key) {
		case "verperfil":
			return routers.VerPerfil(request)
		case "leoTweets":
			return routers.LeoTweets(ctx, request)
		case "consultaRelacion":
			return routers.ConsultaRelacion(ctx, request, claim)
		case "listaUsuarios":
			return routers.ListaUsuarios(ctx, request, claim)
		case "leoTweetsSeguidores":
			return routers.LeoTweetsSeguidores(ctx, request, claim)
		case "obtenerAvatar":
			return routers.ObtenerImagen(ctx, "A", request, claim)
		case "obtenerBanner":
			return routers.ObtenerImagen(ctx, "B", request, claim)
		}
	case "PUT":
		switch ctx.Value("path").(key) {
		case "modificarPerfil":
			return routers.ModificarPerfil(ctx, claim)
		}
	case "DELETE":
		switch ctx.Value("path").(key) {
		case "eliminarTweet":
			return routers.EliminarTweet(ctx, request, claim)
		case "bajaRelacion":
			return routers.BajaRelacion(ctx, request, claim)
		}
	}

	r.Status = 400
	r.Message = "Method Invalid"
	return r
}

func validoAuthorization(ctx context.Context, request events.APIGatewayV2HTTPRequest) (bool, int, string, models.Claim) {

	path := ctx.Value("path").(key)
	if path == "registro" || path == "login" || path == "obtenerAvatar" || path == "obtenerBanner" {
		return true, 200, "", models.Claim{}
	}

	token := request.Headers["authorization"]
	if len(token) == 0 {
		return false, 401, "Token requerido", models.Claim{}
	}

	claim, todoOK, msg, err := routers.ProcesoToken(token, string(ctx.Value("jwtSign").(key)))
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
