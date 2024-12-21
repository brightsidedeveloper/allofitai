package handler

import (
	"allofitai/pkg/sb"
	"allofitai/pkg/util"
	"allofitai/view/auth"
	"log/slog"
	"net/http"

	"github.com/nedpals/supabase-go"
)

func RenderSignIn(w http.ResponseWriter, r *http.Request) error {
	return render(w, r, auth.SignIn())
}

func RenderCreate(w http.ResponseWriter, r *http.Request) error {
	return render(w, r, auth.Create())
}

func SignIn(w http.ResponseWriter, r *http.Request) error {
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

	setAuthCookies(w, resp.AccessToken)

	cookie, err := r.Cookie("path")
	if err != nil {
		hxRedirect(w, r, "/")
		return nil
	}

	path := cookie.Value
	http.SetCookie(w, &http.Cookie{
		Name:     "path",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	hxRedirect(w, r, path)
	return nil
}

func SignInWithGoogle(w http.ResponseWriter, r *http.Request) error {
	resp, err := sb.Client.Auth.SignInWithProvider(supabase.ProviderSignInOptions{
		Provider:   "google",
		RedirectTo: "http://localhost:8888/auth/callback",
	})
	if err != nil {
		slog.Error("failed to sign in with google", "err", err)
		return render(w, r, auth.SignIn())
	}

	http.Redirect(w, r, resp.URL, http.StatusSeeOther)
	return nil
}

func Create(w http.ResponseWriter, r *http.Request) error {
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

func AuthCallback(w http.ResponseWriter, r *http.Request) error {
	accessToken := r.URL.Query().Get("access_token")
	if len(accessToken) == 0 {
		return render(w, r, auth.CallbackScript())
	}
	setAuthCookies(w, accessToken)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func Logout(w http.ResponseWriter, r *http.Request) error {
	cookie := &http.Cookie{
		Name:     "at",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
		Path:     "/",
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func setAuthCookies(w http.ResponseWriter, accessToken string) {
	cookie := &http.Cookie{
		Name:     "at",
		Value:    accessToken,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	}

	http.SetCookie(w, cookie)
}
