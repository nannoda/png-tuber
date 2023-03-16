package main

import (
	"fmt"

	utils "nannnoda.com/pngtuber/internal"
)

// import signal.go from the same directory as this file

func main() {

	remoteIdChan := make(chan string)
	onMessageString := make(chan string)

	println("Starting server...")
	go utils.ServeWeb(remoteIdChan)

	println("Waiting for remoteId...")
	remoteId := <-remoteIdChan

	// fmt.Printf("remoteId: [%s]\n", remoteId)

	utils.CreatePeerConnection(remoteId, onMessageString)

	for {
		select {
		case msg := <-onMessageString:
			fmt.Println(msg)
		}
	}
}
