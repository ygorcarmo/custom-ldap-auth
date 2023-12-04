package api

import (
	"custom/ldap-auth/utils"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

type response struct {
	Username string
	Token    string
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	fmt.Printf("the usernameis %s", username)

	fmt.Printf("username %s, password %s", username, password)

	err := utils.Authenticate(username, password)

	if err != nil {
		fmt.Println("Error")
		// http.Error(w, "Invalid username or password.", http.StatusBadRequest)
		htmlStr := fmt.Sprint("<div class='my-5'><span class='p-3 bg-red-400 rounded text-white'>Invalid username or password.</span></div>")
		tmpl, _ := template.New("error").Parse(htmlStr)
		tmpl.Execute(w, nil)
		return
	}
	claims := map[string]interface{}{"username": username, "hostname": "test123"}
	jwtauth.SetExpiry(claims, time.Now().Add(time.Hour*1))
	_, tokenString, _ := utils.TokenAuth.Encode(claims)

	res := response{Username: username, Token: tokenString}
	jsonRes, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &http.Cookie{Name: "jwt", Value: tokenString})
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}
