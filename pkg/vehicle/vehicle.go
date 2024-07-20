package vehicle

import (
	"log"
	"sync"

	"github.com/ArduCrow/go-msp/internal/msp"
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
	log.Println("GMSP-VEH: Vehicle initialized successfully")
	return &Vehicle{MspReader: mspReader, stopChan: make(chan struct{})}, nil
}

// Start the vehicle update loop, reads and updates vehicle states from the MSP
// connection.
func (v *Vehicle) Start() {
	log.Println("GMSP-VEH: Starting Vehicle:", v)
	v.wg.Add(1)

	go func() {
		defer v.wg.Done()
		for {
			select {
			case <-v.stopChan:
				log.Println("GMSP-VEH: STOP SIGNAL RECEIVED")
				return
			default:
				err := v.updateStates()
				if err != nil {
					log.Println("GMSP-VEH: Failed to read states:", err)
				}
			}
		}
	}()
}

// Stop the vehicle update loop and close the serial port connection. Shuts
// down all goroutines gracefully.
func (v *Vehicle) Stop() {
	log.Println("GMSP-VEH: Stopping Vehicle:", v)
	close(v.stopChan)
	v.wg.Wait()
	v.MspReader.Port.Close()
}

// Sets the raw RC values (channels, which are PWM values). Vehicle must be in
// MSP override mode.
func (v *Vehicle) SetChannels(channels []int) error {
	log.Println("GMSP-VEH: Setting channels:", channels)
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
