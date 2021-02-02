package main

import (
	"./controllers"
	"./models"
	"./views"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
  host = "localhost"
  port = 5432
  user = "postgres"
  dbname= "lenslocked_dev"
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
	//Create a connection db connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
	//Create user service
	userService, err := models.NewUserService(psqlInfo)
  if err != nil{
    log.Panicln(err)
  }
  defer userService.Close()
  userService.DestructiveReset()

	//Controllers
	staticController := controllers.NewStatic()
	userController := controllers.NewUsers(userService)

	//Routing
	r := mux.NewRouter()

	r.Handle("/", staticController.Home).Methods("GET")
	r.Handle("/contact", staticController.Contact).Methods("GET")
	r.Handle("/about", staticController.About).Methods("GET")

	r.HandleFunc("/signup", userController.New).Methods("GET")
	r.HandleFunc("/signup", userController.Create).Methods("POST")

	r.Handle("/login", userController.LoginView).Methods("GET")
	r.HandleFunc("/login", userController.Login).Methods("POST")

	http.ListenAndServe(":3000", r)
}
