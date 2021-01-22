

package controllers
import (
  "../views"
)

func NewStatic() *Static{
  return &Static{
    Home: views.NewView("bootstrap", "views/home.gohtml"),
  	Contact: views.NewView("bootstrap", "views/contact.gohtml"),
    About : views.NewView("bootstrap", "views/about.gohtml"),
  }
}

type Static struct{
  Home *views.View
  Contact *views.View
  About *views.View
}
