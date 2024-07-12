package msp

import (
	"fmt"

	"github.com/tarm/serial"
)

const MSP_ATTITUDE = 108
const MSP_RC = 105

// Reads and writes Multi wii serial protocol (MSP) from a serial port
// in order to communicate with a flight controller.
type MspReader struct {
	Port           *serial.Port
	RcChannels     []int
	ActiveChannels int
	MsgCodes       map[string]int
}

// NewMspReader initializes a new MspReader with the given serial port configuration.
func NewMspReader(portName string, baudRate int) (*MspReader, error) {
	c := &serial.Config{Name: portName, Baud: baudRate}
	port, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}
	fmt.Println("Serial port opened successfully")
	return &MspReader{Port: port}, nil
}

// SendRawMsg sends a raw MSP message through the serial port
func (mr *MspReader) SendRawMsg(code int, data []byte) (int, error) {
	var buf []byte
	if code < 255 { // MSP V1
		buf = make([]byte, 6+len(data))
		buf[0] = 36 // $
		buf[1] = 77 // M
		buf[2] = 60 // <
		buf[3] = byte(len(data))
		buf[4] = byte(code)
		checksum := buf[3] ^ buf[4]
		for i, d := range data {
			buf[5+i] = d
			checksum ^= d
		}
		buf[len(buf)-1] = checksum
	} else {
		// MSP V2 not implemented in this example
		return 0, fmt.Errorf("MSP V2 not supported in this example")
	}

	n, err := mr.Port.Write(buf)
	if err != nil {
		return n, err
	}
	return n, nil
}

// ReadAttitude requests and reads the vehicle's attitude (roll, pitch, yaw)
func (mr *MspReader) ReadAttitude() ([]float64, error) {
	mr.Port.Flush()
	_, err := mr.SendRawMsg(MSP_ATTITUDE, nil)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 14) // 6 bytes header + 6 bytes data + 2 bytes for potential MSP V2
	n, err := mr.Port.Read(buf)
	if err != nil {
		return nil, err
	}

	if n < 12 { // Not enough data
		return nil, nil
	}

	// Assuming MSP V1 and data starts at index 5
	data := buf[5 : 5+6]
	roll := float64(int16(data[0])|(int16(data[1])<<8)) / 10.0
	pitch := float64(int16(data[2])|(int16(data[3])<<8)) / 10.0
	yaw := float64(int16(data[4]) | (int16(data[5]) << 8))

	return []float64{roll, pitch, yaw}, nil
}

func (mr *MspReader) ReadRcChannels() ([]int, error) {
	mr.Port.Flush()
	_, err := mr.SendRawMsg(MSP_RC, nil)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 32) // Adjust buffer size as needed
	n, err := mr.Port.Read(buf)
	if err != nil {
		return nil, err
	}

	if n < 16 { // Not enough data
		return nil, nil
	}

	// Assuming MSP V1 and data starts at index 5
	data := buf[5 : 5+16]
	channels := make([]int, 8)
	for i := 0; i < 8; i++ {
		channels[i] = int(int16(data[i*2]) | (int16(data[i*2+1]) << 8))
	}

	return channels, nil
}

// ListenAndPrint listens on the serial port and prints received data.
// func (mr *MspReader) ListenAndPrint() {
// 	buf := make([]byte, 128) // Adjust buffer size as needed
// 	println("Listening for data...")
// 	for {
// 		n, err := mr.Port.Read(buf)
// 		fmt.Println("Read data")
// 		if err != nil {
// 			fmt.Printf("Error reading from port: %v", err)
// 			continue
// 		}
// 		fmt.Printf("Received: %q\n", buf[:n])
// 	}
// }
