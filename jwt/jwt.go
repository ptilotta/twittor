package jwt

import (
	"context"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/ptilotta/twittor/models"
)

/*GeneroJWT genera el encriptado con JWT */
func GeneroJWT(ctx context.Context, t models.Usuario) (string, error) {

	//miClave := []byte("MastersdelDesarrollo_grupodeFacebook")
	miClave := []byte(ctx.Value("jwtSign").(string))

	payload := jwt.MapClaims{
		"email":            t.Email,
		"nombre":           t.Nombre,
		"apellidos":        t.Apellidos,
		"fecha_nacimiento": t.FechaNacimiento,
		"biografia":        t.Biografia,
		"ubicacion":        t.Ubicacion,
		"sitioweb":         t.SitioWeb,
		"_id":              t.ID.Hex(),
		"exp":              time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenStr, err := token.SignedString(miClave)
	if err != nil {
		return tokenStr, err
	}
	return tokenStr, nil
}
