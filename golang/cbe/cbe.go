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
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
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

// PersonPayload ...
type PersonPayload struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// InteractionPayload ...
type InteractionPayload struct {
	Comment  string `json:"comment"`
	PersonID string `json:"personId"`
}

// Interaction ...
type Interaction struct {
	Person  string
	Date    string
	Comment string
	ID      string
}

// Person ...
type Person struct {
	Name string
	Type string
	ID   string
}

// PersonType ...
type PersonType struct {
	ID, Type string
}

// PersonInfo ...
type PersonInfo struct {
	ID, Name, Type, TypeName string
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
			auth(getPersonTypes, enterYourUserNamePassword),
		},

		Route{
			"getPersons",
			"GET",
			"/persons",
			auth(getPersons, enterYourUserNamePassword),
		},

		Route{
			"getFamilyInteractions",
			"GET",
			"/familyinteractions",
			auth(getFamilyInteractions, enterYourUserNamePassword),
		},

		Route{
			"getFriendInteractions",
			"GET",
			"/friendinteractions",
			auth(getFriendInteractions, enterYourUserNamePassword),
		},

		Route{
			"getCoworkersInteractions",
			"GET",
			"/coworkersinteractions",
			auth(getCoworkersInteractions, enterYourUserNamePassword),
		},

		Route{
			"getPersonsPerType",
			"GET",
			"/personspertype/{type}",
			auth(getPersonsPerType, enterYourUserNamePassword),
		},

		Route{
			"getPersonsInfo",
			"GET",
			"/personinfo/{id}",
			auth(getPersonInformation, enterYourUserNamePassword),
		},

		Route{
			"addPerson",
			"POST",
			"/addperson",
			auth(addPerson, enterYourUserNamePassword),
		},

		Route{
			"addInteraction",
			"POST",
			"/addinteraction",
			auth(addInteraction, enterYourUserNamePassword),
		},
	}
)

func getPersons(w http.ResponseWriter, r *http.Request) {

	types := []struct {
		ID   string
		Name string
	}{}

	rows, err := db.Query("SELECT id, name FROM person")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var id, name string
	for rows.Next() {
		rows.Scan(&id, &name)
		types = append(types, struct{ ID, Name string }{id, name})
	}

	json.NewEncoder(w).Encode(types)
}

func getPersonsPerType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	arg := vars["type"]

	personType, err := strconv.Atoi(arg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	persons, err := personsPerType(personType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(persons)
}

func getPersonInformation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	arg := vars["id"]

	personType, err := strconv.Atoi(arg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	info, err := personInfo(personType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(info)
}

func getFamilyInteractions(w http.ResponseWriter, r *http.Request) {
	interactions, err := getInteractionsPerType(1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(interactions)
}

func getCoworkersInteractions(w http.ResponseWriter, r *http.Request) {
	interactions, err := getInteractionsPerType(3)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(interactions)
}

func getFriendInteractions(w http.ResponseWriter, r *http.Request) {
	interactions, err := getInteractionsPerType(2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(interactions)
}

func getPersonTypes(w http.ResponseWriter, r *http.Request) {
	types, err := personTypes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(types)
}

func readForm(r *http.Request) *PersonPayload {
	r.ParseForm()
	person := new(PersonPayload)
	decoder := schema.NewDecoder()
	decodeErr := decoder.Decode(person, r.PostForm)
	if decodeErr != nil {
		log.Println("error mapping parsed form data to struct: ", decodeErr)
	}
	return person
}

func addPerson(w http.ResponseWriter, r *http.Request) {
	person := readForm(r)

	stmt, err := db.Prepare("INSERT INTO person (name, type) values(?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(person.Name, person.Type)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(person)
}

func addInteraction(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	interaction := new(InteractionPayload)
	decoder := schema.NewDecoder()
	decodeErr := decoder.Decode(interaction, r.PostForm)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusInternalServerError)
		return
	}

	stmt, err := db.Prepare("INSERT INTO interaction (comment, date, person_id) VALUES(?, date(now()), ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(interaction.Comment, interaction.PersonID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(interaction)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	parsedTemplates, _ := template.ParseFiles("templates/index.html")
	err := parsedTemplates.Execute(w, nil)
	if err != nil {
		log.Print("Error occurred while executing the template or writing its output: ", err)
		return
	}
}

func personasPage(w http.ResponseWriter, r *http.Request) {
	parsedTemplates, _ := template.ParseFiles("templates/personas.html")
	err := parsedTemplates.Execute(w, nil)
	if err != nil {
		log.Print("Error occurred while executing the template or writing its output: ", err)
		return
	}
}

func personInfoPage(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	personID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	info, err := personInfo(personID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	parsedTemplates, _ := template.ParseFiles("templates/persona.html")
	err = parsedTemplates.Execute(w, info)
	if err != nil {
		log.Print("Error occurred while executing the template or writing its output: ", err)
		return
	}
}

func interactionInfoPage(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	interID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	info, err := getInteractions(interID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	parsedTemplates, _ := template.ParseFiles("templates/interaction.html")
	err = parsedTemplates.Execute(w, info)
	if err != nil {
		log.Print("Error occurred while executing the template or writing its output: ", err)
		return
	}
}

func auth(handler http.HandlerFunc, realm string) http.HandlerFunc {
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

	router := mux.NewRouter().StrictSlash(false)
	router = addRoutes(router)

	defer db.Close()

	router.HandleFunc("/", auth(homePage, enterYourUserNamePassword))
	router.HandleFunc("/personas", auth(personasPage, enterYourUserNamePassword))
	router.HandleFunc("/person/{id}", auth(personInfoPage, enterYourUserNamePassword))
	router.HandleFunc("/interaction/{id}", auth(interactionInfoPage, enterYourUserNamePassword))

	router.PathPrefix("/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static/"))))

	err := http.ListenAndServe(connHost+":"+connPort, router)
	if err != nil {
		log.Fatal("error starting http server: ", err)
		return
	}

}
