/**
leogtzr | leogutierrezramirez@gmail.com
*/
package main

import (
	"crypto/subtle"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const (
	connHost = "localhost"
	connPort = "8080"

	userEnvVar     = "CBE_USER"
	passwordEnvVar = "CBE_PASSWORD"

	enterYourUserNamePassword = "Please enter your username and password"
)

var (
	cbeUser     string
	cbePassword string
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

// Person ...
type Person struct {
	Name string
	Age  string
}

func homePage(w http.ResponseWriter, r *http.Request) {
	person := Person{Age: "1", Name: "Foo"}
	parsedTemplates, _ := template.ParseFiles("templates/index.html")
	err := parsedTemplates.Execute(w, person)
	if err != nil {
		log.Print("Error occurred while executing the template or writing its output: ", err)
		return
	}
}

func personasPage(w http.ResponseWriter, r *http.Request) {
	person := Person{Age: "1", Name: "Foo"}
	parsedTemplates, _ := template.ParseFiles("templates/personas.html")
	err := parsedTemplates.Execute(w, person)
	if err != nil {
		log.Print("Error occurred while executing the template or writing its output: ", err)
		return
	}
}

// BasicAuth ...
func BasicAuth(handler http.HandlerFunc, realm string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(cbeUser)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), []byte(cbePassword)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("You are Unauthorized to access the application.\n"))
			return
		}
		handler(w, r)
	}
}

func init() {
	if user, isSet := os.LookupEnv(userEnvVar); isSet {
		cbeUser = user
	} else {
		log.Fatalf("%s env variable not set.", userEnvVar)
	}
	if password, isSet := os.LookupEnv(passwordEnvVar); isSet {
		cbePassword = password
	} else {
		log.Fatalf("%s env variable not set.", passwordEnvVar)
	}
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/", BasicAuth(homePage, enterYourUserNamePassword))
	router.HandleFunc("/personas", BasicAuth(personasPage, enterYourUserNamePassword))

	router.PathPrefix("/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static/"))))

	err := http.ListenAndServe(connHost+":"+connPort, router)
	if err != nil {
		log.Fatal("error starting http server: ", err)
		return
	}

}
