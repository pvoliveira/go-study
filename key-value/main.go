package main

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello gorilla/mux!\n"))
}

func keyPutHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)

	_, errGet := Get(vars["key"])

	err = Put(vars["key"], string(b))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if errors.Is(errGet, ErrorNoSuchKey) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func keyGetHandler(w http.ResponseWriter, r *http.Request) {
	logVars(r)
}

func keyDeleteHandler(w http.ResponseWriter, r *http.Request) {
	logVars(r)
}

func logVars(r *http.Request) {
	vars := mux.Vars(r)
	for k, v := range vars {
		log.Printf("%v: %v", k, v)
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", handler)

	r.HandleFunc("/v1/key/{key}", keyPutHandler).Methods("PUT")
	r.HandleFunc("/v1/key/{key}", keyGetHandler).Methods("GET")
	r.HandleFunc("/v1/key/{key}", keyDeleteHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
