package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/objx"

	"github.com/stretchr/gomniauth"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")
	if err == http.ErrNoCookie || cookie.Value == "" {
		//user is not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		//other errors
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//success
	h.next.ServeHTTP(w, r)
}

//MustAuth helper function creates an authHandler that wraps any other http.Handler
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// login handler handles the third-party login process.
// format: /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")

	if len(segs) != 4 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Malicious path: %s", r.URL.Path)
		return
	}

	action := segs[2]

	provider := segs[3]

	fmt.Println(action)
	fmt.Println(provider)

	switch action {
	case "login":
		log.Println("TODO handle login for ", provider)
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("Error when trying to get provider %s: %s", provider, err), http.StatusBadRequest)
			return
		}
		loginURL, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("Error when trying to GetBeginAuthURL for %s: %s", provider, err),
				http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", loginURL)
		w.WriteHeader(http.StatusTemporaryRedirect)

	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("Error when trying to get provider %s: %s", provider, err),
				http.StatusInternalServerError)
			return
		}

		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			http.Error(w,
				fmt.Sprintf("Error when trying to complete auth for %s: %s", provider, err),
				http.StatusInternalServerError)
			return
		}

		user, err := provider.GetUser(creds)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("Error when trying to get user from %s: %s", provider, err),
				http.StatusInternalServerError)
			return
		}

		authCookoeValue := objx.New(map[string]interface{}{
			"name":       user.Name(),
			"avatar_url": user.AvatarURL(),
		}).MustBase64()

		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookoeValue,
			Path:  "/",
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}
