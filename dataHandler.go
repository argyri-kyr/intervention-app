package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Declare and initialize the global variable
const gamesPerSession = 4

var allSectors = []string{
	"phonological_awareness",
	"decoding",
	"reading_comprehension",
	"processing_speed",
	"memory",
}

func initializesavedsectosvector() map[string]bool {
	var savedSectors = map[string]bool{
		"phonological_awareness": false,
		"decoding":               false,
		"reading_comprehension":  false,
		"processing_speed":       false,
		"memory":                 false,
	}
	return savedSectors
}

func storeAnswers(answers []Answer, sessionNumber int, gameNumber int) {
	log.Println("Storing answers for session", sessionNumber, "and game", gameNumber)
	jsonData, err := json.Marshal(answers)
	if err != nil {
		log.Fatalf("Failed to encode answers to JSON: %v", err)
	}

	// Identify which game i am saving
	var sectorDetails []SectorDetails
	if gameNumber == 1 {
		_, sectorDetails, err = LoadData(sessionNumber)
		if err != nil {
			log.Fatalf("Failed to load data: %v", err)
		}
	}
	sessionfilename := "./data/sessions.json"
	currentSector := getSectorNameByIndex(sessionNumber, gameNumber-1, sessionfilename)
	if currentSector == "" {
		fmt.Println("Sector not found or error occurred.")
		return
	}
	sectorDetail, err := getSectorDetails(currentSector, sectorDetails)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println("sectorDetail = ", sectorDetail)
	current_section := sectorDetail.Section
	current_difficulty := sectorDetail.Difficulty
	current_gameNo := sectorDetail.GameNo

	filename := fmt.Sprintf("answers/%d_%d_%d_answers.json", current_section, current_difficulty, current_gameNo)
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

func loadSessionData(sessions []Sessions, progression Progression, sessionNumber int, pa_gameNumber int, de_gameNumber int, rc_gameNumber int, ps_gameNumber int, me_gameNumber int) []SectorDetails {
	var details []SectorDetails
	for _, session := range sessions {
		if session.Session == sessionNumber {
			sectors := session.Sectors
			for _, sector := range sectors {
				var firstElement SectorEntry
				switch sector {
				case "phonological_awareness":
					if len(progression.PhonologicalAwareness) > 0 {
						firstElement = progression.PhonologicalAwareness[pa_gameNumber-1]
					}
				case "decoding":
					if len(progression.Decoding) > 0 {
						firstElement = progression.Decoding[de_gameNumber-1]
					}
				case "reading_comprehension":
					if len(progression.ReadingComprehension) > 0 {
						firstElement = progression.ReadingComprehension[rc_gameNumber-1]
					}
				case "memory":
					if len(progression.Memory) > 0 {
						firstElement = progression.Memory[me_gameNumber-1]
					}
				case "processing_speed":
					if len(progression.ProcessingSpeed) > 0 {
						firstElement = progression.ProcessingSpeed[ps_gameNumber-1]
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

func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
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

func loadGameData(sectorDetails []SectorDetails, sessionNumber int) (map[int][]Game, error) {
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

		selectedGames, _ := getRandomElements(gameSlice, maxExamples, sector, sessionNumber)
		// fmt.Println("\nsector = ", sector)
		// fmt.Println("selectedIndices = ", selectedIndices)
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
	byteValue, _ := io.ReadAll(file)
	var answers []Answer
	json.Unmarshal(byteValue, &answers)
	//log.Println("Answers for game", gameNumber, ":", answers)
	return answers, nil
}

func loadAllSectorLevels(sessionNumber int) (int, int, int, int, int) {
	var pa_gameNumber, de_gameNumber, rc_gameNumber, ps_gameNumber, me_gameNumber int
	for _, sector := range allSectors {
		switch sector {
		case "phonological_awareness":
			pa_gameNumber = getSectorLevel(sessionNumber, sector)

		case "decoding":
			de_gameNumber = getSectorLevel(sessionNumber, sector)

		case "reading_comprehension":
			rc_gameNumber = getSectorLevel(sessionNumber, sector)

		case "processing_speed":
			ps_gameNumber = getSectorLevel(sessionNumber, sector)

		case "memory":
			me_gameNumber = getSectorLevel(sessionNumber, sector)
		}
	}
	return pa_gameNumber, de_gameNumber, rc_gameNumber, ps_gameNumber, me_gameNumber
}

func LoadData(sessionNumber int) ([]GameToHtml, []SectorDetails, error) {
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

	err = createSessionFolder(sessionNumber)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Folder for session %d created or already exists.\n", sessionNumber)
	}

	savedSectors := initializesavedsectosvector()

	// Get answers, get score and determine next games
	filename1 := fmt.Sprintf("saved_data/session%d/SectorLevel.json", sessionNumber)
	deleteOldFile(filename1)

	// Get sector level until now
	pa_gameNumber, de_gameNumber, rc_gameNumber, ps_gameNumber, me_gameNumber := loadAllSectorLevels(sessionNumber)

	// iterate throught CURRENT sectors
	for gameNumber := 1; gameNumber <= gamesPerSession; gameNumber++ {
		sector := sessions[sessionNumber-1].Sectors[gameNumber-1]

		switch sector {
		case "phonological_awareness":
			pa_gameNumber = getSectorLevel(sessionNumber, sector)
			if sessionNumber == 1 || pa_gameNumber == 0 {
				pa_gameNumber = 1
			}

		case "decoding":
			de_gameNumber = getSectorLevel(sessionNumber, sector)
			if sessionNumber == 1 || de_gameNumber == 0 {
				de_gameNumber = 1
			}

		case "reading_comprehension":
			rc_gameNumber = getSectorLevel(sessionNumber, sector)
			if sessionNumber == 1 || rc_gameNumber == 0 {
				rc_gameNumber = 1
			}

		case "processing_speed":
			ps_gameNumber = getSectorLevel(sessionNumber, sector)
			if sessionNumber == 1 || ps_gameNumber == 0 {
				ps_gameNumber = 1
			}

		case "memory":
			me_gameNumber = getSectorLevel(sessionNumber, sector)
			if sessionNumber == 1 || me_gameNumber == 0 {
				me_gameNumber = 1
			}
		}
	}

	// Calculate score for PREVIOUS session, in order to do so i have to load the games from the PREVIOUS sectors
	if sessionNumber > 1 {
		//Load answers
		for gameNumber := 1; gameNumber <= gamesPerSession; gameNumber++ {
			log.Println("Checking game number ", gameNumber, " for session ", sessionNumber-1)

			/////////////////////////////// load data for PREVIOUS session ///////////////////////////////////

			// iterate throught PREVIOUS sectors
			sector := sessions[sessionNumber-2].Sectors[gameNumber-1]
			log.Println("Previous Sector: ", sector)
			fmt.Println("\nsector = ", sector)
			fmt.Println("gameNumber = ", gameNumber)
			previousTemplate := ""
			switch sector {
			case "phonological_awareness":
				pa_gameNumber = getSectorLevel(sessionNumber, sector)
				previousTemplate = progression.PhonologicalAwareness[pa_gameNumber-1].Template
				// pa_gameNumber++
			case "decoding":
				de_gameNumber = getSectorLevel(sessionNumber, sector)
				previousTemplate = progression.Decoding[de_gameNumber-1].Template
				// de_gameNumber++
			case "reading_comprehension":
				rc_gameNumber = getSectorLevel(sessionNumber, sector)
				previousTemplate = progression.ReadingComprehension[rc_gameNumber-1].Template
				// rc_gameNumber++
			case "processing_speed":
				ps_gameNumber = getSectorLevel(sessionNumber, sector)
				previousTemplate = progression.ProcessingSpeed[ps_gameNumber-1].Template
				// ps_gameNumber++
			case "memory":
				me_gameNumber = getSectorLevel(sessionNumber, sector)
				previousTemplate = progression.Memory[me_gameNumber].Template
				// me_gameNumber++
			}

			// log.Println("Progression: ", progression[sessions[sessionNumber-1].Sectors[gameNumber]])
			// log.Println("Has progression: ", progression.)
			answers, err := loadAnswers(sessionNumber-1, gameNumber)
			if previousTemplate == "AudioRecord" {
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

			// Get score
			score := getScore(sessions, progression, sessionNumber-1, gameNumber, answers, pa_gameNumber, de_gameNumber, rc_gameNumber, ps_gameNumber, me_gameNumber)
			log.Println("Score for game", gameNumber, ":", score)

			/////////////////////////////// save all SECTORS each time in order not to lose the sequence ///////////////////////////////////

			// at first i will save the sectors that exist in the previous session, in order not to lose to increase the level of a sector if this sector not exist in the current session
			filename1 := fmt.Sprintf("saved_data/session%d/SectorLevel.json", sessionNumber)
			// deleteOldFile(filename1)
			switch sector {
			case "phonological_awareness":
				if score > 80 {
					pa_gameNumber++
				}
				saveSessionSectorLevelToFile(savedSectors, filename1, sector, pa_gameNumber, score)

			case "decoding":
				if score > 80 {
					de_gameNumber++
				}
				saveSessionSectorLevelToFile(savedSectors, filename1, sector, de_gameNumber, score)

			case "reading_comprehension":
				if score > 80 {
					rc_gameNumber++
				}
				saveSessionSectorLevelToFile(savedSectors, filename1, sector, rc_gameNumber, score)

			case "processing_speed":
				if score > 80 {
					ps_gameNumber++
				}
				saveSessionSectorLevelToFile(savedSectors, filename1, sector, ps_gameNumber, score)

			case "memory":
				if score > 80 {
					me_gameNumber++
				}
				saveSessionSectorLevelToFile(savedSectors, filename1, sector, me_gameNumber, score)
			}

		}

		// save the rest of the sectors
		checkAndSaveMissingSectors(savedSectors, sessionNumber, pa_gameNumber, de_gameNumber, rc_gameNumber, ps_gameNumber, me_gameNumber)

	}

	// sif sessionNumber==1 i will need to save all sectors as none will have saved
	if sessionNumber == 1 {
		filename1 := fmt.Sprintf("saved_data/session%d/SectorLevel.json", sessionNumber)
		deleteOldFile(filename1)
		checkAndSaveMissingSectors(savedSectors, sessionNumber, pa_gameNumber, de_gameNumber, rc_gameNumber, ps_gameNumber, me_gameNumber)
	}

	// Process data (example: print it)
	//log.Println("sessions = ", sessions)
	//log.Println("progression = ", progression)
	//log.Println("sessionNumber = ", sessionNumber)
	//TODO: Progression per sector
	sectorDetails := loadSessionData(sessions, progression, sessionNumber, pa_gameNumber, de_gameNumber, rc_gameNumber, ps_gameNumber, me_gameNumber)
	//log.Println("sectorDetails = ", sectorDetails)

	// Load game data based on sector details
	gameData, err := loadGameData(sectorDetails, sessionNumber)
	//log.Println("gameData = ", gameData)
	for i := 0; i < len(gameData); i++ {
		gameDataToHtml = append(gameDataToHtml, GameToHtml{Template: sectorDetails[i].Template, Games: gameData[i+1]})
	}

	if err != nil {
		log.Fatalf("error loading game data: %v", err)
	}

	return gameDataToHtml, sectorDetails, nil
}
