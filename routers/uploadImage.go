package routers

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"

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
type readSeeker struct {
	io.Reader
}

func (rs *readSeeker) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func UploadImage(ctx context.Context, uploadType string, request events.APIGatewayProxyRequest, claim models.Claim) models.RespApi {

	var r models.RespApi
	r.Status = 400

	var filename string
	var usuario models.Usuario

	bucket := aws.String(ctx.Value(models.Key("bucketName")).(string))
	switch uploadType {
	case "A":
		filename = "avatars/" + claim.ID.Hex() + ".jpg"
		usuario.Avatar = claim.ID.Hex() + ".jpg"
	case "B":
		filename = "banners/" + claim.ID.Hex() + ".jpg"
		usuario.Banner = claim.ID.Hex() + ".jpg"
	}

	fmt.Println("paso 1")
	fmt.Println(request.Headers)
	mediaType, params, err := mime.ParseMediaType(request.Headers["Content-Type"])
	if err != nil {
		r.Status = 500
		r.Message = err.Error()
		return r
	}

	fmt.Println("paso 2")
	if strings.HasPrefix(mediaType, "multipart/") {
		body, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			r.Status = 500
			r.Message = err.Error()
			return r
		}
		fmt.Println("paso 3")

		mr := multipart.NewReader(bytes.NewReader(body), params["boundary"])
		p, err := mr.NextPart()
		if err != nil && err != io.EOF {
			r.Status = 500
			r.Message = err.Error()
			return r
		}
		if err != io.EOF {
			if p.FileName() != "" {
				buf := bytes.NewBuffer(nil)
				if _, err := io.Copy(buf, p); err != nil {
					r.Status = 500
					r.Message = err.Error()
					return r
				}
				sess, err := session.NewSession(&aws.Config{
					Region: aws.String("us-east-1")},
				)

				if err != nil {
					r.Status = 500
					r.Message = err.Error()
					return r
				}

				uploader := s3manager.NewUploader(sess)
				_, err = uploader.Upload(&s3manager.UploadInput{
					Bucket: bucket,
					Key:    aws.String(filename),
					Body:   &readSeeker{buf},
				})

				if err != nil {
					r.Status = 500
					r.Message = err.Error()
					return r
				}

			}
		}

	} else {
		r.Status = 400
		r.Message = "Debe enviar una imagen con el 'Content-Type' de tipo 'multipart/' en el Header"
		return r
	}

	r.Status = 200
	r.Message = "Image Upload OK"
	return r
}
