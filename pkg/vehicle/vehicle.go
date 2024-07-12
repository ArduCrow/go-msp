package vehicle

import (
	"fmt"
	"go-msp/pkg/msp"
	"sync"
)

type Vehicle struct {
	MspReader     *msp.MspReader
	ChannelValues []int
	Attitude      []float64
	stopChan      chan struct{}
	wg            sync.WaitGroup
}

// NewVehicle initializes a new Vehicle with the given serial port configuration.
func NewVehicle(portName string, baudRate int) (*Vehicle, error) {
	mspReader, err := msp.NewMspReader(portName, baudRate)
	if err != nil {
		return nil, err
	}
	fmt.Println("Vehicle initialized successfully")
	return &Vehicle{MspReader: mspReader, stopChan: make(chan struct{})}, nil
}

func (v *Vehicle) Start() {
	fmt.Println("Starting Vehicle:", v)
	v.wg.Add(1)

	go func() {
		defer v.wg.Done()
		for {
			select {
			case <-v.stopChan:
				fmt.Println("STOP SIGNAL RECEIVED")
				return
			default:
				err := v.updateStates()
				if err != nil {
					fmt.Println("Failed to read states:", err)
				}
				fmt.Println("Attitude:", v.Attitude)
			}
		}
	}()
}

func (v *Vehicle) Stop() {
	fmt.Println("Stopping Vehicle:", v)
	close(v.stopChan)
	v.wg.Wait()
	v.MspReader.Port.Close()
}

func (v *Vehicle) SetChannels(channels []int) error {
	fmt.Println("Setting channels:", channels)
	_, err := v.MspReader.SendRawRC(channels)
	if err != nil {
		return err
	}
	return nil
}

func (v *Vehicle) updateStates() error {
	err := v.readAttitude()
	if err != nil {
		return err
	}
	err = v.readChannels()
	if err != nil {
		return err
	}
	return nil
}

func (v *Vehicle) readAttitude() error {
	att, err := v.MspReader.ReadAttitude()
	if err != nil {
		return err
	}
	if att != nil {
		v.Attitude = att
	}
	return nil
}

func (v *Vehicle) readChannels() error {
	ch, err := v.MspReader.ReadRcChannels()
	if err != nil {
		return err
	}
	if ch != nil {
		v.ChannelValues = ch
	}
	return nil
}
