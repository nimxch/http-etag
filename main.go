package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type Cache struct {
	Checksum string
	FileData []byte
	mu       sync.RWMutex
}

var cache = Cache{}

func generateChecksum(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func getCachedChecksum() (string, []byte, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cacheData := cache.FileData

	// Read current file
	currentData, err := os.ReadFile("test.json")
	if err != nil {
		return "", nil, err
	}

	if string(cacheData) != string(currentData) {
		cache.FileData = currentData
		cache.Checksum = generateChecksum(currentData)
	}
	return cache.Checksum, cache.FileData, nil

}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=0, must-revalidate")
	resp := Response{
		Message: "Hello World",
		Data:    nil,
	}

	chk, data, err := getCachedChecksum()
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println(chk)

	if r.Header.Get("If-None-Match") == chk {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.Header().Set("ETag", chk)
	mapData := make(map[string]string)
	err = json.Unmarshal(data, &mapData)
	resp.Data = mapData
	jsonData, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprint(w, "Json Error")
		return
	}
	time.Sleep(1 * time.Second)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from any origin (for development)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, If-None-Match")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	serv := http.Server{
		Addr: ":8080",
	}
	handler := corsMiddleware(http.HandlerFunc(handleRequest))
	http.Handle("/", handler)
	fmt.Println(serv.ListenAndServe())
}
