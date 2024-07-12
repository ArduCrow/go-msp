package main

import (
	"bufio"
	"log"
	"os"

	"github.com/ArduCrow/gomsp/pkg/vehicle"
)

func main() {
	vehicle, err := vehicle.NewVehicle("/dev/ttyACM0", 115200)

	if err != nil {
		log.Fatalf("Failed to start vehicle: %v", err)
	}

	vehicle.Start()

	// press ENTER to terminate and close vehicle gracefully
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		vehicle.Stop()
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error waiting for input: %v", err)
	}

}
