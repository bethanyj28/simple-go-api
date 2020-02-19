package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var community *People

func main() {
	community = &People{} // initialize an array of people

	http.HandleFunc("/hello", sayHelloHandler)
	http.HandleFunc("/add", addPeopleHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Person represents a person we can say hello to
type Person struct {
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name,omitempty"`
	LastName   string `json:"last_name"`
	Age        int32  `json:"age"`
}

// People is a group of persons
type People struct {
	People []*Person `json:"people"`
}

func addPeopleHandler(w http.ResponseWriter, r *http.Request) {
	var p Person
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		log.Fatal("failed to decode person")
	}

	community.People = append(community.People, &p)
}

func sayHelloHandler(w http.ResponseWriter, r *http.Request) {
	for _, p := range community.People {
		fmt.Fprintf(w, "Hello, %s %s!", p.FirstName, p.LastName)
	}
}
