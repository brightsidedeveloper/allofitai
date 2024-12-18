package handler

import (
	"allofitai/pkg/sb"
	"allofitai/pkg/util"
	"allofitai/view/auth"
	"log/slog"
	"net/http"

	"github.com/nedpals/supabase-go"
)

func HandleAuthIndex(w http.ResponseWriter, r *http.Request) error {
	return render(w, r, auth.SignIn())
}

func HandleSignInCreate(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	if ok := util.IsValidEmail(credentials.Email); !ok {
		return render(w, r, auth.SignInForm(credentials, auth.SignInErrors{
			Email: "Enter a Valid Email Address",
		}))
	}

	if ok := util.IsValidPassword(credentials.Password); !ok {
		return render(w, r, auth.SignInForm(credentials, auth.SignInErrors{
			Password: "Password must be at least 6 characters",
		}))
	}

	resp, err := sb.Client.Auth.SignIn(r.Context(), credentials)
	if err != nil {
		slog.Error("failed to sign in", "err", err)
		return render(w, r, auth.SignInForm(credentials, auth.SignInErrors{
			InvalidCredentials: "Invalid Email or Password",
		}))
	}

	cookie := &http.Cookie{
		Name:     "at",
		Value:    resp.AccessToken,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	}

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}
