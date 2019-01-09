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
}

// Person ...
type Person struct {
	Name string
	Age  string
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
			BasicAuth(getPersonTypes, enterYourUserNamePassword),
		},

		Route{
			"getPersons",
			"GET",
			"/persons",
			BasicAuth(getPersons, enterYourUserNamePassword),
		},

		Route{
			"getFamilyInteractions",
			"GET",
			"/familyinteractions",
			BasicAuth(getFamilyInteractions, enterYourUserNamePassword),
		},

		Route{
			"getFriendInteractions",
			"GET",
			"/friendinteractions",
			BasicAuth(getFriendInteractions, enterYourUserNamePassword),
		},

		Route{
			"getCoworkersInteractions",
			"GET",
			"/coworkersinteractions",
			BasicAuth(getCoworkersInteractions, enterYourUserNamePassword),
		},

		Route{
			"addPerson",
			"POST",
			"/addperson",
			BasicAuth(addPerson, enterYourUserNamePassword),
		},

		Route{
			"addInteraction",
			"POST",
			"/addinteraction",
			BasicAuth(addInteraction, enterYourUserNamePassword),
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

func getInteractionsPerType(personType int) ([]Interaction, error) {
	interactions := []Interaction{}

	stmt, err := db.Query(`select concat(p.name, ' (', pt.type, ')') person,inter.date, inter.comment
	FROM interaction inter
	inner join person p on p.id = inter.person_id
	inner join person_type pt on pt.id = p.type
	where p.type = ?`, personType)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var person, date, comment string
	for stmt.Next() {
		stmt.Scan(&person, &date, &comment)
		interactions = append(interactions, Interaction{Person: person, Date: date, Comment: comment})
	}

	return interactions, nil
}

func getPersonsPerType(personType int) ([]string, error) {
	persons := []string{}

	stmt, err := db.Query(`select distinct(p.name)
	FROM person p
	inner join person_type pt on pt.id = p.type
	where p.type = ?`, personType)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var name string
	for stmt.Next() {
		stmt.Scan(&name)
		persons = append(persons, name)
	}

	return persons, nil
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
	types := []struct {
		ID   string
		Type string
	}{}

	rows, err := db.Query("SELECT id, type FROM person_type")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var tp, id string
	for rows.Next() {
		rows.Scan(&tp, &id)
		types = append(types, struct{ ID, Type string }{tp, id})
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
