package main

import (
	"fmt"
	"log"

	"gopkg.in/ldap.v2"
)

func main() {

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "ldap.forumsys.com", 389))
	if err != nil {
		log.Fatal(err)
	}
	// Reconnect with TLS
	// err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	defer l.Close()

	// First bind with a read only user
	err = l.Bind("cn=read-only-admin,dc=example,dc=com", "password")
	if err != nil {
		log.Fatal(err)
	}

	// Search for the given username

	user := "newton"
	// password := "password"
	baseDN := "dc=example,dc=com"

	searchReq := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", user),
		[]string{"dn"},
		nil,
	)

	result, err := l.Search(searchReq)
	if err != nil {
		log.Fatal("failed to query LDAP: %w", err)
	}

	log.Println("Got", len(result.Entries), "search results")

	if len(result.Entries) != 1 {
		log.Fatal("User does not exist or too many entries returned")
	}
	// Search for the given username
	log.Println(result.Entries[0])

	userdn := result.Entries[0].DN

	// Bind as the user to verify their password
	err = l.Bind(userdn, "password")
	if err != nil {
		log.Fatal(err)
	}

}
