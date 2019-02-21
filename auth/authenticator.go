package auth

import "net/http"

type Authenticator interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
