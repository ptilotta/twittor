package models

/*Tweet captura del Body, el mensaje que nos llega */
type Tweet struct {
	Mensaje string `bson:"mensaje" json:"mensaje"`
}
