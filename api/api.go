package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	db "../db"
	"github.com/gorilla/mux"
)

type JsonKeyPair struct {
	CodeName   string `json:"code_name,omitempty"`
	PublicKey  string `json:"public_key,omitempty"`
	PrivateKey string `json:"private_key,omitempty"`
}

func sendNextKey(w http.ResponseWriter, r *http.Request) {
	keys, err := db.PullNextKey()
	if err != nil {
		log.Fatal(err)
	}
	jsonKeys := JsonKeyPair{CodeName: keys.CodeName,
		PublicKey:  keys.PublicKey,
		PrivateKey: keys.PrivateKey,
	}
	json.NewEncoder(w).Encode(jsonKeys)
}

func Server() {
	router := mux.NewRouter()
	router.HandleFunc("/getkey", sendNextKey).Methods("GET")

	ip := "127.0.0.1"
	port := "9001"
	serverStr := fmt.Sprintf("%s:%s", ip, port)

	log.Fatal(http.ListenAndServe(serverStr, router))
}
