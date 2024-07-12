package vehicle

import (
	"fmt"
	"go-msp/pkg/msp"
)

// Vehicle represents a vehicle that can be controlled by a remote controller.
type Vehicle struct {
	MspReader *msp.MspReader
}

// NewVehicle initializes a new Vehicle with the given serial port configuration.
func NewVehicle(portName string, baudRate int) (*Vehicle, error) {
	mspReader, err := msp.NewMspReader(portName, baudRate)
	if err != nil {
		return nil, err
	}
	fmt.Println("Vehicle initialized successfully")
	return &Vehicle{MspReader: mspReader}, nil
}

func (v *Vehicle) ReadAttitude() ([]float64, error) {
	return v.MspReader.ReadAttitude()
}