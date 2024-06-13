package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const gamesPerSession = 4

func storeAnswers(answers []Answer, sessionNumber int, gameNumber int) {
	log.Println("Storing answers for session", sessionNumber, "and game", gameNumber)
	jsonData, err := json.Marshal(answers)
	if err != nil {
		log.Fatalf("Failed to encode answers to JSON: %v", err)
	}
	filename := fmt.Sprintf("answers/%d_%d_answers.json", sessionNumber, gameNumber)
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write answers to JSON file: %v", err)
	}
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
	log.Printf("Handling answers for session %d and game %d\n", sessionNumber, gameNumber)

	var answers []Answer
	err := json.NewDecoder(r.Body).Decode(&answers)
	log.Println("Answers:", r.Body)
	if err != nil {
		log.Println("Failed to decode JSON body")
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
					log.Printf("Unknown sector: %s\n", sector)
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

func loadAllGameData(sectorDetails []SectorDetails) (map[int][]Game, error) {
	gameData := make(map[int][]Game)
	counter := 1
	for _, detail := range sectorDetails {
		//log.Println("Sector", detail.Sector)
		sector := detail.Sector
		section := detail.Section
		difficulty := detail.Difficulty
		game := detail.GameNo

		datafilePath := filepath.Join(basePath, sector, fmt.Sprintf("%d_%d_%d.json", section, difficulty, game))

		// Check if the path exists
		if _, err := os.Stat(datafilePath); os.IsNotExist(err) {
			log.Printf("File does not exist: %s\n", datafilePath)
			continue
		}

		file, err := os.Open(datafilePath)
		if err != nil {
			log.Fatalf("file exists but failed to open: %v", err)
		}
		defer file.Close()

		var gameSlice []Game
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&gameSlice)
		if err != nil {
			log.Fatalf("failed to decode JSON data: %v", err)
		}

		// Append games to the appropriate array in gameData
		gameData[counter] = append(gameData[counter], gameSlice...)
		counter++
	}

	return gameData, nil
}

func loadGameData(sectorDetails []SectorDetails) (map[int][]Game, error) {
	gameData := make(map[int][]Game)
	counter := 1
	for _, detail := range sectorDetails {
		//log.Println("Sector", detail.Sector)
		sector := detail.Sector
		section := detail.Section
		difficulty := detail.Difficulty
		maxExamples := detail.MaxExamples
		game := detail.GameNo

		datafilePath := filepath.Join(basePath, sector, fmt.Sprintf("%d_%d_%d.json", section, difficulty, game))

		// Check if the path exists
		if _, err := os.Stat(datafilePath); os.IsNotExist(err) {
			log.Printf("File does not exist: %s\n", datafilePath)
			continue
		}

		file, err := os.Open(datafilePath)
		if err != nil {
			log.Fatalf("file exists but failed to open: %v", err)
		}
		defer file.Close()

		var gameSlice []Game
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&gameSlice)
		if err != nil {
			log.Fatalf("failed to decode JSON data: %v", err)
		}

		selectedGames, _ := getRandomElements(gameSlice, maxExamples)
		// log.Println("\n\nselectedIndices = ", selectedIndices)
		// log.Println("\n\nrandomGames = ", selectedGames)

		// Append games to the appropriate array in gameData
		gameData[counter] = append(gameData[counter], selectedGames...)
		counter++
	}

	return gameData, nil
}

func loadAnswers(sessionNumber int, gameNumber int) ([]Answer, error) {
	log.Println("Loading answers for session", sessionNumber, "and game", gameNumber)
	filename := fmt.Sprintf("answers/%d_%d_answers.json", sessionNumber, gameNumber)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	var answers []Answer
	json.Unmarshal(byteValue, &answers)
	//log.Println("Answers for game", gameNumber, ":", answers)
	return answers, nil
}

