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
		err := vehicle.UpdateAttitude()
		if err != nil {
			log.Fatalf("Failed to read attitude: %v", err)
		}
		log.Printf("Attitude: %v", vehicle.Attitude)
		err = vehicle.UpdateChannels()
		if err != nil {
			log.Fatalf("Failed to read channels: %v", err)
		}
		log.Printf("Channels: %v", vehicle.ChannelValues)
	}
}
