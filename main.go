package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/stretchr/objx"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"

	"github.com/adrianbrad/chat/trace"
)

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

	authCookie, err := r.Cookie("auth")
	if err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	_ = t.template.Execute(w, data)
}

func main() {
	addr := flag.String("addr", ":8080", "The address of the application.")
	flag.Parse()

	//setup gomniauth
	gomniauth.SetSecurityKey("test")
	gomniauth.WithProviders(
		google.New("875697936126-mtpa207uap62gcf7d534qv7dv4b3o6li.apps.googleusercontent.com", "G1Il9Yv-QFPwg-XEnm-xOVKD",
			"http://localhost:8080/auth/callback/google"),
	)

	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("/home/brad/workspace/go/src/github.com/adrianbrad/chat/assets"))))

	// http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/room", r)

	go r.run() //get the room going in another thread
	//the chatting operation occur in the background
	//the main goroutine is running the web server

	log.Println("Starting web server on", *addr)

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}
