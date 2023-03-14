package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func get404(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}

func setContentType(fileExtension string, w http.ResponseWriter) {
	if fileExtension == "html" {
		w.Header().Set("Content-Type", "text/html")
	}
	if fileExtension == "css" {
		w.Header().Set("Content-Type", "text/css")
	}
	if fileExtension == "js" {
		w.Header().Set("Content-Type", "text/javascript")
	}
	if fileExtension == "png" {
		w.Header().Set("Content-Type", "image/png")
	}
	if fileExtension == "jpg" {
		w.Header().Set("Content-Type", "image/jpeg")
	}
	if fileExtension == "jpeg" {
		w.Header().Set("Content-Type", "image/jpeg")
	}
	if fileExtension == "gif" {
		w.Header().Set("Content-Type", "image/gif")
	}
	if fileExtension == "svg" {
		w.Header().Set("Content-Type", "image/svg+xml")
	}
}

func returnFile(filePath string, w http.ResponseWriter) {
	// file extension is the characters after the last dot
	fileExtension := filePath[strings.LastIndex(filePath, ".")+1:]
	dat, err := os.ReadFile(filePath)
	if err != nil {
		get404(w)
		// return err
	}
	setContentType(fileExtension, w)
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func handleHttp(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Request received: " + r.URL.Path + " " + r.Method + " " + r.RemoteAddr + " ")
	// remove the first slash
	filePath := r.URL.Path[1:]
	log.Default().Println("Path: " + filePath)
	returnFile(filePath, w)
}

func main() {
	port := 5100
	if len(os.Args) > 1 {
		portStr := os.Args[1]
		port = int(portStr[0])

	}

	http.HandleFunc("/", handleHttp)
	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(":%d", port), nil))
}
