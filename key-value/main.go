package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var logger TransactionLogger

func initTransactionLogger() error {
	var err error

	logger, err = NewFileTransactionLogger("transaction.log")
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errors := logger.ReadEvents()
	e, ok := EventLog{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors:
			log.Writer()
		case e, ok = <-events:
			switch e.Type {
			case EventDelete:
				err = Delete(e.Key)
			case EventPut:
				err = Put(e.Key, e.Value)
			}
		}
	}

	logger.Run()

	return err
}

func keyPutHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)

	_, errGet := Get(vars["key"])

	err = Put(vars["key"], string(b))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.WritePut(vars["key"], string(b))

	if errors.Is(errGet, ErrorNoSuchKey) {
		w.WriteHeader(http.StatusCreated)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func keyGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	value, err := Get(vars["key"])

	if err != nil {
		if errors.Is(err, ErrorNoSuchKey) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	w.Write([]byte(value))
}

func keyDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := Delete(vars["key"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.WriteDelete(vars["key"])

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	initTransactionLogger()

	r := mux.NewRouter()

	r.HandleFunc("/v1/key/{key}", keyPutHandler).Methods("PUT")
	r.HandleFunc("/v1/key/{key}", keyGetHandler).Methods("GET")
	r.HandleFunc("/v1/key/{key}", keyDeleteHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))

	logger.Close()
}
