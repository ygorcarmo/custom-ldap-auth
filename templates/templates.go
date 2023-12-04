package templates

import (
	"html/template"
	"log"
	"net/http"
)

func HandleLogin(w http.ResponseWriter, r *http.Request){
	templ, err := template.ParseFiles("templates/index.html")
	if err != nil{
		log.Fatal(err)
	}
	err = templ.Execute(w, nil)
}