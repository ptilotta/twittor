package bd

import (
	"context"
	"fmt"

	"github.com/ptilotta/twittor/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*MongoCN es el objeto de conexión a la BD */
var MongoCN *mongo.Client
var DatabaseName string

func ConectarBD(ctx context.Context) error {

	// sss
	user := ctx.Value("user").(models.Key)
	passwd := ctx.Value("password").(models.Key)
	host := ctx.Value("host").(models.Key)
	connStr := fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority", string(user), string(passwd), string(host))
	fmt.Println(connStr)
	var clientOptions = options.Client().ApplyURI(connStr)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("Conexión Exitosa con la BD")
	MongoCN = client
	db := ctx.Value("database").(models.Key)
	DatabaseName = string(db)
	return nil
}

func BaseConectada() bool {
	err := MongoCN.Ping(context.TODO(), nil)
	return err == nil
}
