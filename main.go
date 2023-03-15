package pngtuber

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pion/webrtc/v3"
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

	http.HandleFunc("/", handleHttp)
	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(":%d", port), nil))
}
