package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello gorilla/mux!\n"))
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	logVars(r)
}

func articlesCategoryHandler(w http.ResponseWriter, r *http.Request) {
	logVars(r)
}

func articleHandler(w http.ResponseWriter, r *http.Request) {
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

	r.HandleFunc("/products/{key}", productsHandler)
	r.HandleFunc("/articles/{category}/", articlesCategoryHandler)
	r.HandleFunc("/articles/{category}/{id:[0-9]+}", articleHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}
