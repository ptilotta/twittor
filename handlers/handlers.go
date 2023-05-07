package handlers

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ptilotta/twittor/models"
	"github.com/ptilotta/twittor/routers"
)

/*Manejadores seteo mi puerto, el Handler y pongo a escuchar al Servidor */
func Manejadores(ctx context.Context, request events.APIGatewayV2HTTPRequest) (int, string) {

	/*	router.HandleFunc("/subirAvatar", middlew.ChequeoBD(middlew.ValidoJWT(routers.SubirAvatar))).Methods("POST")
		router.HandleFunc("/obtenerAvatar", middlew.ChequeoBD(routers.ObtenerAvatar)).Methods("GET")
		router.HandleFunc("/subirBanner", middlew.ChequeoBD(middlew.ValidoJWT(routers.SubirBanner))).Methods("POST")
		router.HandleFunc("/obtenerBanner", middlew.ChequeoBD(routers.ObtenerBanner)).Methods("GET")
	*/
	//================================================================================================
	fmt.Println("Voy a procesar " + ctx.Value("path").(string) + " > " + ctx.Value("method").(string))

	isOk, statusCode, msg, claim := validoAuthorization(ctx)
	if !isOk {
		return statusCode, msg
	}

	switch ctx.Value("method").(string) {
	case "POST":
		switch ctx.Value("path").(string) {
		case "registro":
			return routers.Registro(ctx)
		case "login":
			return routers.Login(ctx)
		case "tweet":
			return routers.GraboTweet(ctx)
		case "altaRelacion":
			return routers.AltaRelacion(ctx, request, claim)
		}
	case "GET":
		switch ctx.Value("path").(string) {
		case "verperfil":
			return routers.VerPerfil(ctx, request)
		case "leoTweets":
			return routers.LeoTweets(ctx, request)
		case "consultaRelacion":
			return routers.ConsultaRelacion(ctx, request, claim)
		case "listaUsuarios":
			return routers.ListaUsuarios(ctx, request, claim)
		case "leoTweetsSeguidores":
			return routers.LeoTweetsSeguidores(ctx, request, claim)
		}
	case "PUT":
		switch ctx.Value("path").(string) {
		case "modificarPerfil":
			return routers.ModificarPerfil(ctx, claim)
		}
	case "DELETE":
		switch ctx.Value("path").(string) {
		case "eliminarTweet":
			return routers.EliminarTweet(ctx, request, claim)
		case "bajaRelacion":
			return routers.BajaRelacion(ctx, request, claim)
		}
	}

	return 400, "Method Invalid"

}

func validoAuthorization(ctx context.Context) (bool, int, string, models.Claim) {

	path := ctx.Value("path").(string)
	if path == "registro" || path == "login" || path == "obtenerAvatar" || path == "obtenerBanner" {
		return true, 200, "", models.Claim{}
	}

	head := ctx.Value("headers").(map[string]string)
	token := head["authorization"]
	if len(token) == 0 {
		return false, 401, "Token requerido", models.Claim{}
	}

	claim, todoOK, msg, err := routers.ProcesoToken(token, ctx.Value("jwtSign").(string))
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
