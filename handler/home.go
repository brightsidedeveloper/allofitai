package handler

import (
	"allofitai/view/home"
	"net/http"
)

func RenderHome(w http.ResponseWriter, r *http.Request) error {
	return render(w, r, home.Index())
}
