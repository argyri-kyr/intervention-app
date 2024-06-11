package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func decodeAudio(base64data string) {
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
	err = os.WriteFile("output.ogg", data, 0644)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}
}

func storeAnswers(answers []Answer, sessionNumber int, gameNumber int) {
	// Get gameData
	//gameData, err := LoadData(1)
	jsonData, err := json.Marshal(answers)
	if err != nil {
		log.Fatalf("Failed to encode answers to JSON: %v", err)
	}
	filename := fmt.Sprintf("%d_%d_answers.json", sessionNumber, gameNumber)
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write answers to JSON file: %v", err)
	}
}

func getScore(answers []Answer) int {
	fmt.Println(answers[0].Question)
	fmt.Println(answers[0].Index)
	fmt.Println(answers[0].Answer)
	//decodeAudio(answers[0].Answer)
	return 0
}

func HandleAnswer(w http.ResponseWriter, r *http.Request) {
	sessionNumber := 0
	gameNumber := 0
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) > 3 {
		var err error
		sessionNumber, err = strconv.Atoi(parts[2])
		if err != nil {
			sessionNumber = 1
			// handle the error
		}
		gameNumber, err = strconv.Atoi(parts[3])
		if err != nil {
			gameNumber = 1
			// handle the error
		}
		// Perform actions with sessionNumber and gameNumber
	}
	fmt.Printf("Handling answers for session %d and game %d\n", sessionNumber, gameNumber)
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var answers []Answer
	err := json.NewDecoder(r.Body).Decode(&answers)
	fmt.Println(len(answers))

	if err != nil {
		http.Error(w, "Failed to decode JSON body", http.StatusBadRequest)
		return
	}
	storeAnswers(answers, sessionNumber, gameNumber)
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
					Sector:      sector,
					Section:     firstElement.Section,
					Difficulty:  firstElement.Difficulty,
					GameNo:      firstElement.GameNo,
					MaxExamples: firstElement.MaxExamples,
					Template:    firstElement.Template,
				}
				details = append(details, detail)
			}
		}
	}
	return details
}

func getRandomElements(gameSlice []Game, count int) ([]Game, []int) {
	// Check if gameSlice has less than count elements
	if len(gameSlice) < count {
		return gameSlice, nil
	}

	// Create a slice to store the selected games
	selectedGames := make([]Game, count)
	// Create a slice to store the indices of selected games
	selectedIndices := make([]int, count)

	// Create a map to track used indices
	usedIndices := make(map[int]bool)

	// Loop to select random elements
	for i := 0; i < count; i++ {
		// Generate a random index
		randIndex := rand.Intn(len(gameSlice))
		// Check if the index is already used
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

func loadGameData(sectorDetails []SectorDetails) (map[int][]Game, error) {
	gameData := make(map[int][]Game)
	counter := 1
	for _, detail := range sectorDetails {
		//fmt.Println("Sector", detail.Sector)
		sector := detail.Sector
		section := detail.Section
		difficulty := detail.Difficulty
		maxExamples := detail.MaxExamples
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

		selectedGames, _ := getRandomElements(gameSlice, maxExamples)
		// fmt.Println("\n\nselectedIndices = ", selectedIndices)
		// fmt.Println("\n\nrandomGames = ", selectedGames)

		// Append games to the appropriate array in gameData
		gameData[counter] = append(gameData[counter], selectedGames...)
		counter++
	}

	return gameData, nil
}

func LoadData(sessionNumber int) ([]GameToHtml, error) {
	fmt.Println("Loading data for session number:", sessionNumber)
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
	// Get answers, get score and determine next games
	if sessionNumber > 1 {

	}
	// Process data (example: print it)
	//fmt.Println("sessions = ", sessions)
	//fmt.Println("progression = ", progression)
	//fmt.Println("sessionNumber = ", sessionNumber)
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
