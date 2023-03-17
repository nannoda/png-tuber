package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/MarkKremer/microphone"
	"github.com/faiface/beep/wav"
	utils "nannnoda.com/pngtuber/internal"
)

// import signal.go from the same directory as this file

func main() {

	fmt.Println("Recording. Press Ctrl-C to stop.")

	err := microphone.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer microphone.Terminate()

	stream, format, err := microphone.OpenDefaultStream(44100, 2)
	if err != nil {
		log.Fatal(err)
	}
	// Close the stream at the end if it hasn't already been
	// closed explicitly.
	defer stream.Close()

	filename := "output"
	if !strings.HasSuffix(filename, ".wav") {
		filename += ".wav"
	}
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Stop the stream when the user tries to quit the program.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	go func() {
		<-sig
		stream.Stop()
		stream.Close()
	}()

	stream.Start()

	// Encode the stream. This is a blocking operation because
	// wav.Encode will try to drain the stream. However, this
	// doesn't happen until stream.Close() is called.
	err = wav.Encode(f, stream, format)
	if err != nil {
		log.Fatal(err)
	}

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
