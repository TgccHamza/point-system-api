package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"
)

// decodeTime decodes a timestamp retrieved from the timeclock
func decodeTime(t uint32) time.Time {
	second := t % 60
	t = t / 60

	minute := t % 60
	t = t / 60

	hour := t % 24
	t = t / 24

	day := t%31 + 1
	t = t / 31

	month := t%12 + 1
	t = t / 12

	year := t + 2000

	return time.Date(int(year), time.Month(month), int(day), int(hour), int(minute), int(second), 0, time.UTC)
}

func main() {
	// Convert the hex string to bytes
	hexCode := "2F00323030000000000000000000000000000000000000000000011EBAFE2F050000000000000000"
	byteData, err := hex.DecodeString(hexCode)
	if err != nil {
		panic(err)
	}

	// Define the struct format
	type Record struct {
		UID      uint16
		UserID   [24]byte
		Status   uint8
		Timestamp [4]byte
		Punch    uint8
		Space    [8]byte
	}

	// Unpack the byte data into the struct
	var record Record
	buf := bytes.NewReader(byteData)
	err = binary.Read(buf, binary.LittleEndian, &record)
	if err != nil {
		panic(err)
	}

	// Clean the UserID by removing null bytes and converting to a string
	userIDClean := string(bytes.TrimRight(record.UserID[:], "\x00"))

	// Convert UserID to an integer
	var userIDNumber int
	_, err = fmt.Sscanf(userIDClean, "%d", &userIDNumber)
	if err != nil {
		panic(err)
	}

	// Decode the timestamp
	timestamp := binary.LittleEndian.Uint32(record.Timestamp[:])
	decodedTime := decodeTime(timestamp)

	// Print the unpacked data
	fmt.Println("UID:", record.UID)
	fmt.Println("User ID (clean):", userIDClean)
	fmt.Println("User ID (number):", userIDNumber)
	fmt.Println("Status:", record.Status)
	fmt.Println("Punch:", record.Punch)
	fmt.Println("Timestamp:", decodedTime)
}