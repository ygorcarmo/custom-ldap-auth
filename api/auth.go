package api

import (
	"custom/ldap-auth/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type response struct{
	Username string
	Token string
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	fmt.Printf("username %s, password %s", username, password)

	err := utils.Authenticate(username, password);

	if err != nil{
		http.Error(w, "Invalid username or password.", http.StatusBadRequest)
		return
	}

	_, tokenString, _ := utils.TokenAuth.Encode(map[string]interface{}{"username": username})

	res := response{Username: username, Token: tokenString}
	jsonRes, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}