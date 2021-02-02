package controllers

import (
	"../views"
	"fmt"
	"log"
	"net/http"
	"../models"
)

type Users struct {
	NewView *views.View
	LoginView *views.View
	us *models.UserService
}

type SignupForm struct{
	Name string `schema:"name"`
	Email string `schema:"email"`
	Password string `schema: "password"`
}

type LoginForm struct{
	Email string `schema:"email"`
	Password string `schema:"password"`
}

func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView: views.NewView("bootstrap","users/signup"),
		LoginView: views.NewView("bootstrap","users/login"),
		us: us,
	}
}

// Get /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		log.Panicln(err)
	}
}

// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {

	form := SignupForm{}
	if err := parseForm(r, &form); err != nil{
		log.Panicln(err)
	}
	user := models.User{
		Name: form.Name,
		Email: form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, user)
}
//POST /login
//User is logged in and redirected back to the home page
func (u *Users) Login(w http.ResponseWriter, r *http.Request){
	//Parse the form
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil{
		log.Panicln(err)
	}

	//Authenticate user
	user, err := u.us.Authenticate(form.Email, form.Password)
	switch err {
	case models.ErrNotFound:
		fmt.Fprintln(w, "Invalid email provided!")
	case models.ErrInvalidPassword:
		fmt.Fprintln(w, "Invalid password provided!")
	case nil:
		fmt.Fprintln(w, user)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
