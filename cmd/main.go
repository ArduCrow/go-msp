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

	for {
		println("Setting channels")
		vehicle.SetChannels([]int{888, 999, 1000, 1500, 1500, 1500, 1500, 1500})
	}
	// select {}

}
