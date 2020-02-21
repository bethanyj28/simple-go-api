package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// community stores the data for this demo so that we can represent data transfer without needing to dive into databases
var community *People

func main() {
	community = &People{} // initialize an array of people

	http.HandleFunc("/hello", sayHelloHandler)
	http.HandleFunc("/add", addPeopleHandler)
	http.HandleFunc("/person/", personHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Person represents a person we can say hello to
type Person struct {
	Nickname   string `json:"nickname"` // represents a unique id
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name,omitempty"`
	LastName   string `json:"last_name"`
	Age        int32  `json:"age"`
}

// People is a group of persons
type People struct {
	People []*Person `json:"people"`
}

// Handlers
func addPeopleHandler(w http.ResponseWriter, r *http.Request) {
	var p Person
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		log.Fatal("failed to decode person")
	}

	addPerson(&p)
}

func personHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		nickname := r.URL.Path[len("/person/"):] // grabs everything after the person route
		p := getPerson(nickname)
		if p == nil { // no person found
			http.Error(w, "no person found", http.StatusNoContent)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	case http.MethodPost:
		var p Person
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, "failed to decode person", http.StatusBadRequest)
			return
		}
		addPerson(&p)
	case http.MethodDelete:
		nickname := r.URL.Path[len("/person/"):] // grabs everything after the person route
		p := deletePerson(nickname)
		if p == nil { // no person found
			http.Error(w, "no person found", http.StatusNoContent)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	case http.MethodPut:
		var p Person
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, "failed to decode person", http.StatusBadRequest)
			return
		}
		updatePerson(&p)
	default:
		http.Error(w, "method not supported", http.StatusInternalServerError)
	}
}

func sayHelloHandler(w http.ResponseWriter, r *http.Request) {
	for _, p := range community.People {
		fmt.Fprintf(w, "Hello, %s %s!", p.FirstName, p.LastName)
	}
}

// Service methods

func addPerson(p *Person) {
	community.People = append(community.People, p)
}

func getPerson(nickname string) *Person {
	for _, p := range community.People {
		if p.Nickname == nickname {
			return p
		}
	}

	return nil
}

func deletePerson(nickname string) *Person {
	for i, p := range community.People {
		if p.Nickname == nickname {
			community.People = append(community.People[:i], community.People[i+1:]...)
			return p
		}
	}

	return nil
}

func updatePerson(p *Person) {
	for i, op := range community.People {
		if op.Nickname == p.Nickname {
			community.People[i] = p
			return
		}
	}
}
