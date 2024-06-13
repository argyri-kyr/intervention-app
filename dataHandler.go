package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sort"
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
	fmt.Println("Storing answers for session", sessionNumber, "and game", gameNumber)
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

func compareSlices(a []string, b string) bool {
	trimmedb := strings.ReplaceAll(b, "[", "")
	trimmedb = strings.ReplaceAll(trimmedb, "]", "")
	bSlice := strings.Split(trimmedb, " ")
	if len(a) != len(bSlice) {
		return false
	}
	sort.Strings(a)
	sort.Strings(bSlice)
	for i, v := range a {
		if v != bSlice[i] {
			return false
		}
	}
	return true
}

func getScore(sessions []Sessions, progression Progression, sessionNumber int, answers []Answer) float64 {
	fmt.Println("Calculating score for session", sessionNumber)
	sectorDetails := loadSessionData(sessions, progression, sessionNumber)
	gameData, err := loadAllGameDate(sectorDetails)
	if err != nil {
		fmt.Println("Error loading game data:", err)
	}

	score := 0.0
	maxExamples := 0
	var selectedGames []Game
	for _, answer := range answers {
		// fmt.Println("Answer:", answer)
		for _, games := range gameData {
			// fmt.Println("Games:", games)
			if games[0].Question[0] == answer.Question {
				// fmt.Println("Found game")
				selectedGames = games
				maxExamples = sectorDetails[0].MaxExamples
				break
			}
		}
	}

	//fmt.Println("Selected games:", selectedGames)
	answerIndex := 0
	for _, game := range selectedGames {
		//fmt.Println("Game:", game)
		//Find the answer for the game
		//fmt.Println("Game: ", game.MainWord.Text, " and answer: ", answers[answerIndex].MainwordText)
		if compareSlices(game.MainWord.Text, answers[answerIndex].MainwordText) {
			//fmt.Println("Found Mainword match")
			//Check if the answer is correct
			fmt.Println("Question: ", game.Question[0], " and submitted question: ", answers[answerIndex].Question)
			fmt.Println("Options: ", game.MainWord.Text, " and submitted options: ", answers[answerIndex].MainwordText)
			fmt.Println("Answer: ", game.Answer[0], " and submitted answer: ", answers[answerIndex].Answer)
			if game.Answer[0] == answers[answerIndex].Answer {
				fmt.Println("Correct answer")
				score++
			}
			answerIndex++
			if answerIndex == len(answers) {
				// Get the MaxExamples for the sector, in order to calculate the score in %
				score = (score / float64(maxExamples)) * 100
				fmt.Println("Score:", score)
				return score
			}
		}
	}
	//decodeAudio(answers[0].Answer)
	return score
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

	var answers []Answer
	err := json.NewDecoder(r.Body).Decode(&answers)
	fmt.Println("Answers:", r.Body)
	if err != nil {
		fmt.Println("Failed to decode JSON body")
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

func loadAllGameDate(sectorDetails []SectorDetails) (map[int][]Game, error) {
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

func loadAnswers(sessionNumber int, gameNumber int) ([]Answer, error) {
	fmt.Println("Loading answers for session", sessionNumber, "and game", gameNumber)
	filename := fmt.Sprintf("%d_%d_answers.json", sessionNumber, gameNumber)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	var answers []Answer
	json.Unmarshal(byteValue, &answers)
	fmt.Println("Answers for game", gameNumber, ":", answers)
	return answers, nil
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
		//Load answers
		for gameNumber := 1; gameNumber <= 5; gameNumber++ {
			answers, err := loadAnswers(sessionNumber-1, gameNumber)
			if err != nil {
				fmt.Println("Error loading answers:", err)
				continue
			}
			fmt.Println()
			//Get score
			score := getScore(sessions, progression, gameNumber, answers)
			fmt.Println("Score for game", gameNumber, ":", score)
		}
		//Set progression per sector
	}
	// Process data (example: print it)
	//fmt.Println("sessions = ", sessions)
	//fmt.Println("progression = ", progression)
	//fmt.Println("sessionNumber = ", sessionNumber)
	//TODO: Progression per sector
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
