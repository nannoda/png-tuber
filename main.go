package main

import (
	"fmt"
	"os"
	"time"

	"github.com/pion/webrtc/v3"

	// import signal.go from the same directory as this file
	"nannnoda.com/pngtuber/utils/signal"
)

func main() {
	// Everything below is the Pion WebRTC API! Thanks for using it ❤️.

	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if cErr := peerConnection.Close(); cErr != nil {
			fmt.Printf("cannot close peerConnection: %v\n", cErr)
		}
	}()

	// Set the handler for Peer connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		fmt.Printf("Peer Connection State has changed: %s\n", s.String())

		if s == webrtc.PeerConnectionStateFailed {
			// Wait until PeerConnection has had no network activity for 30 seconds or another failure. It may be reconnected using an ICE Restart.
			// Use webrtc.PeerConnectionStateDisconnected if you are interested in detecting faster timeout.
			// Note that the PeerConnection may come back from PeerConnectionStateDisconnected.
			fmt.Println("Peer Connection has gone to failed exiting")
			os.Exit(0)
		}
	})

	// Register data channel creation handling
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

		// Register channel opening handling
		d.OnOpen(func() {
			fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", d.Label(), d.ID())

			for range time.NewTicker(5 * time.Second).C {
				message := signal.RandSeq(15)
				fmt.Printf("Sending '%s'\n", message)

				// Send the message as text
				sendErr := d.SendText(message)
				if sendErr != nil {
					panic(sendErr)
				}
			}
		})

		// Register text message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
		})
	})

	// Wait for the offer to be pasted
	offer := webrtc.SessionDescription{}
	// base64Desc := "eyJ0eXBlIjoib2ZmZXIiLCJzZHAiOiJ2PTBcclxubz0tIDE0Njk3NDE2MzU5MjUyMTA4MjIgMiBJTiBJUDQgMTI3LjAuMC4xXHJcbnM9LVxyXG50PTAgMFxyXG5hPWdyb3VwOkJVTkRMRSAwXHJcbmE9ZXh0bWFwLWFsbG93LW1peGVkXHJcbmE9bXNpZC1zZW1hbnRpYzogV01TXHJcbm09YXBwbGljYXRpb24gNTE5MzcgVURQL0RUTFMvU0NUUCB3ZWJydGMtZGF0YWNoYW5uZWxcclxuYz1JTiBJUDQgNzUuMTI4LjIxNy4xMTRcclxuYT1jYW5kaWRhdGU6MjE5Njk5NjIzMyAxIHVkcCAyMTEzOTM3MTUxIDRkNjM1Yjg1LTczYjYtNDMzYy1hYzFhLTI5Y2UwMTc5NzhlZS5sb2NhbCA1MTkzNyB0eXAgaG9zdCBnZW5lcmF0aW9uIDAgbmV0d29yay1jb3N0IDk5OVxyXG5hPWNhbmRpZGF0ZTozMDk4NDQ2NDQ5IDEgdWRwIDE2Nzc3Mjk1MzUgNzUuMTI4LjIxNy4xMTQgNTE5MzcgdHlwIHNyZmx4IHJhZGRyIDAuMC4wLjAgcnBvcnQgMCBnZW5lcmF0aW9uIDAgbmV0d29yay1jb3N0IDk5OVxyXG5hPWljZS11ZnJhZzpaTURYXHJcbmE9aWNlLXB3ZDp6WEJpUmRMVXQybm1KZldOSzB4SlRuaGRcclxuYT1pY2Utb3B0aW9uczp0cmlja2xlXHJcbmE9ZmluZ2VycHJpbnQ6c2hhLTI1NiA3QzoyODo3NDo5MDo1MTpERDo1RDpERTpERDo4ODoyNDo3NjoyQzowNTpBRTpCNjpERjowQzo4RDozRjoyRjo4MDpGRjpCQzowRTozRjo3NjoxNDo4Mjo3MjowNTpFMlxyXG5hPXNldHVwOmFjdHBhc3NcclxuYT1taWQ6MFxyXG5hPXNjdHAtcG9ydDo1MDAwXHJcbmE9bWF4LW1lc3NhZ2Utc2l6ZToyNjIxNDRcclxuIn0="
	// signal.Decode(signal.MustReadStdin(), &offer)
	signal.Decode(base64Desc, &offer)

	// Set the remote SessionDescription
	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		panic(err)
	}

	// Create an answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	// Output the answer in base64 so we can paste it in browser
	fmt.Println(signal.Encode(*peerConnection.LocalDescription()))

	// Block forever
	select {}
}
