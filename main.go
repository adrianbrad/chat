package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"hash"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/adrianbrad/chat/channel"
	"github.com/adrianbrad/chat/messageProcessor"

	"github.com/adrianbrad/chat/repository"

	"github.com/adrianbrad/chat/auth"
	"github.com/adrianbrad/chat/config"
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

var db *sql.DB
var userRepository repository.UserRepository
var usersChannelsRepository repository.UsersChannelsRepository
var messagesRepository repository.MessageRepository
var secret string
var h hash.Hash
var channelIdentifier string

type templateHandler struct {
	once     sync.Once
	filename string
	template *template.Template
}

// ServerHTTP handles the HTTP request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//we make sure the template initialization is only executed once
	t.once.Do(
		func() {
			t.template = template.Must(
				template.ParseFiles(
					filepath.Join("templates", t.filename),
				),
			)
		},
	)

	data := map[string]interface{}{
		"Host": r.Host,
	}

	_ = t.template.Execute(w, data)
}

func main() {
	c := config.Load("config")
	secret = c.Server.Secret
	channelIdentifier = c.Server.Channel
	h = hmac.New(sha256.New, []byte(secret))
	// * initDB is called first, then the return value is assigned to the defer
	defer initDB(c.Database)()
	// * removing all previous subscribtions as there is nio chance to recover that in case of an unexpected application shutdown
	db.Exec("TRUNCATE TABLE Users_Rooms")

	userRepository = repository.NewDbUsersRepository(db)
	messagesRepository = repository.NewDbMessagesRepository(db)
	usersChannelsRepository = repository.NewDbUsersChannelsRepository(db)

	channel := channel.New(usersChannelsRepository, 1, userRepository, messageProcessor.New(), []int64{1, 2}, messagesRepository)
	go channel.Run() //get the channel going in another thread
	//the chatting operation occur in the background
	//the main goroutine is running the web server

	log.Println("Starting web server on", c.Server.Port)
	err := http.ListenAndServe(c.Server.Port, routes(channel))
	if err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}

func routes(channel channel.Channel) (r *chi.Mux) {
	r = chi.NewRouter()
	// * middlewares
	r.Use(logging)
	r.Use(authRequests)
	// * http.HandleFunc(s)
	// r.Post("/users/")

	// * http.Handler(s)
	r.Method(http.MethodGet, "/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("/home/brad/workspace/go/src/github.com/adrianbrad/chat/assets"))))
	r.Method(http.MethodGet, "/chat", &templateHandler{filename: "chat.html"})
	r.Handle("/talk/{channel}", validChannel(auth.TokenAuth(10, channel, userRepository.CheckIfExists)))
	return r
}

//TODO Add origin verification
func validChannel(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var channelExists bool
		_ = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
		FROM "Channels"
		WHERE "Name"=$1)`, chi.URLParam(r, "channel")).Scan(&channelExists)
		fmt.Println(channelExists)
		if !channelExists {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func authRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ! we ignore this if it's a websocket conn
		if r.Method == http.MethodGet && r.URL.String() == fmt.Sprintf("/talk/%s", channelIdentifier) {
			log.Println("No auth")
			next.ServeHTTP(w, r)
		}
		requestTimeString := r.Header.Get("X-Time")
		timeNow := time.Now().UTC().UnixNano() / int64(time.Millisecond)
		requestTime, err := strconv.ParseInt(requestTimeString, 10, 64)
		if err != nil {
			// w.WriteHeader(http.StatusInternalServerError)
			// return
		}

		// ! we accept only requests from 10 maximum 10 seconds ago
		if (timeNow-requestTime)/1000 > 10 {
			log.Println("wrong timestamp")
			// w.WriteHeader(http.StatusUnauthorized)
			// return
		}

		if calculateHash([]byte(requestTimeString)) != r.Header.Get("X-Authorization") {
			log.Println("Unauthorized request")
			// w.WriteHeader(http.StatusUnauthorized)
			// return
		}
		next.ServeHTTP(w, r)
	})
}

func calculateHash(data []byte) string {
	h.Reset()
	h.Write(data)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func redirectToChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "/chat")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func initDB(config config.DatabaseConfiguration) func() error {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port,
		config.User, config.Pass, config.Name)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to the database!")

	return db.Close
}
