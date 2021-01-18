package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

var homeTemplate, aboutTemplate, contactTemplate *template.Template

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := homeTemplate.Execute(w, nil); err != nil {
		fmt.Fprintf(w, "<h1>Error</h1>")
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := contactTemplate.Execute(w, nil); err != nil {
		log.Fatal(err)
	}
}

func about(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := aboutTemplate.Execute(w, nil); err != nil {
		log.Fatal(err)
	}
}

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>FAQ</h1>")
}

func handleError(err error) {
	log.Panicln(err)
}

func main() {

	var err error
	homeTemplate, err = template.ParseFiles("views/home.gohtml")
	if err != nil {
		log.Fatal(err)
	}
	contactTemplate, err = template.ParseFiles("views/contact.gohtml")
	if err != nil {
		log.Fatal(err)
	}
	aboutTemplate, err = template.ParseFiles("views/aboutTemplate.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/about", about)
	r.HandleFunc("/faq", faq)
	//http.HandleFunc("/", handlerFunc)

	http.ListenAndServe(":3000", r)
}
