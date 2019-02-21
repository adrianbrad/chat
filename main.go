package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/adrianbrad/chat/auth"
	"github.com/adrianbrad/chat/room"
)

var users = map[string]*auth.User{
	"1": &auth.User{Name: "brad", Role: true},
	"2": &auth.User{Name: "john", Role: false},
	"3": &auth.User{Name: "eusebiu", Role: true},
}

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
	addr := flag.String("addr", ":8080", "The address of the application.")
	flag.Parse()

	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("/home/brad/workspace/go/src/github.com/adrianbrad/chat/assets"))))

	http.Handle("/chat", &templateHandler{filename: "chat.html"})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	r := room.New()

	http.Handle("/room", auth.TokenAuth(10, r))

	go r.Run() //get the room going in another thread
	//the chatting operation occur in the background
	//the main goroutine is running the web server

	log.Println("Starting web server on", *addr)

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}
