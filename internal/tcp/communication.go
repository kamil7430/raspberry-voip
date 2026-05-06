package tcp

import (
	"log"
	"net"
	"time"

	"github.com/kamil7430/raspberry-voip/internal/audio"
)

const (
	bufferSize  = 8196
	dialingTime = 15 * time.Second
)

func receiveAndPlay(conn net.Conn, audio *audio.AudioHandler) {
	receiveBuffer := make([]byte, bufferSize) // TODO: get bufferSize from audio

	n, err := conn.Read(receiveBuffer)
	if err != nil {
		log.Printf("Error reading from connection: %s\n", err)
		return
	}

	audio.In <- receiveBuffer[:n]
}

func sendFromAudioBuffer(conn net.Conn, audio *audio.AudioHandler) {
	_, err := conn.Write(<-audio.Out)
	if err != nil {
		log.Printf("Error sending voice: %s\n", err)
	}
}
