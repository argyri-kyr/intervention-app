package main

import (
	"encoding/base64"
	"log"
	"os"
	"strings"
)

func decodeAudio(base64data string, filename string) {
	// Split the base64 string. The actual data is after the comma.
	splitData := strings.Split(base64data, ",")
	if len(splitData) != 2 {
		log.Fatalf("Invalid base64 data")
	}

	// Decode the base64 string
	data, err := base64.StdEncoding.DecodeString(splitData[1])
	if err != nil {
		log.Fatalf("Failed to decode base64 string: %v", err)
	}

	// Write the data to an ogg file
	filename = strings.TrimSpace(filename)
	err = os.WriteFile("answers/audio/"+filename+".ogg", data, 0644)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}
}