func LoadData(sessionNumber int) ([]GameToHtml, error) {
	log.Println("Loading data for session number:", sessionNumber)
	var gameDataToHtml []GameToHtml
	// Load sessions data
	sessionsFilePath := filepath.Join(basePath, "sessions.json")
	sessionsFile, err := os.Open(sessionsFilePath)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer sessionsFile.Close()
	var sessions []Sessions
	bytes, err := io.ReadAll(sessionsFile)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	err = json.Unmarshal(bytes, &sessions)
	if err != nil {
		log.Fatalf("failed to unmarshal JSON data: %v", err)
	}

	// Load progressions data
	progressionsFilePath := filepath.Join(basePath, "progression.json")
	progressionsFile, err := os.Open(progressionsFilePath)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer progressionsFile.Close()

	var progression Progression
	bytes, err = io.ReadAll(progressionsFile)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	err = json.Unmarshal(bytes, &progression)
	if err != nil {
		log.Fatalf("failed to unmarshal JSON data: %v", err)
	}
	// Get answers, get score and determine next games
	if sessionNumber > 1 {
		//Load answers
		for gameNumber := 1; gameNumber <= gamesPerSession; gameNumber++ {
			log.Println("Checking game number ", gameNumber, " for session ", sessionNumber-1)
			log.Println("Sector: ", sessions[sessionNumber-2].Sectors[gameNumber-1])
			sector := sessions[sessionNumber-2].Sectors[gameNumber-1]
			currentTemplate := ""
			pa_gameNumber := 0
			de_gameNumber := 0
			rc_gameNumber := 0
			ps_gameNumber := 0
			me_gameNumber := 0
			switch sector {
			case "phonological_awareness":
				currentTemplate = progression.PhonologicalAwareness[pa_gameNumber].Template
				pa_gameNumber++
			case "decoding":
				currentTemplate = progression.Decoding[de_gameNumber].Template
				de_gameNumber++
			case "reading_comprehension":
				currentTemplate = progression.ReadingComprehension[rc_gameNumber].Template
				rc_gameNumber++
			case "processing_speed":
				currentTemplate = progression.ProcessingSpeed[ps_gameNumber].Template
				ps_gameNumber++
			case "memory":
				currentTemplate = progression.Memory[me_gameNumber].Template
				me_gameNumber++
			}

			// log.Println("Progression: ", progression[sessions[sessionNumber-1].Sectors[gameNumber]])
			// log.Println("Has progression: ", progression.)
			answers, err := loadAnswers(sessionNumber-1, gameNumber)
			if currentTemplate == "AudioRecord" {
				log.Println("Audio record template, exit")
				//TODO: Audio should be decoded from answers, operator to listen to it and mark the score
				for _, answer := range answers {
					filename := answer.MainwordText
					filename = strings.ReplaceAll(filename, "[", "")
					filename = strings.ReplaceAll(filename, "]", "")
					filename = strings.TrimSpace(filename)
					decodeAudio(answer.Answer, filename)
				}
				continue
			}
			if err != nil {
				log.Println("Error loading answers:", err)
				continue
			}

			//Get score
			score := getScore(sessions, progression, sessionNumber-1, gameNumber, answers)
			log.Println("Score for game", gameNumber, ":", score)
		}
		//Set progression per sector
	}
	// Process data (example: print it)
	//log.Println("sessions = ", sessions)
	//log.Println("progression = ", progression)
	//log.Println("sessionNumber = ", sessionNumber)
	//TODO: Progression per sector
	sectorDetails := loadSessionData(sessions, progression, sessionNumber)
	//log.Println("sectorDetails = ", sectorDetails)

	// Load game data based on sector details
	gameData, err := loadGameData(sectorDetails)
	//log.Println("gameData = ", gameData)
	for i := 0; i < len(gameData); i++ {
		gameDataToHtml = append(gameDataToHtml, GameToHtml{Template: sectorDetails[i].Template, Games: gameData[i+1]})
	}

	if err != nil {
		log.Fatalf("error loading game data: %v", err)
	}

	return gameDataToHtml, nil
}
