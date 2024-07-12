package vehicle

import (
	"fmt"
	"go-msp/pkg/msp"
)

type Vehicle struct {
	MspReader     *msp.MspReader
	ChannelValues []int
	Attitude      []float64
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

func (v *Vehicle) UpdateAttitude() error {
	att, err := v.MspReader.ReadAttitude()
	if err != nil {
		return err
	}
	if att != nil {
		v.Attitude = att
	}
	return nil
}

func (v *Vehicle) UpdateChannels() error {
	ch, err := v.MspReader.ReadRcChannels()
	if err != nil {
		return err
	}
	if ch != nil {
		v.ChannelValues = ch
	}
	return nil
}