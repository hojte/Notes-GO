package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var nt = &Note{
	Title:       "Test Title",
	Description: "Test Description",
}

func TestCreateNote(t *testing.T) {
	setUp()

	var jsonReq = []byte(`{"title": "` + nt.Title + `", "description": "` + nt.Description + `"}`)
	req, _ := http.NewRequest("POST", "/notes", bytes.NewBuffer(jsonReq))

	handler := http.HandlerFunc(createNote)

	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)

	status := response.Code

	if status != http.StatusOK {
		t.Errorf("Handler returned a wrong status code: got %v want %v", status, http.StatusOK)
	}

	var note Note
	_ = json.NewDecoder(io.Reader(response.Body)).Decode(&note)

	assert.NotNil(t, note.ID)
	assert.Equal(t, nt.Title, note.Title)
	assert.Equal(t, nt.Description, note.Description)

	tearDown()
}

func TestGetPosts(t *testing.T) {
	setUp()

	req, _ := http.NewRequest("GET", "/posts", nil)

	handler := http.HandlerFunc(getNotes)

	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)

	status := response.Code

	if status != http.StatusOK {
		t.Errorf("Handler returned a wrong status code: got %v want %v", status, http.StatusOK)
	}

	var notes []Note
	_ = json.NewDecoder(io.Reader(response.Body)).Decode(&notes)

	assert.NotNil(t, notes[0].ID)
	assert.Equal(t, nt.Title, notes[0].Title)
	assert.Equal(t, nt.Description, notes[0].Description)

	tearDown()
}

func TestUpdateNote(t *testing.T) {
	setUp()

	var jsonReq = []byte(`{"title": "Updated Title", "description": "Updated Description"}`)
	req, _ := http.NewRequest("PUT", "/notes/1", bytes.NewBuffer(jsonReq))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	handler := http.HandlerFunc(updateNote)

	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)

	status := response.Code

	if status != http.StatusOK {
		t.Errorf("Handler returned a wrong status code: got %v expected %v", status, http.StatusOK)
	}

	var note Note
	_ = json.NewDecoder(io.Reader(response.Body)).Decode(&note)

	assert.NotNil(t, note.ID)
	assert.Equal(t, "Updated Title", note.Title)
	assert.Equal(t, "Updated Description", note.Description)

	tearDown()
}

func TestUpdateNonExistingNote(t *testing.T) {
	setUp()

	var jsonReq = []byte(`{"title": "Updated Title", "description": "Updated Description"}`)
	req, _ := http.NewRequest("PUT", "/notes/4", bytes.NewBuffer(jsonReq))
	req = mux.SetURLVars(req, map[string]string{"id": "4"})

	handler := http.HandlerFunc(updateNote)

	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)

	status := response.Code

	if status != http.StatusBadRequest {
		t.Errorf("Handler returned a wrong status code: got %v expected %v", status, http.StatusOK)
	}

	tearDown()
}

func TestDeleteNote(t *testing.T) {
	setUp()
	//Create a new HTTP POST request
	req, _ := http.NewRequest("DELETE", "/notes/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	handler := http.HandlerFunc(deleteNote)

	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)

	status := response.Code

	if status != http.StatusOK {
		t.Errorf("Handler returned a wrong status code: got %v want %v", status, http.StatusOK)
	}

	tearDown()
}

func setUp() {
	InitDB("root:root@tcp(127.0.0.1:3306)/notes_test?charset=utf8mb4&parseTime=True&loc=Local")
	DB.Create(nt)
}

func tearDown() {
	DB.Exec("TRUNCATE TABLE notes;")
}
