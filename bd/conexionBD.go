package bd

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*MongoCN es el objeto de conexión a la BD */
var MongoCN *mongo.Client
var DatabaseName string

func ConectarBD(ctx context.Context) error {
	connStr := fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority", ctx.Value("user").(string), ctx.Value("password").(string), ctx.Value("host").(string))
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
	DatabaseName = ctx.Value("database").(string)
	return nil
}

func BaseConectada() bool {
	err := MongoCN.Ping(context.TODO(), nil)
	return err == nil
}
