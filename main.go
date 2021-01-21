package main

import (
	"./views"
	"fmt"
	"github.com/gorilla/mux"
	//"html/template"
	"log"
	"net/http"
)

var (
	homeView    *views.View
	aboutView   *views.View
	contactView *views.View
	signupView  *views.View
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	handleError(homeView.Render(w, nil))
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	handleError(contactView.Render(w, nil))

}

func about(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	handleError(aboutView.Render(w, nil))
}

func signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	handleError(signupView.Render(w, nil))
}

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>FAQ</h1>")
}

func handleError(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func main() {

	//var err error
	homeView = views.NewView("bootstrap", "views/home.gohtml")
	contactView = views.NewView("bootstrap", "views/contact.gohtml")
	aboutView = views.NewView("bootstrap", "views/about.gohtml")
	signupView = views.NewView("bootstrap", "views/signup.gohtml")

	r := mux.NewRouter()

	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/about", about)
	r.HandleFunc("/signup", signup)
	r.HandleFunc("/faq", faq)
	//http.HandleFunc("/", handlerFunc)

	http.ListenAndServe(":3000", r)
}
