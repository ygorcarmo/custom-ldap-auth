package api

import (
	"custom/ldap-auth/utils"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
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
		http.Error(w, "Invalid username or password.", http.StatusBadRequest)
		htmlStr := fmt.Sprint("<span class='p-2 border border-red-400 rounded'>Invalid username or password.</span>")
		tmpl, _ := template.New("error").Parse(htmlStr)
		tmpl.Execute(w, nil)
		return
	}

	_, tokenString, _ := utils.TokenAuth.Encode(map[string]interface{}{"username": username, "hostname": "test123"})

	res := response{Username: username, Token: tokenString}
	jsonRes, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &http.Cookie{Name: "jwt", Value: tokenString})
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}
