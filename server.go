package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

func getLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res, err := json.Marshal(dRecs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(res)
}

func displayGraph(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("./web/index.html")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer f.Close()
	w.Header().Set("Content-Type", "text/html")

	_, err = io.Copy(w, f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Server() {
	http.HandleFunc("/logs", getLogs)
	http.HandleFunc("/", displayGraph)
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
