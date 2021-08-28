package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
)

var DB *gorm.DB
var err error

const URL = "localadmin:admin@tcp(127.0.0.1:3306)/notes?charset=utf8mb4&parseTime=True&loc=Local"

// Note struct (Model)
type Note struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
}

func getNotes(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var notes []Note
	DB.Find(&notes)
	err = json.NewEncoder(w).Encode(notes)
	if err != nil {
		log.Println(err)
	}
}

func createNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var note Note
	_ = json.NewDecoder(r.Body).Decode(&note)
	DB.Create(&note)
	err = json.NewEncoder(w).Encode(note)
	if err != nil {
		log.Println(err)
	}
}

func updateNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var note Note
	DB.First(&note, params["id"])
	_ = json.NewDecoder(r.Body).Decode(&note)
	err = json.NewEncoder(w).Encode(note)
	if err != nil {
		log.Println(err)
	}
}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var note Note
	DB.Delete(&note, params["id"])
	err = json.NewEncoder(w).Encode("note deleted successfully")
	if err != nil {
		log.Println(err)
	}
}

func main() {
	DB, err = gorm.Open(mysql.Open(URL), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to DB")
	}
	err := DB.AutoMigrate(&Note{})
	if err != nil {
		log.Fatal(err)
	}
	router := mux.NewRouter()

	router.HandleFunc("/notes", getNotes).Methods("GET")
	router.HandleFunc("/notes", createNote).Methods("POST")
	router.HandleFunc("/notes/{id}", updateNote).Methods("PUT")
	router.HandleFunc("/notes/{id}", deleteNote).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
