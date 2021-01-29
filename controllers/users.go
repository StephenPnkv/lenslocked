package controllers

import (
	"../views"
	"fmt"
	"log"
	"net/http"
)

type Users struct {
	NewView *views.View
}

type SignupForm struct{
	Email string `schema:"email"`
	Password string `schema: "password"`
}

func NewUser() *Users {
	return &Users{
		NewView: views.NewView("bootstrap",
			"users/signup"),
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
	fmt.Fprintln(w, form)

}
