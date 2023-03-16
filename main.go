package main

import (
	"fmt"

	"github.com/pion/mediadevices"
	_ "github.com/pion/mediadevices/pkg/driver/microphone"
	utils "nannnoda.com/pngtuber/internal"
)

// import signal.go from the same directory as this file

func main() {

	stream, err := mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
		Audio: func(c *mediadevices.MediaTrackConstraints) {

		},
	})
	if err != nil {
		panic(err)
	}
	track := stream.GetAudioTracks()[0]
	audioTrack := track.(*mediadevices.AudioTrack)
	defer audioTrack.Close()

	audioReader := audioTrack.NewReader(false)

	chunk, release, _ := audioReader.Read()
	defer release()

	fmt.Println(chunk.ChunkInfo())

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
