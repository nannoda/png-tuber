package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func handle404(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "404 not found: %q", r.URL.Path)
}

func handleIndexHTML(w http.ResponseWriter, r *http.Request) {
	filePath := "web/index.html"

	dat, err := os.ReadFile(filePath)
	if err != nil {
		handle404(w, r)
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func handleIndexJS(w http.ResponseWriter, r *http.Request) {
	filePath := "web/index.js"

	dat, err := os.ReadFile(filePath)
	if err != nil {
		handle404(w, r)
	}
	w.Header().Set("Content-Type", "text/javascript")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func handleRemoteIdApi(removeId chan string, w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("handleRemoteIdApi POST err: ", err)
			return
		}
		bodyStr := "" + string(body)
		// fmt.Printf("handleRemoteIdApi POST body: [%s]\n", bodyStr)
		removeId <- bodyStr
		fmt.Printf("sent remoteId to removeId channel.\n")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	} else {
		fmt.Println("Method not allowed: ", r.Method)
	}
}

func ServeWeb(removeId chan string) {
	http.HandleFunc("/", handle404)
	http.HandleFunc("/index.html", handleIndexHTML)
	http.HandleFunc("/index.js", handleIndexJS)
	http.HandleFunc("/api/remoteId", func(w http.ResponseWriter, r *http.Request) {
		handleRemoteIdApi(removeId, w, r)
	})

	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf(":%d", 6060),
			nil,
		),
	)
}
