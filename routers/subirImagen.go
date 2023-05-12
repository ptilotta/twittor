package routers

import (
	"bytes"
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/ptilotta/twittor/awsgo"
	"github.com/ptilotta/twittor/bd"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ptilotta/twittor/models"
)

/*
en el bucket S3 creado hay que habilitar la política del bucket y colocar el siguiente JSON

{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": "*",
            "Action": [
                "s3:GetObject",
                "s3:ListBucket"
            ],
            "Resource": [
                "arn:aws:s3:::ecomgo",
                "arn:aws:s3:::ecomgo/*"
            ]
        }
    ]
}
*/

func UploadImage(ctx context.Context, uploadType string, request events.APIGatewayV2HTTPRequest, claim models.Claim) models.RespApi {

	var r models.RespApi
	r.Status = 400

	var filename string

	fileBytes := []byte(ctx.Value(models.Key("body")).(string))

	var usuario models.Usuario

	switch uploadType {
	case "A":
		filename = "avatars/" + claim.ID.Hex() + ".jpg"
		usuario.Avatar = claim.ID.Hex() + ".jpg"
	case "B":
		filename = "banners/" + claim.ID.Hex() + ".jpg"
		usuario.Banner = claim.ID.Hex() + ".jpg"
	}

	// Crear un objeto de archivo multipart/form-data
	/*	fileBody := &bytes.Buffer{}
		writer := multipart.NewWriter(fileBody)
		part, err := writer.CreateFormFile("file", filename)
		if err != nil {
			r.Status = 500
			r.Message = "Error realizando el CreateFormFile: " + err.Error()
			return r
		}

		part.Write(fileBytes)
		writer.Close()
	*/
	svc := s3.NewFromConfig(awsgo.Cfg)
	// Subir el archivo al bucket S3
	_, err := svc.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(ctx.Value(models.Key("bucketName")).(string)),
		Key:         aws.String(filename),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String("image/jpeg"),
	})

	if err != nil {
		r.Status = 500
		r.Message = err.Error()
		return r
	}

	var status bool

	status, err = bd.ModificoRegistro(usuario, IDUsuario)
	if err != nil || !status {
		r.Message = "Error al grabar el avatar en la BD ! " + err.Error()
		return r
	}

	r.Status = 200
	r.Message = "Archivo subido correctamente"

	return r

}
