package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
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

func Server(ctx context.Context, tunn ngrok.Tunnel) {
	mux := http.NewServeMux()
	mux.HandleFunc("/logs", getLogs)
	mux.HandleFunc("/", displayGraph)

	log.Println("TUNNEL CREATED AT:", tunn.URL())
	log.Fatal(http.Serve(tunn, mux))
}

// extract NGROK initialization to function
// So as I can access the URL [tunn.URL()]
func createNgrokListener(ctx context.Context, token string) (ngrok.Tunnel, error) {
	tunn, err := ngrok.Listen(
		ctx,
		config.HTTPEndpoint(),
		ngrok.WithAuthtoken(token),
	)
	return tunn, err
}
