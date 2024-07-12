package main

import (
	"go-msp/pkg/vehicle"
	"log"
)

func main() {
	vehicle, err := vehicle.NewVehicle("/dev/ttyACM0", 115200)

	if err != nil {
		log.Fatalf("Failed to start vehicle: %v", err)
	}

	vehicle.Start()
	select {}

}
