package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func createSessionFolder(sessionNumber int) error {
	folderPath := fmt.Sprintf("saved_data/session%d", sessionNumber)

	// Check if the folder already exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// Create the folder and all necessary parent directories
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create folder: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to access folder: %v", err)
	}

	return nil
}

func deleteOldFile(filename string) {
	// Check if the file exists
	if _, err := os.Stat(filename); err == nil {
		// File exists, delete it
		if err := os.Remove(filename); err != nil {
			fmt.Printf("Failed to delete file: %v\n", err)
			return
		}
	}

	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		return
	}
	defer file.Close()

	fmt.Printf("File '%s' created successfully.\n", filename)
}

// Function to save session game number to a file
func saveSessionSectorLevelToFile(savedSectors map[string]bool, filename string, sector string, SectorLevel int, score float64) {
	// Create or open the file
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		return
	}
	defer file.Close()

	// Read existing JSON data from file
	var data []SessionSectorLevel
	bytes, err := os.ReadFile(filename)
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("Failed to read file: %v\n", err)
		return
	}

	// Unmarshal JSON data if file exists
	if len(bytes) > 0 {
		err = json.Unmarshal(bytes, &data)
		if err != nil {
			fmt.Printf("Failed to unmarshal JSON data: %v\n", err)
			return
		}
	}

	// Append new data to the existing slice or create a new slice if empty
	data = append(data, SessionSectorLevel{
		Sector:      sector,
		SectorLevel: SectorLevel,
		Score:       score,
	})

	// Rewind file pointer to the beginning for writing
	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Printf("Failed to seek file: %v\n", err)
		return
	}

	// Encode the data as JSON
	encodedData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Printf("Failed to encode JSON: %v\n", err)
		return
	}

	// Write JSON data to file
	err = os.WriteFile(filename, encodedData, 0644)
	if err != nil {
		fmt.Printf("Failed to write JSON data to file: %v\n", err)
		return
	}

	savedSectors[sector] = true
	// fmt.Println("Data saved successfully.")
}

// Function to check and save missing sectors
func checkAndSaveMissingSectors(savedSectors map[string]bool, sessionNumber int, pa_gameNumber int, de_gameNumber int, rc_gameNumber int, ps_gameNumber int, me_gameNumber int) {
	for _, sector := range allSectors {
		// Check if sector has been saved already
		if savedSectors[sector] {
			// fmt.Printf("Sector '%s' has already been saved for this session\n", sector)
			continue
		}
		switch sector {
		case "phonological_awareness":
			score := getSectorScore(sessionNumber, sector)
			saveSessionSectorLevelToFile(savedSectors, fmt.Sprintf("saved_data/session%d/SectorLevel.json", sessionNumber), sector, pa_gameNumber, score)

		case "decoding":
			score := getSectorScore(sessionNumber, sector)
			saveSessionSectorLevelToFile(savedSectors, fmt.Sprintf("saved_data/session%d/SectorLevel.json", sessionNumber), sector, de_gameNumber, score)

		case "reading_comprehension":
			score := getSectorScore(sessionNumber, sector)
			saveSessionSectorLevelToFile(savedSectors, fmt.Sprintf("saved_data/session%d/SectorLevel.json", sessionNumber), sector, rc_gameNumber, score)

		case "processing_speed":
			score := getSectorScore(sessionNumber, sector)
			saveSessionSectorLevelToFile(savedSectors, fmt.Sprintf("saved_data/session%d/SectorLevel.json", sessionNumber), sector, ps_gameNumber, score)

		case "memory":
			score := getSectorScore(sessionNumber, sector)
			saveSessionSectorLevelToFile(savedSectors, fmt.Sprintf("saved_data/session%d/SectorLevel.json", sessionNumber), sector, me_gameNumber, score)
		}
	}
}
