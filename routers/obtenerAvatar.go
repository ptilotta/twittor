package routers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ptilotta/twittor/bd"
)

/*ObtenerAvatar envia el Avatar al HTTP */
func ObtenerAvatar(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Entre en Obtener 1")
	ID := r.URL.Query().Get("id")
	fmt.Println("Entre en Obtener 1 a")
	if len(ID) < 1 {
		fmt.Println("Entre en Obtener 1 b")
		http.Error(w, "Debe enviar el parámetro ID", http.StatusBadRequest)
		fmt.Println("Entre en Obtener 1 c")
		return
	}

	fmt.Println("Entre en Obtener 2")
	perfil, err := bd.BuscoPerfil(ID)
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusBadRequest)
		return
	}

	fmt.Println("Entre en Obtener 3")
	OpenFile, err := os.Open("uploads/avatars/" + perfil.Avatar)
	if err != nil {
		http.Error(w, "Imagen no encontrada", http.StatusBadRequest)
		return
	}

	fmt.Println("Entre en Obtener 4")
	_, err = io.Copy(w, OpenFile)
	if err != nil {
		http.Error(w, "Error al copiar la imagen", http.StatusBadRequest)
	}
}
