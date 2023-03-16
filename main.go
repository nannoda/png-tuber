package main

import (
	"fmt"

	utils "nannnoda.com/pngtuber/internal"
)

// import signal.go from the same directory as this file

func main() {

	remoteIdChan := make(chan string)
	localIdChan := make(chan string)
	onMessageChan := make(chan string)

	println("Starting server...")
	go utils.ServeWeb(remoteIdChan, localIdChan)

	println("Waiting for remoteId...")
	remoteId := <-remoteIdChan

	go utils.CreatePeerConnection(remoteId, localIdChan, onMessageChan)

	for {
		select {
		case msg := <-onMessageChan:
			fmt.Println(msg)
		}
	}
}
