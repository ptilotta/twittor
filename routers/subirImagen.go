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

	body := ctx.Value(models.Key("body")).(string)
	decodedBody, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		r.Message = "Error decodificando cuerpo Base64: " + err.Error()
		return r
	}

	var filename string
	var usuario models.Usuario

	switch uploadType {
	case "A":
		filename = "avatars/" + claim.ID.Hex() + ".jpg"
		usuario.Avatar = claim.ID.Hex() + ".jpg"
	case "B":
		filename = "banners/" + claim.ID.Hex() + ".jpg"
		usuario.Banner = claim.ID.Hex() + ".jpg"
	}

	svc := s3.NewFromConfig(awsgo.Cfg)

	if err := uploadToS3(ctx, svc, filename, decodedBody, ".jpg"); err != nil {
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

	/*
			    // Leer el archivo JPEG de la solicitud HTTP
		    body := bytes.NewReader([]byte(request.Body))
		    fileBytes, err := ioutil.ReadAll(body)
		    if err != nil {
		        return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
		    }

		    // Crear un objeto de archivo multipart/form-data
		    fileBody := &bytes.Buffer{}
		    writer := multipart.NewWriter(fileBody)
		    part, err := writer.CreateFormFile("file", "image.jpg")
		    if err != nil {
		        return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
		    }
		    part.Write(fileBytes)
		    writer.Close()

		    // Configurar una sesión de AWS
		    sess := session.Must(session.NewSession())
		    svc := s3.New(sess)

		    // Subir el archivo al bucket S3
		    _, err = svc.PutObject(&s3.PutObjectInput{
		        Bucket: aws.String(s3Bucket),
		        Key:    aws.String("image.jpg"),
		        Body:   fileBody,
		    })
		    if err != nil {
		        return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
		    }
	*/
}

func uploadToS3(ctx context.Context, svc *s3.Client, filename string, data []byte, contentType string) error {
	_, err := svc.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(ctx.Value(models.Key("bucketName")).(string)),
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
