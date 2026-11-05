package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Alasdair-Roddick/beam/internal/protocol"
)

func handleOpen(w http.ResponseWriter, r *http.Request) {
	var req protocol.OpenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("open request: user=%s host=%s path=%s", req.User, req.Host, req.Path)

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /open", handleOpen)

	addr := "127.0.0.1:7777"
	log.Printf("beamd listening on %s", addr)

	log.Fatal(http.ListenAndServe(addr, mux))
}
