package main

import (
	"go-msp/pkg/vehicle"
	"log"
)

func main() {
	vehicle, err := vehicle.NewVehicle("/dev/ttyACM0", 115200)

	if err != nil {
		log.Fatalf("Failed to open serial port: %v", err)
	}
	for {
		attitude, err := vehicle.ReadAttitude()
		if err != nil {
			log.Fatalf("Failed to read attitude: %v", err)
		}
		if attitude != nil {
			log.Printf("Attitude: %v", attitude)
		}
	}
}
