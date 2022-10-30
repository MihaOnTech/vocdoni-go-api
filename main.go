package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type election struct {
	ID        int    `json:"ID"`
	Title     string `json:"Title"`
	StartDate string `json:"StartDate"`
	EndDate   string `json:"EndDate"`
	Question  string `json:"Question"`
	Choice1   string `json:"Choice1"`
	Choice2   string `json:"Choice2"`
	Choice3   string `json:"Choice3"`
	Status    string `json:"Status"`
}

type allElections []election

var elections = allElections{
	{
		ID:        1,
		Title:     "Test_Tribu election 1",
		StartDate: "23-01-2022 10:00",
		EndDate:   "25-01-2022 10:00",
		Question:  "Deberia Roy ser el Rey?",
		Choice1:   "Sí",
		Choice2:   "No",
		Choice3:   "",
		Status:    "OPEN",
	},
}

func getElections(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(elections)
}

func getElectionById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	electionId, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
	}

	for _, election := range elections {
		if election.ID == electionId {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(elections[electionId-1])
		}
	}

}

func deleteElection(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	electionId, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
	}

	for i, election := range elections {
		if election.ID == electionId {
			// Junta todo lo que hay ANTES de i y DESPUES de i
			elections = append(elections[:i], elections[i+1:]...)
			fmt.Fprintf(w, "Election with ID %v has benn removed succesfully", electionId)
		}
	}

}

func createElection(w http.ResponseWriter, r *http.Request) {
	var newElection election
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Invalid Data")
	}

	json.Unmarshal(reqBody, &newElection)

	newElection.ID = len(elections) + 1
	elections = append(elections, newElection)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newElection)
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to vocdoni API")
}

func main() {
	fmt.Println("Starting API...")
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/elections", getElections).Methods("GET")
	router.HandleFunc("/elections/{id}", getElectionById).Methods("GET")
	router.HandleFunc("/elections/{id}", deleteElection).Methods("DELETE")
	router.HandleFunc("/elections", createElection).Methods("POST")
	log.Fatal(http.ListenAndServe(":3000", router))

}
