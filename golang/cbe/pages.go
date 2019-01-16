package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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

func statsPage(w http.ResponseWriter, r *http.Request) {
	parsedTemplates, _ := template.ParseFiles("templates/stats.html")
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
