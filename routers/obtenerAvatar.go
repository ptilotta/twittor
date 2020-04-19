package routers

import (
	"io"
	"net/http"
	"os"

	"github.com/ptilotta/twittor/bd"
)

/*ObtenerAvatar envia el Avatar al HTTP */
func ObtenerAvatar(w http.ResponseWriter, r *http.Request) {

	ID := r.URL.Query().Get("id")
	if len(ID) < 1 {
		http.Error(w, "Debe enviar el parámetro ID", http.StatusBadRequest)
		return
	}

	perfil, err := bd.BuscoPerfil(ID)
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusBadRequest)
		return
	}

	OpenFile, err := os.Open("uploads/avatars/" + perfil.Avatar)
	if err != nil {
		http.Error(w, "Imagen no encontrada", http.StatusBadRequest)
		return
	}

	_, err = io.Copy(w, OpenFile)
	if err != nil {
		http.Error(w, "Error al copiar la imagen", http.StatusBadRequest)
	}
}
