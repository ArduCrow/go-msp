package main

import (
	"go-msp/pkg/msp"
	"log"
)

func main() {
	mspReader, err := msp.NewMspReader("/dev/ttyACM0", 115200) // Use your serial port name and baud rate
	if err != nil {
		log.Fatalf("Failed to open serial port: %v", err)
	}
	for {
		attitude, err := mspReader.ReadAttitude()
		if err != nil {
			log.Fatalf("Failed to read attitude: %v", err)
		}
		if attitude != nil {
			log.Printf("Attitude: %v", attitude)
		}
	}
}
