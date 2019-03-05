package auth

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	cache "github.com/patrickmn/go-cache"
)

type User struct {
	Name string `json:"name,omitempty"`
	Role bool   `json:"role,omitempty"`
}

type authenticator struct {
	tokens         *cache.Cache
	next           http.Handler
	validateUserID func(int) bool
}

func TokenAuth(seconds time.Duration, handler http.Handler, validateUserID func(int) bool) Authenticator {
	return &authenticator{cache.New(seconds*time.Second, seconds*time.Second), handler, validateUserID}
}

func (a *authenticator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.verify(w, r)
	case http.MethodPost:
		a.authenticate(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (a *authenticator) verify(w http.ResponseWriter, r *http.Request) {
	subprotocols := websocket.Subprotocols(r)
	if len(subprotocols) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token := subprotocols[0]
	userID, ok := a.tokens.Get(token)

	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	userIDstr := strconv.Itoa(userID.(int))
	r.Header.Add("User-Id", userIDstr)
	//get the websocket function handler from the channel
	a.next.ServeHTTP(w, r)
}

func (a *authenticator) authenticate(w http.ResponseWriter, r *http.Request) {
	userIDbytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	userID, err := strconv.Atoi(string(userIDbytes))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	if ok := a.validateUserID(userID); !ok {
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

func readBody(body io.ReadCloser) *User {
	fmt.Println(body)
	if body == nil {
		return nil
	}

	decoder := json.NewDecoder(body)

	var u User
	if err := decoder.Decode(&u); err != nil {
		return nil
	}
	fmt.Println(u)
	return &u
}

func createToken() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
