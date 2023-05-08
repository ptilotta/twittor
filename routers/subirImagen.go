package routers

import (
	"context"
	"encoding/base64"
	"strings"

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

	decodedBody, err := base64.StdEncoding.DecodeString(ctx.Value("body").(string))
	if err != nil {
		r.Message = "Error decodificando cuerpo Base64: " + err.Error()
		return r
	}

	ext := getFileExtension(request.Headers["Content-Type"])
	if ext == "" {
		r.Message = "Error obteniendo extensión del archivo"
		return r
	}

	var filename string
	var usuario models.Usuario

	switch uploadType {
	case "A":
		filename = "avatars/" + claim.ID.Hex() + ext
		usuario.Avatar = claim.ID.Hex() + ext
	case "B":
		filename = "banners/" + claim.ID.Hex() + ext
		usuario.Banner = claim.ID.Hex() + ext
	}

	svc := s3.NewFromConfig(awsgo.Cfg)

	if err := uploadToS3(ctx, svc, filename, decodedBody, request.Headers["Content-Type"]); err != nil {
		r.Status = 500
		r.Message = "Error cargando archivo a S3: " + err.Error()
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

func getFileExtension(contentType string) string {
	switch {
	case strings.Contains(contentType, "jpeg"):
		return ".jpg"
	case strings.Contains(contentType, "png"):
		return ".png"
	default:
		return ""
	}
}

func uploadToS3(ctx context.Context, svc *s3.Client, filename string, data []byte, contentType string) error {
	_, err := svc.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(ctx.Value("bucketName").(string)),
		Key:           aws.String(filename),
		Body:          aws.ReadSeekCloser(strings.NewReader(string(data))),
		ContentType:   aws.String(contentType),
		ContentLength: *aws.Int64(int64(len(data))),
	})
	if err != nil {
		return err
	}
	return nil
}
