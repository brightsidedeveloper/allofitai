package handler

import (
	"allofitai/pkg/sb"
	"allofitai/pkg/util"
	"allofitai/view/auth"
	"log/slog"
	"net/http"

	"github.com/nedpals/supabase-go"
)

func HandleSignInIndex(w http.ResponseWriter, r *http.Request) error {
	return render(w, r, auth.SignIn())
}

func HandleCreateIndex(w http.ResponseWriter, r *http.Request) error {
	return render(w, r, auth.Create())
}

func HandleSignIn(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	if ok := util.IsValidEmail(credentials.Email); !ok {
		return render(w, r, auth.AuthForm("/signin", credentials, auth.AuthErrors{
			Email: "Enter a Valid Email Address",
		}))
	}

	if ok := util.IsValidPassword(credentials.Password); !ok {
		return render(w, r, auth.AuthForm("/signin", credentials, auth.AuthErrors{
			Password: "Password must be at least 6 characters",
		}))
	}

	resp, err := sb.Client.Auth.SignIn(r.Context(), credentials)
	if err != nil {
		slog.Error("failed to sign in", "err", err)
		return render(w, r, auth.AuthForm("/signin", credentials, auth.AuthErrors{
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

func HandleCreate(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	if ok := util.IsValidEmail(credentials.Email); !ok {
		return render(w, r, auth.AuthForm("/create", credentials, auth.AuthErrors{
			Email: "Enter a Valid Email Address",
		}))
	}

	if ok := util.IsValidPassword(credentials.Password); !ok {
		return render(w, r, auth.AuthForm("/create", credentials, auth.AuthErrors{
			Password: "Password must be at least 6 characters",
		}))
	}

	sbUser, err := sb.Client.Auth.SignUp(r.Context(), credentials)
	if err != nil {
		slog.Error("failed to sign up", "err", err)
		return render(w, r, auth.AuthForm("/create", credentials, auth.AuthErrors{
			InvalidCredentials: "Supabase ain't happy",
		}))
	}

	return render(w, r, auth.CreateSuccess(sbUser.Email))
}
