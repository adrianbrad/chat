package auth

import "net/http"

//Authenthicator implements http.Handler
type Authenticator interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
