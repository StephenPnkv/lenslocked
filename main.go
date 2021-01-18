package main

import (
	"./views"
	"fmt"
	"github.com/gorilla/mux"
	//"html/template"
	"log"
	"net/http"
)

var homeView, aboutView, contactView *views.View

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := homeView.Template.ExecuteTemplate(w, homeView.Layout, nil); err != nil {
		fmt.Fprintf(w, "<h1>Error</h1>")
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := contactView.Template.ExecuteTemplate(w, contactView.Layout, nil); err != nil {
		log.Fatal(err)
	}
}

func about(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := aboutView.Template.ExecuteTemplate(w, aboutView.Layout, nil); err != nil {
		log.Fatal(err)
	}
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

	r := mux.NewRouter()

	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/about", about)
	r.HandleFunc("/faq", faq)
	//http.HandleFunc("/", handlerFunc)

	http.ListenAndServe(":3000", r)
}
