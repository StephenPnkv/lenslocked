package controllers

import (
	"lenslocked/views"
	"fmt"
	"log"
	"net/http"
	"lenslocked/models"
	"lenslocked/rand"
)

type Users struct {
	NewView *views.View
	LoginView *views.View
	us models.UserService
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

// @route: GET /signup
// @desc: Renders sign in view
// @access: Public

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		log.Panicln(err)
	}
}

// @route: POST /signup
// @desc: User is able to fill out the form and create an account
// @access: Public
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

	err := u.signIn(w, &user)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// @route: POST /login
// @desc: User is logged in and redirected back to the home page
// @access: Public
func (u *Users) Login(w http.ResponseWriter, r *http.Request){
	//Parse the form
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil{
		log.Panicln(err)
	}

	//Authenticate user
	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil{
		switch err {
		case models.ErrNotFound:
			fmt.Fprintln(w, "Invalid email provided!")
		case models.ErrInvalidPassword:
			fmt.Fprintln(w, "Invalid password provided!")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	//Sign in user
	err = u.signIn(w, user)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (u *Users) signIn(w http.ResponseWriter, user *models.User) error{
	if user.Remember == ""{
		token, err := rand.RememberToken()
		if err != nil{
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil{
			return err
		}
	}
		cookie := http.Cookie{
			Name: "remember_token",
			Value: user.Remember,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		return nil

}

func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie("remember_token")
	if err != nil{
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}

	user, err := u.us.ByRemember(cookie.Value)
	if err != nil{
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, user)
}
