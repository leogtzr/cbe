package main

import (
	"database/sql"
	"net/http"
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
	Name      string `json:"name"`
	Type      string `json:"type"`
	EveryDays string `json:"everydays"`
}

// InteractionPayload ...
type InteractionPayload struct {
	Comment  string `json:"comment"`
	PersonID string `json:"personId"`
	Date     string `json:"date"`
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
	ID, Name, Type, TypeName, EveryDays string
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
