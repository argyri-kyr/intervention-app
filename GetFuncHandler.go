package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
)

func getRandomElements(gameSlice []Game, count int, sector string, sessionNumber int) ([]Game, []int) {
	// Check if gameSlice has less than count elements
	if len(gameSlice) < count {
		count = len(gameSlice)
		fmt.Println("The max number of games you want to serve is more that the available. So I will serve as mush as I have.")
		// return gameSlice, nil
	}

	// Create a slice to store the selected games
	selectedGames := make([]Game, count)
	// Create a slice to store the indices of selected games
	selectedIndices := make([]int, count)
	// Create a map to track used indices
	usedIndices := make(map[int]bool)
	// Array that containes Indices(Example number) that I used in the last session of this sector. I want to exclude these Indices from the new ones

	// var retrievedIndices []int
	// if sessionNumber == 1 {
	// 	retrievedIndices = []int{}
	// } else {
	// 	filename := fmt.Sprintf("saved_data/session%d/SectorSelectedIndices.json", sessionNumber)
	// 	retrievedIndices = getSelectedIndicesFromFile(filename, sector)
	// }

	// Loop to select random elements
	for i := 0; i < count; i++ {
		// Generate a random index
		randIndex := rand.Intn(len(gameSlice))
		// Check if the index is already used or has used in the PREVIOUS session of this sector
		// for usedIndices[randIndex] || (len(retrievedIndices) > 0 && contains(retrievedIndices, randIndex)) {
		for usedIndices[randIndex] {
			randIndex = rand.Intn(len(gameSlice))
		}
		// Mark the index as used
		usedIndices[randIndex] = true
		// Store the selected game and its index
		selectedGames[i] = gameSlice[randIndex]
		selectedIndices[i] = randIndex
	}

	return selectedGames, selectedIndices
}

// Function to get the sector level for a specific sector from a file
func getSectorLevel(sessionNumber int, sector string) int {
	sectorLevel := 0 // Default to 0 if sector is not found

	for sessionNo := 1; sessionNo < sessionNumber; sessionNo++ {
		filename := fmt.Sprintf("saved_data/session%d/SectorLevel.json", sessionNo)

		// Read the JSON file
		fileContent, err := os.ReadFile(filename)
		if err != nil {
			log.Printf("Failed to read file: %v\n", err)
			continue
		}

		// Unmarshal the JSON data into a slice of SessionSectorLevel
		var data []SessionSectorLevel
		err = json.Unmarshal(fileContent, &data)
		if err != nil {
			log.Printf("Failed to unmarshal JSON data: %v\n", err)
			continue
		}

		// Iterate through the slice to find the sector
		for _, session := range data {
			if session.Sector == sector {
				sectorLevel = session.SectorLevel
			}
		}
	}

	// If sector not found in any session, log a message
	if sectorLevel == 0 {
		log.Printf("Sector '%s' not found in any session up to session number %d\n", sector, sessionNumber)
	}

	return sectorLevel
}

func getSectorScore(sessionNumber int, sector string) float64 {
	sectorscore := -1.0 // Default to 0 if sector is not found

	for sessionNo := 1; sessionNo < sessionNumber; sessionNo++ {
		filename := fmt.Sprintf("saved_data/session%d/SectorLevel.json", sessionNo)

		// Read the JSON file
		fileContent, err := os.ReadFile(filename)
		if err != nil {
			log.Printf("Failed to read file: %v\n", err)
			continue
		}

		// Unmarshal the JSON data into a slice of SessionSectorLevel
		var data []SessionSectorLevel
		err = json.Unmarshal(fileContent, &data)
		if err != nil {
			log.Printf("Failed to unmarshal JSON data: %v\n", err)
			continue
		}

		// Iterate through the slice to find the sector
		for _, session := range data {
			if session.Sector == sector {
				sectorscore = session.Score
			}
		}
	}

	// If sector not found in any session, log a message
	if sectorscore == -1.0 {
		log.Printf("Sector '%s' not found in any session up to session number %d\n", sector, sessionNumber)
	}

	return sectorscore
}

func getSessionData(sessionNumber int, filename string) *Sessions {
	// Read the JSON file
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		return nil
	}

	// Unmarshal the JSON data into a slice of Session
	var sessions []Sessions
	err = json.Unmarshal(fileContent, &sessions)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON data: %v\n", err)
		return nil
	}

	// Iterate through the sessions to find the matching session number
	for _, session := range sessions {
		if session.Session == sessionNumber {
			return &session
		}
	}

	// Return nil if the session is not found
	fmt.Printf("Session %d not found\n", sessionNumber)
	return nil
}

func getSectorNameByIndex(sessionNumber int, sectorIndex int, filename string) string {
	sessionData := getSessionData(sessionNumber, filename)
	if sessionData == nil {
		return ""
	}

	// Check if the sector index is within bounds
	if sectorIndex < 0 || sectorIndex >= len(sessionData.Sectors) {
		fmt.Printf("Sector index %d is out of bounds for session %d\n", sectorIndex, sessionNumber)
		return ""
	}

	return sessionData.Sectors[sectorIndex]
}

func getSectorDetails(currentSector string, sectorsDetails []SectorDetails) (*SectorDetails, error) {
	for _, details := range sectorsDetails {
		if details.Sector == currentSector {
			return &details, nil
		}
	}
	return nil, fmt.Errorf("sector '%s' not found in sectors details", currentSector)
}

func getSelectedIndicesFromFile(filename string, sector string) []int {
	// Try to open the file
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		return []int{}
	}
	defer file.Close()

	// Read the file content into a byte slice
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		return []int{}
	}

	// Unmarshal the JSON data
	var data SectorSelectedIndices
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		log.Printf("Failed to unmarshal JSON data: %v", err)
		return []int{}
	}

	// Check if the sector matches
	if data.Sector == sector {
		return data.SelectedIndices
	}

	log.Printf("Sector '%s' not found in file '%s'", sector, filename)
	return []int{}
}
