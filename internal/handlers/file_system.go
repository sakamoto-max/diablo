package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sakamoto-max/diablo/internal/dto"
	"github.com/sakamoto-max/diablo/internal/services"
)

type FileSystem struct {
	service *services.Synchronizer
}

func (f *FileSystem) CreateNewFileSystem(w http.ResponseWriter, r *http.Request) {

	log.Println("req received in handler")

	var input dto.NewSuiteReq

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("error occured %w", err)
		return
	}

	err = f.service.New(r.Context(), &input)
	if err != nil {

		// todo
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (f *FileSystem) Sync(w http.ResponseWriter, r *http.Request) {

	var input dto.EventsReq
	json.NewDecoder(r.Body).Decode(&input)

	err := f.service.Sync(r.Context(), &input)
	if err != nil {
		fmt.Println(err)
		// todo
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (f *FileSystem) Suite(w http.ResponseWriter, r *http.Request) {

	var input dto.Event
	json.NewDecoder(r.Body).Decode(&input)

	suite, err := f.service.GetSuite(r.Context(), input)
	if err != nil {
		fmt.Println(err)
		// todo
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(suite)
}

func (f *FileSystem) Ping(w http.ResponseWriter, r *http.Request) {

	var input dto.UserIp

	json.NewDecoder(r.Body).Decode(&input)

	allEvents, err := f.service.Ping(r.Context(), input)
	if err != nil {
		fmt.Println(err)
		// todo
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allEvents)
}
