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

var receiveBuffer = make([]byte, bufferSize) // TODO: get bufferSize from audio

func receiveAndPlay(conn net.Conn, audio *audio.AudioHandler) {
	_, err := conn.Read(receiveBuffer)
	if err != nil {
		log.Printf("Error reading from connection: %s\n", err)
	}
	audio.In <- receiveBuffer
}

func sendFromAudioBuffer(conn net.Conn, audio *audio.AudioHandler) {
	_, err := conn.Write(<-audio.Out)
	if err != nil {
		log.Printf("Error sending voice: %s\n", err)
	}
}
