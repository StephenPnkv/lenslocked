
package main

import (
	"net/http"
	"fmt"
)

func handlerFunc(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html")
	if r.URL.Path == "/"{
		fmt.Fprint(w, "<h1>Welcome to LensLocked!</h1>")
	}else if r.URL.Path == "/contact"{
		fmt.Fprint(w, "<h1>Contact</h1>")
	}else if r.URL.Path == "/about"{
		fmt.Fprint(w, "<h1>About</h1>")
	}
}

func main(){
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":3000", nil)
}