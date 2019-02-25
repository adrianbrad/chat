package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/go-chi/render"

	"github.com/adrianbrad/chat/config"
	"github.com/adrianbrad/chat/model"
	"github.com/adrianbrad/chat/repository"
	"github.com/adrianbrad/chat/room"
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var db *sql.DB

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
	config := loadConfig()

	// * initDB is called first, then the return value is assigned to the defer
	defer initDB(config.Database)()

	r := chi.NewRouter()
	r.Use(logging)

	r.Method(http.MethodGet, "/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("/home/brad/workspace/go/src/github.com/adrianbrad/chat/assets"))))

	r.Method(http.MethodGet, "/chat", &templateHandler{filename: "chat.html"})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	repo := repository.NewDbUsersRepository(db)
	r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		res, _ := repo.GetOne(id)
		render.JSON(w, r, res)
	})
	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		res := repo.GetAll()
		render.JSON(w, r, res)
	})
	r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		var u model.User
		render.DecodeJSON(r.Body, &u)
		id, err := repo.Create(u)
		fmt.Println(err)
		render.JSON(w, r, id)
	})

	roomRepo := repository.NewDbRoomsRepository(db)
	r.Get("/rooms/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		res, _ := roomRepo.GetOne(id)
		render.JSON(w, r, res)
	})
	r.Get("/rooms", func(w http.ResponseWriter, r *http.Request) {
		res := roomRepo.GetAll()
		render.JSON(w, r, res)
	})
	r.Post("/rooms", func(w http.ResponseWriter, r *http.Request) {
		var u model.Room
		render.DecodeJSON(r.Body, &u)
		id, err := roomRepo.Create(u)
		fmt.Println(err)
		render.JSON(w, r, id)
	})

	messageRepo := repository.NewDbMessagesRepository(db)
	r.Get("/messages/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		res, _ := messageRepo.GetOne(id)
		render.JSON(w, r, res)
	})
	r.Get("/messages", func(w http.ResponseWriter, r *http.Request) {
		res := messageRepo.GetAll()
		render.JSON(w, r, res)
	})
	r.Post("/messages", func(w http.ResponseWriter, r *http.Request) {
		var u model.Message
		render.DecodeJSON(r.Body, &u)
		id, err := messageRepo.Create(u)
		fmt.Println(err)
		render.JSON(w, r, id)
	})

	roleRepo := repository.NewDbRolesRepository(db)
	r.Get("/roles/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		res, _ := roleRepo.GetOne(id)
		render.JSON(w, r, res)
	})
	r.Get("/roles", func(w http.ResponseWriter, r *http.Request) {
		res := roleRepo.GetAll()
		render.JSON(w, r, res)
	})
	r.Post("/roles", func(w http.ResponseWriter, r *http.Request) {
		var u model.Role
		render.DecodeJSON(r.Body, &u)
		id, err := roleRepo.Create(u)
		fmt.Println(err)
		render.JSON(w, r, id)
	})

	permissionRepo := repository.NewDbPermissionsRepository(db)
	r.Get("/permissions/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		res, _ := permissionRepo.GetOne(id)
		render.JSON(w, r, res)
	})
	r.Get("/permissions", func(w http.ResponseWriter, r *http.Request) {
		res := permissionRepo.GetAll()
		render.JSON(w, r, res)
	})
	r.Post("/permissions", func(w http.ResponseWriter, r *http.Request) {
		var u model.Permission
		render.DecodeJSON(r.Body, &u)
		id, err := permissionRepo.Create(u)
		fmt.Println(err)
		render.JSON(w, r, id)
	})

	room := room.New()
	// http.Handle("/rooms/1", auth.TokenAuth(10, room))

	go room.Run() //get the room going in another thread
	//the chatting operation occur in the background
	//the main goroutine is running the web server

	log.Println("Starting web server on", config.Server.Port)

	err := http.ListenAndServe(config.Server.Port, r)
	if err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}

func loadConfig() config.Configuration {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	var config config.Configuration

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	return config
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
