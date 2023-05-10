package bd

import (
	"context"
	"fmt"
	"time"

	"github.com/ptilotta/twittor/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
LeoUsuariosTodos Lee los usuarios registrados en el sistema, si se recibe "R" en quienes

	trae solo los que se relacionan conmigo
*/
func LeoUsuariosTodos(ID string, page int64, search string, tipo string) ([]*models.Usuario, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoCN.Database(DatabaseName)
	col := db.Collection("usuarios")

	var results []*models.Usuario

	findOptions := options.Find()
	findOptions.SetLimit(20)
	findOptions.SetSkip((page - 1) * 20)

	query := bson.M{
		"nombre": bson.M{"$regex": `(?i)` + search},
	}

	cur, err := col.Find(ctx, query, findOptions)
	if err != nil {
		return results, false
	}

	var incluir bool

	for cur.Next(ctx) {
		var s models.Usuario
		err := cur.Decode(&s)
		if err != nil {
			fmt.Println("Decode = " + err.Error())
			return results, false
		}

		var r models.Relacion
		r.UsuarioID = ID
		r.UsuarioRelacionID = s.ID.Hex()

		incluir = false

		encontrado := ConsultoRelacion(r)
		fmt.Println("===================================")
		fmt.Println("encontrado = ")
		fmt.Println(encontrado)
		fmt.Println("===================================")

		if tipo == "new" && !encontrado {
			fmt.Println("tipo == 'new' && !encontrado")
			incluir = true
		}
		if tipo == "follow" && encontrado {
			fmt.Println("tipo == 'follow' && encontrado")
			incluir = true
		}

		if r.UsuarioRelacionID == ID {
			fmt.Println("r.UsuarioRelacionID == ID")
			incluir = false
		}

		fmt.Println("===================================")
		fmt.Println("incluir = ")
		fmt.Println(incluir)
		fmt.Println("===================================")
		if incluir {
			s.Password = ""
			s.Biografia = ""
			s.SitioWeb = ""
			s.Ubicacion = ""
			s.Banner = ""
			s.Email = ""
			results = append(results, &s)
		}
	}

	err = cur.Err()
	if err != nil {
		fmt.Println("cur.Err() = " + err.Error())
		return results, false
	}
	cur.Close(ctx)
	return results, true
}
