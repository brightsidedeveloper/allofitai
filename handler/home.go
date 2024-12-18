package handler

import (
	"allofitai/view/home"
	"net/http"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	return render(w, r, home.Index())
}
