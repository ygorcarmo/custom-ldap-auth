package utils

import (
	"fmt"
	"log"

	"github.com/go-chi/jwtauth/v5"
	"gopkg.in/ldap.v2"
)
 
var TokenAuth *jwtauth.JWTAuth

func init() {
	// change secret as it is easily guessable and remove it from source code
	// fmt.Println(os.Getenv("big big secret :)"))
	TokenAuth = jwtauth.New("HS256", []byte("secret"), nil )

	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:123` here:
	_, tokenString, _ := TokenAuth.Encode(map[string]interface{}{"user_id": 123})
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
}

func Authenticate(username string, password string) error{
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "ldap.forumsys.com", 389))
	if err != nil {
		log.Fatal(err)
	}
	// Reconnect with TLS
	// err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// close connection when done
	defer l.Close()

	// First bind with a read only user
	err = l.Bind("cn=read-only-admin,dc=example,dc=com", "password")
	if err != nil {
		log.Fatal(err)
	}

	// Search for the given username

	baseDN := "dc=example,dc=com"

	searchReq := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", username),
		[]string{"dn"},
		nil,
	)

	result, err := l.Search(searchReq)
	if err != nil {
		fmt.Println("failed to query LDAP: %w", err)
		return err
	}

	log.Println("Got", len(result.Entries), "search results")

	if len(result.Entries) != 1 {
		fmt.Println("User does not exist or too many entries returned")
		return err
	}
	// Search for the given username
	log.Println(result.Entries[0])

	userdn := result.Entries[0].DN

	// Bind as the user to verify their password
	err = l.Bind(userdn, password)
	if err != nil {
		fmt.Println(err)
		return err

	}
	log.Println("Logged in")

	return nil
}