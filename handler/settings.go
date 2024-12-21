package handler

import (
	"allofitai/view/settings"
	"net/http"
)

func RenderSettings(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	return render(w, r, settings.Index(user))
}
