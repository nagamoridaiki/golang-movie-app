package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
)

// Temps is template structure.
type Temps struct {
	notemp *template.Template
	indx   *template.Template
	helo   *template.Template
}

// Template for no-template.
func notemp() *template.Template {
	src := "<html><body><h1>NO TEMPLATE.</h1></body></html>"
	tmp, _ := template.New("index").Parse(src)
	return tmp
}

// setup template function.
func setupTemp() *Temps {
	temps := new(Temps)

	temps.notemp = notemp()

	// set index template.
	indx, er := template.ParseFiles("templates/index.html")
	if er != nil {
		indx = temps.notemp
	}
	temps.indx = indx

	// set hello template.
	helo, er := template.ParseFiles("templates/hello.html")
	if er != nil {
		helo = temps.notemp
	}
	temps.helo = helo

	return temps
}

// index handler.
func index(w http.ResponseWriter, rq *http.Request, tmp *template.Template) {
	er := tmp.Execute(w, nil)
	if er != nil {
		log.Fatal(er)
	}
}

var cs *sessions.CookieStore = sessions.NewCookieStore([]byte("secret-key-12345"))

// hello handler.
func hello(w http.ResponseWriter, rq *http.Request,
	tmp *template.Template) {
	msg := "login name & password:"

	ses, _ := cs.Get(rq, "hello-session")

	if rq.Method == "POST" {
		ses.Values["login"] = nil
		ses.Values["name"] = nil
		nm := rq.PostFormValue("name")
		pw := rq.PostFormValue("pass")
		if nm == pw {
			ses.Values["login"] = true
			ses.Values["name"] = nm
		}
		ses.Save(rq, w)
	}

	flg, _ := ses.Values["login"].(bool)
	lname, _ := ses.Values["name"].(string)
	if flg {
		msg = "logined: " + lname
	}

	item := struct {
		Title   string
		Message string
	}{
		Title:   "Session",
		Message: msg,
	}

	er := tmp.Execute(w, item)
	if er != nil {
		log.Fatal(er)
	}
}

// main program.
func main() {
	temps := setupTemp()
	// index handling.
	http.HandleFunc("/", func(w http.ResponseWriter, rq *http.Request) {
		index(w, rq, temps.indx)
	})
	// hello handling
	http.HandleFunc("/hello", func(w http.ResponseWriter, rq *http.Request) {
		hello(w, rq, temps.helo)
	})

	http.ListenAndServe("", nil)
}
