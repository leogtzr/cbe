/**
leogtzr | leogutierrezramirez@gmail.com
*/
package main

import (
	"crypto/subtle"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	connHost = "localhost"
	connPort = "8080"

	userEnvVar       = "CBE_USER"
	passwordEnvVar   = "CBE_PASSWORD"
	userDBEnvVar     = "DB_USER"
	passwordDBEnvVar = "DB_PASSWORD"
	dbNameEnvVar     = "DB_NAME"
	driverName       = "mysql"

	enterYourUserNamePassword = "Please enter your username and password"
)

// Route ...
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes ...
type Routes []Route

var (
	cbeUser     string
	cbePassword string
	// DB variables:
	dbUser     string
	dbPassword string
	dbName     string

	db              *sql.DB
	connectionError error

	routes = Routes{
		Route{
			"getPersonTypes",
			"GET",
			"/persontypes",
			BasicAuth(personTypes, enterYourUserNamePassword),
		},
	}
)

func personTypes(w http.ResponseWriter, r *http.Request) {

	types := []string{}

	rows, err := db.Query("SELECT type FROM person_type")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var tp string
	for rows.Next() {
		rows.Scan(&tp)
		types = append(types, tp)
	}

	json.NewEncoder(w).Encode(types)
}

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
	// check if the necessary env variables are set:
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

	if user, isSet := os.LookupEnv(userDBEnvVar); isSet {
		dbUser = user
	} else {
		log.Fatalf("%s env variable not set.", userDBEnvVar)
	}
	if password, isSet := os.LookupEnv(passwordDBEnvVar); isSet {
		dbPassword = password
	} else {
		log.Fatalf("%s env variable not set.", passwordDBEnvVar)
	}

	if password, isSet := os.LookupEnv(passwordDBEnvVar); isSet {
		cbePassword = password
	} else {
		log.Fatalf("%s env variable not set.", passwordDBEnvVar)
	}

	if db, isSet := os.LookupEnv(dbNameEnvVar); isSet {
		dbName = db
	} else {
		log.Fatalf("%s env variable not set.", dbNameEnvVar)
	}

	db, connectionError = sql.Open(driverName, fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName))
	if connectionError != nil {
		log.Fatal("error connecting to database :: ", connectionError)
	}
}

func addRoutes(router *mux.Router) *mux.Router {
	for _, route := range routes {
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)

	}

	return router
}

func main() {

	router := mux.NewRouter()
	router = addRoutes(router)

	defer db.Close()

	router.HandleFunc("/", BasicAuth(homePage, enterYourUserNamePassword))
	router.HandleFunc("/personas", BasicAuth(personasPage, enterYourUserNamePassword))

	router.PathPrefix("/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static/"))))

	err := http.ListenAndServe(connHost+":"+connPort, router)
	if err != nil {
		log.Fatal("error starting http server: ", err)
		return
	}

}
