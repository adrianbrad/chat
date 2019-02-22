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

var rooms map[string]room.Room

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

	room := room.New()
	http.Handle("/rooms/1", auth.TokenAuth(10, room))

	go room.Run() //get the room going in another thread
	//the chatting operation occur in the background
	//the main goroutine is running the web server

	log.Println("Starting web server on", *addr)

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}

var counter = 1

//*TODO: find a way to create dynamic links for rooms, or handle the requests based on the url
//format can be /rooms or /rooms/{id}
func roomHandler(w http.ResponseWriter, r *http.Request) {
	//if path is /rooms we expect a post request to create a room
	// if path.Base(r.URL.Path) == "rooms" {
	// 	if r.Method != http.MethodPost {
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		return
	// 	}
	// }
	room := room.New()
	http.Handle("/rooms/"+string(counter), auth.TokenAuth(10, room))
	go room.Run() //get the room going in another thread
	//the chatting operation occur in the background
	//the main goroutine is running the web server

}
