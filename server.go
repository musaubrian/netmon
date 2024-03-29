package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	g "github.com/musaubrian/netmon/gno"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

func Server(ctx context.Context, tunn ngrok.Tunnel) {
	mux := http.NewServeMux()
	mux.HandleFunc("/logo", logo)
	mux.HandleFunc("/favicon", favicon)
	mux.HandleFunc("/lats", getLatencies)
	mux.HandleFunc("/", displayGraph)

	g.Log(g.INFO, "TUNNEL CREATED AT: "+tunn.URL())
	log.Fatal(http.Serve(tunn, mux))
}

func getLatencies(w http.ResponseWriter, r *http.Request) {
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

func logo(w http.ResponseWriter, r *http.Request) {
	f, err := getLogo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(f) < 1 {
		f = "./web/static/netmon.png"
	}
	http.ServeFile(w, r, f)
}

func favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/static/favicon.ico")
}

func createNgrokListener(ctx context.Context, token string) (ngrok.Tunnel, error) {

	tunn, err := ngrok.Listen(
		ctx,
		config.HTTPEndpoint(),
		ngrok.WithAuthtoken(token),
	)
	return tunn, err
}
