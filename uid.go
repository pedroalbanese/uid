package main

import (
	"bufio"
	"encoding/base32"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	// Parse command line flags
	decodeFlag := flag.String("d", "", "Decode the provided UID")
	flag.Parse()

	if *decodeFlag != "" {
		if *decodeFlag == "-" {
			// Read from stdin for decoding
			reader := bufio.NewReader(os.Stdin)
			encodedUID, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading encoded UID from stdin:", err)
				return
			}
			encodedUID = strings.TrimSpace(encodedUID)

			// Decode provided UID
			decodedUID, err := decodeUID(encodedUID)
			if err != nil {
				fmt.Println("Error decoding UID:", err)
				return
			}
			fmt.Println("UID:", encodedUID)
			fmt.Println("TMP:", decodedUID)
		} else {
			// Decode provided UID
			decodedUID, err := decodeUID(*decodeFlag)
			if err != nil {
				fmt.Println("Error decoding UID:", err)
				return
			}
			fmt.Println("UID:", *decodeFlag)
			fmt.Println("TMP:", decodedUID)
		}
	} else {
		// Encode UID
		encodedUID := encodeUID()
		fmt.Println(encodedUID)
	}
}

func encodeUID() string {
	// Get MAC address
	mac, err := getMACAddress()
	if err != nil {
		fmt.Println("Error obtaining MAC address:", err)
		return ""
	}

	// Reverse MAC address byte order
	reverseMAC(mac)

	// Get current timestamp
	timestamp := time.Now().Unix()

	// Create UID based on MAC and timestamp
	uidBytes := make([]byte, 10)
	copy(uidBytes[:6], mac)
	binary.LittleEndian.PutUint32(uidBytes[6:], uint32(timestamp))

	// Encode UID bytes to base32 string
	encoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	uid := encoder.EncodeToString(uidBytes)

	return uid
}

func decodeUID(encodedUID string) (string, error) {
	// Decode UID from base32 string to bytes
	decoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	uidBytes, err := decoder.DecodeString(encodedUID)
	if err != nil {
		return "", err
	}

	// Extract MAC address from UID bytes
	mac := uidBytes[:6]

	// Reverse MAC address byte order
	reverseMAC(mac)

	// Extract timestamp from UID bytes
	timestamp := int64(binary.LittleEndian.Uint32(uidBytes[6:]))

	// Convert timestamp to time
	t := time.Unix(timestamp, 0).Local()

	// Convert MAC address to string
	macString := fmt.Sprintf("%x", mac)

	// Format timestamp as desired format
	tFormatted := t.Format("2006-01-02 15:04:05 -0700")

	return fmt.Sprintf("%s\nMAC: %s", tFormatted, macString), nil
}

func getMACAddress() ([]byte, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		if len(iface.HardwareAddr) == 6 {
			return iface.HardwareAddr, nil
		}
	}

	return nil, fmt.Errorf("MAC address not found")
}

func reverseMAC(mac []byte) {
	for i := 0; i < len(mac)/2; i++ {
		j := len(mac) - i - 1
		mac[i], mac[j] = mac[j], mac[i]
	}
}
