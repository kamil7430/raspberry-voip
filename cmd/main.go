package main

import "log"

func main() {
	log.Println("Starting RPi VoIP!")
	runHttpServer() // will be a go routine
}
