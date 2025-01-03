package msp

import (
	"fmt"
	"log"
	"sync"

	"github.com/tarm/serial"
)

// Reads and writes Multi wii serial protocol (MSP) from a serial port
// in order to communicate with a flight controller.
type MspReader struct {
	Port           *serial.Port
	RcChannels     []int
	ActiveChannels int
	MsgCodes       map[string]int
	mu             sync.Mutex
}

// Initializes a new MspReader with the given serial port configuration.
func NewMspReader(portName string, baudRate int) (*MspReader, error) {
	c := &serial.Config{Name: portName, Baud: baudRate}
	port, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}
	log.Println("GMSP-MSP: Serial port opened successfully")
	return &MspReader{Port: port}, nil
}

// Sends a raw MSP message through the serial port
func (mr *MspReader) SendRawMsg(code int, data []byte) (int, error) {
	mr.Port.Flush()
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
		// MSP V2 not implemented
		return 0, fmt.Errorf("GMSP-MSP: MSP V2 not supported")
	}

	n, err := mr.Port.Write(buf)
	if err != nil {
		return n, err
	}
	return n, nil
}

// Send raw RC channel values to flight controller
func (mr *MspReader) SendRawRC(data []int) (int, error) {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	// Ensure the data slice contains 8 RC channels
	// if len(data) != 8 {
	// 	return 0, fmt.Errorf("GMSP-MSP: Expected 8 RC channels, got %d", len(data))
	// }

	// Convert RC channel values to bytes (low-byte/high-byte order)
	byteData := make([]byte, 16) // 8 channels * 2 bytes each
	for i, val := range data {
		byteData[i*2] = byte(val & 0xFF)          // Low byte
		byteData[i*2+1] = byte((val >> 8) & 0xFF) // High byte
	}

	// Send the message using the existing SendRawMsg method
	return mr.SendRawMsg(200, byteData) // 200 = MSP_SET_RAW_RC
}

// Requests and reads the vehicle's attitude (roll, pitch, yaw)
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

// Requests and reads the vehicle's RC channels
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

// Listen for and receive messages from the serial port. Uses a mutex to prevent
// concurrent writes to the serial port. Returns if an error occurs while reading.
func (mr *MspReader) receiveRawMessage() ([]byte, error) {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	buf := make([]byte, 32) // Adjust buffer size as needed
	n, err := mr.Port.Read(buf)
	if err != nil {
		log.Println("GMSP-MSP: Failed to read from serial port:", err)
		return nil, err
	}

	if n < 6 {
		log.Println("GMSP-MSP: Not enough data received")
		return nil, nil
	}

	// Assuming MSP V1 and data starts at index 5
	data := buf[5 : 5+6]
	log.Printf("GMSP-MSP: Received data: %v", data)
	return data, nil
}

// Check if the flight controller is ready to arm and fly
func (mr *MspReader) CheckReady() bool {
	// 3 - msp version
	// 61 - msp arming config
	// 101 - MSP status

	_, err := mr.SendRawMsg(3, nil)
	if err != nil {
		log.Println("GMSP-MSP: Check ready - failed to communicate with FC:", err)
		return false
	}
	data, err := mr.receiveRawMessage()

	switch {
	case err != nil:
		log.Println("GMSP-MSP: Failed to get FC state:", err)
		return false
	case data == nil:
		log.Println("GMSP-MSP: Check ready - no data received")
		return false
	case len(data) < 7:
		log.Println("GMSP-MSP: Check ready - not enough data received")
		return false
	case data[6] != 100:
		log.Println("GMSP-MSP: Flight controller not ready", data[6])
		return false
	}

	mspFcVersion := int(data[6])
	log.Println("GMSP-MSP: FC ready, MSP version:", mspFcVersion)
	return true
}
