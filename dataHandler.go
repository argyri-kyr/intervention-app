package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var answers []Answer

func HandleAnswer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	question := r.FormValue("question")
	answer := r.FormValue("answer")

	answers = append(answers, Answer{Question: question, Answer: answer})

	jsonData, err := json.Marshal(answers)
	if err != nil {
		log.Fatalf("Failed to encode answers to JSON: %v", err)
	}

	err = os.WriteFile("answers.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write answers to JSON file: %v", err)
	}

	log.Println("Stored answer")
}

func loadSessionData(sessions []Sessions, progression Progression, sessionNumber int) []SectorDetails {
	var details []SectorDetails
	for _, session := range sessions {
		if session.Session == sessionNumber {
			sectors := session.Sectors
			for _, sector := range sectors {
				var firstElement SectorEntry
				switch sector {
				case "phonological_awareness":
					if len(progression.PhonologicalAwareness) > 0 {
						firstElement = progression.PhonologicalAwareness[0]
					}
				case "decoding":
					if len(progression.Decoding) > 0 {
						firstElement = progression.Decoding[0]
					}
				case "reading_comprehension":
					if len(progression.ReadingComprehension) > 0 {
						firstElement = progression.ReadingComprehension[0]
					}
				case "memory":
					if len(progression.Memory) > 0 {
						firstElement = progression.Memory[0]
					}
				case "processing_speed":
					if len(progression.ProcessingSpeed) > 0 {
						firstElement = progression.ProcessingSpeed[0]
					}
				default:
					fmt.Printf("Unknown sector: %s\n", sector)
					continue
				}

				detail := SectorDetails{
					Sector:     sector,
					Section:    firstElement.Section,
					Difficulty: firstElement.Difficulty,
					GameNo:     firstElement.GameNo,
					Template:   firstElement.Template,
				}
				details = append(details, detail)
			}
		}
	}
	return details
}

func loadGameData(sectorDetails []SectorDetails) (map[int][]Game, error) {
	gameData := make(map[int][]Game)
	counter := 1
	for _, detail := range sectorDetails {
		//fmt.Println("Sector", detail.Sector)
		sector := detail.Sector
		section := detail.Section
		difficulty := detail.Difficulty
		game := detail.GameNo

		datafilePath := filepath.Join(basePath, sector, fmt.Sprintf("%d_%d_%d.json", section, difficulty, game))

		// Check if the path exists
		if _, err := os.Stat(datafilePath); os.IsNotExist(err) {
			fmt.Printf("File does not exist: %s\n", datafilePath)
			continue
		}

		file, err := os.Open(datafilePath)
		if err != nil {
			return nil, fmt.Errorf("file exists but failed to open: %w", err)
		}
		defer file.Close()

		var gameSlice []Game
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&gameSlice)
		if err != nil {
			return nil, fmt.Errorf("failed to decode JSON data: %w", err)
		}

		// Append games to the appropriate array in gameData
		gameData[counter] = append(gameData[counter], gameSlice...)
		counter++
	}

	return gameData, nil
}

func LoadData(sessionNumber int) ([]GameToHtml, error) {
	var gameDataToHtml []GameToHtml
	// Load sessions data
	sessionsFilePath := filepath.Join(basePath, "sessions.json")
	sessionsFile, err := os.Open(sessionsFilePath)
	if err != nil {
		return gameDataToHtml, fmt.Errorf("failed to open file: %w", err)
	}
	defer sessionsFile.Close()
	var sessions []Sessions
	bytes, err := io.ReadAll(sessionsFile)
	if err != nil {
		return gameDataToHtml, fmt.Errorf("failed to read file: %w", err)
	}
	err = json.Unmarshal(bytes, &sessions)
	if err != nil {
		return gameDataToHtml, fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}

	// Load progressions data
	progressionsFilePath := filepath.Join(basePath, "progression.json")
	progressionsFile, err := os.Open(progressionsFilePath)
	if err != nil {
		return gameDataToHtml, fmt.Errorf("failed to open file: %w", err)
	}
	defer progressionsFile.Close()

	var progression Progression
	bytes, err = io.ReadAll(progressionsFile)
	if err != nil {
		return gameDataToHtml, fmt.Errorf("failed to read file: %w", err)
	}
	err = json.Unmarshal(bytes, &progression)
	if err != nil {
		return gameDataToHtml, fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}

	// Process data (example: print it)
	sectorDetails := loadSessionData(sessions, progression, sessionNumber)
	//fmt.Println("sectorDetails = ", sectorDetails)

	// Load game data based on sector details
	gameData, err := loadGameData(sectorDetails)
	//fmt.Println("gameData = ", gameData)
	for i := 0; i < len(gameData); i++ {
		gameDataToHtml = append(gameDataToHtml, GameToHtml{Template: sectorDetails[i].Template, Games: gameData[i+1]})
	}

	if err != nil {
		return nil, fmt.Errorf("error loading game data: %w", err)
	}

	return gameDataToHtml, nil
}
