package main

import (
	"./controllers"
	"./views"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var (
	homeView    *views.View
	aboutView   *views.View
	contactView *views.View
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
	staticCon := controllers.NewStatic()
	usersCon := controllers.NewUser()
	r := mux.NewRouter()

	r.Handle("/", staticCon.Home).Methods("GET")
	r.Handle("/contact", staticCon.Contact).Methods("GET")
	r.Handle("/about", staticCon.About).Methods("GET")
	r.HandleFunc("/signup", usersCon.New).Methods("GET")
	r.HandleFunc("/signup", usersCon.Create).Methods("POST")
	//r.HandleFunc("/faq", faq).Methods("GET")
	//http.HandleFunc("/", handlerFunc)

	http.ListenAndServe(":3000", r)
}
