package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gaudiot.com/fonli/core"
	"gaudiot.com/fonli/src"
	"github.com/gorilla/mux"
)

func main() {
	envConfig, err := core.LoadEnvConfig()
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/native-to-foreign", handler).Methods("GET")

	log.Printf("Server is running on port :%s", envConfig.Port)
	log.Fatal(http.ListenAndServe(":"+envConfig.Port, router))
}

func handler(w http.ResponseWriter, r *http.Request) {
	exercises, err := src.CreateNativeToForeignExercise(10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(exercises)
}
