package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

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

	log.Printf("Attemping to open vscode at %v@%v -> %v", req.User, req.Host, req.Path)

	url := fmt.Sprintf("vscode://vscode-remote/ssh-remote+%s@%s%s", req.User, req.Host, req.Path)

	cmd := exec.Command(url)
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to open VSCode @ %s@%s%s: %s", req.User, req.Host, req.Path, err)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /open", handleOpen)

	addr := "127.0.0.1:7777"
	log.Printf("beamd listening on %s", addr)

	log.Fatal(http.ListenAndServe(addr, mux))
}
