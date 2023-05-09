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

type key models.Key

func ConectarBD(ctx context.Context) error {

	user := string(ctx.Value("user").(key))
	passwd := string(ctx.Value("password").(key))
	host := string(ctx.Value("host").(key))
	connStr := fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority", user, passwd, host)
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
	DatabaseName = string(ctx.Value("database").(key))
	return nil
}

func BaseConectada() bool {
	err := MongoCN.Ping(context.TODO(), nil)
	return err == nil
}
