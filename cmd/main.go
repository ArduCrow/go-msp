package main

import (
	"go-msp/pkg/vehicle"
	"log"
	"time"
)

func main() {
	vehicle, err := vehicle.NewVehicle("/dev/ttyACM0", 115200)

	if err != nil {
		log.Fatalf("Failed to start vehicle: %v", err)
	}

	vehicle.Start()

	for i := 0; i < 10; i++ {
		log.Printf("Second %d", i)
		time.Sleep(time.Second)
	}

}
