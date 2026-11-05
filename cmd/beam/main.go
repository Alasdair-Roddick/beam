package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"syscall"
	"time"

	"github.com/Alasdair-Roddick/beam/internal/protocol"
)

func main() {
	log.SetFlags(0)

	url := "http://127.0.0.1:7777/open"

	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to get the current user: %v", err)
	}
	currentHost, err := os.Hostname()
	if err != nil {
		log.Fatalf("Failed to get the Hostname: %v", err)
	}
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get the path: %v", err)
	}

	contentBody := protocol.OpenRequest{
		User: currentUser.Username,
		Host: currentHost,
		Path: dir,
	}

	data, err := json.Marshal(contentBody)

	if err != nil {
		log.Fatalf("Error marshaling to JSON: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if errors.Is(err, syscall.ECONNREFUSED) {
			log.Fatalf("No connection to beamd (%s at %s), is it running on the remote system?", syscall.ECONNREFUSED, url)
		}
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("beamd returned %d: %s", resp.StatusCode, body)
	}

	fmt.Println("sent to beamd")
}
