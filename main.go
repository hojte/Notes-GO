package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

var DB *gorm.DB
var err error

// Note struct (Model)
type Note struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
}

func getNotes(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var notes []Note
	err = DB.Find(&notes).Error
	if err != nil {
		log.Println("Error", err)
		_ = json.NewEncoder(w).Encode(err.Error())
	}
	_ = json.NewEncoder(w).Encode(notes)

}

func createNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var note Note
	_ = json.NewDecoder(r.Body).Decode(&note)
	err = DB.Create(&note).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
	} else {
		_ = json.NewEncoder(w).Encode(note)
	}

}

func updateNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	log.Println(mux.Vars(r))
	var note Note
	if s, err := strconv.Atoi(id); err == nil {
		note.ID = uint(s)
	} else {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}
	err = DB.First(&note).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	var updatedNote Note
	_ = json.NewDecoder(r.Body).Decode(&updatedNote)

	note.Title = updatedNote.Title
	note.Description = updatedNote.Description

	err = DB.Save(&note).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
	} else {
		_ = json.NewEncoder(w).Encode(note)
	}

}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	log.Println("id:", params["id"])
	err := DB.Delete(&Note{}, params["id"]).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
	}
}

func InitDB(URL string) {
	DB, err = gorm.Open(mysql.Open(URL), &gorm.Config{})
	if err != nil {
		log.Println("Error", err)
		panic("Cannot connect to DB")
	}
	err := DB.AutoMigrate(&Note{})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	InitDB("root:root@tcp(db:3306)/notes?charset=utf8mb4&parseTime=True&loc=Local")
	log.Println("Serving restapi on port 8000...")

	router := mux.NewRouter()

	router.HandleFunc("/notes", getNotes).Methods("GET")
	router.HandleFunc("/notes", createNote).Methods("POST")
	router.HandleFunc("/notes/{id}", updateNote).Methods("PUT")
	router.HandleFunc("/notes/{id}", deleteNote).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
