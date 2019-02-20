package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	cache "github.com/patrickmn/go-cache"
)

type user struct {
	Name string `json:"name,omitempty"`
	Role bool   `json:"role,omitempty"`
}

type authenticator struct {
	tokens *cache.Cache
}

func (a *authenticator) authenticate(w http.ResponseWriter, r *http.Request) {
	userIDbytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := string(userIDbytes)

	if ok := validateUserID(userID); !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token := createToken()

	if err = a.tokens.Add(token, userID, cache.DefaultExpiration); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Authorization", token)
}

func readBody(body io.ReadCloser) *user {
	if body == nil {
		return nil
	}

	decoder := json.NewDecoder(body)

	var u user
	if err := decoder.Decode(&u); err != nil {
		return nil
	}

	return &u
}

func createToken() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func validateUserID(userID string) bool {
	if _, ok := users[userID]; ok {
		return true
	}
	return false
}
