package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/adrianbrad/chat/repository"

	"github.com/adrianbrad/chat/auth"
	"github.com/adrianbrad/chat/channel"
	"github.com/adrianbrad/chat/config"
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

var db *sql.DB
var userRepository repository.Repository
var usersChannelsRepository repository.UsersChannelsRepository

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

	// * initDB is called first, then the return value is assigned to the defer
	defer initDB(c.Database)()
	// * removing all previous subscribtions as there is nio chance to recover that in case of an unexpected application shutdown
	db.Exec("TRUNCATE TABLE Users_Rooms")
	userRepository = repository.NewDbUsersRepository(db)

	usersChannelsRepository = repository.NewDbUsersChannelsRepository(db)
	channel := channel.New(usersChannelsRepository, 1)
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
	r.Get("/", redirectToChat)
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

		if !channelExists {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func authRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("auth requests")
		// TODO
		next.ServeHTTP(w, r)
	})
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
