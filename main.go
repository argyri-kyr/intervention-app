// main.go
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type MainWord struct {
	Text  []string `json:"text"`
	Image bool     `json:"image"`
	Audio bool     `json:"audio"`
}

type Options struct {
	Text  []string `json:"text"`
	Image bool     `json:"image"`
	Audio bool     `json:"audio"`
}

type Game struct {
	Question string   `json:"question"`
	MainWord MainWord `json:"mainWord"`
	Options  Options  `json:"options"`
	Answer   []string `json:"answer"`
}

type Session struct {
	Section  int `json:"section"`
	Level    int `json:"level"`
	Exercise int `json:"exercise"`
}

type PhonologicRhymePair struct {
	Question                 string                     `json:"question"`
	Answer                   []string                   `json:"answer"`
	Section                  int                        `json:"section"`
	Level                    int                        `json:"level"`
	ExerciseNumber           int                        `json:"exercise"`
	PhonologicRhymePairItems []PhonologicRhymePairItems `json:"items"`
}

type PhonologicRhymePairItems struct {
	Words   []string `json:"words"`
	Correct int      `json:"correct"`
}

type PhonologicRhymeMatch struct {
	Question                  string                      `json:"question"`
	Section                   int                         `json:"section"`
	Level                     int                         `json:"level"`
	ExerciseNumber            int                         `json:"exercise"`
	PhonologicRhymeMatchItems []PhonologicRhymeMatchItems `json:"items"`
}

type PhonologicRhymeMatchItems struct {
	Word    string   `json:"word"`
	Answers []string `json:"answers"`
	Correct int      `json:"correct"`
}

type PhonologicRhymeMultipleMatch struct {
	Question       string         `json:"question"`
	Section        int            `json:"section"`
	Level          int            `json:"level"`
	ExerciseNumber int            `json:"exercise"`
	Column1        []string       `json:"column1"`
	Column2        []string       `json:"column2"`
	Answers        map[string]int `json:"answers"`
}

type PhonologicRhymeSentence struct {
	Question                     string                         `json:"question"`
	Section                      int                            `json:"section"`
	Level                        int                            `json:"level"`
	ExerciseNumber               int                            `json:"exercise"`
	PhonologicRhymeSentenceItems []PhonologicRhymeSentenceItems `json:"items"`
}

type PhonologicRhymeSentenceItems struct {
	Sentence string   `json:"sentence"`
	Answers  []string `json:"answers"`
	Correct  int      `json:"correct"`
}

type PhonologicRhymeProductionSentence struct {
	Question                               string                                   `json:"question"`
	Section                                int                                      `json:"section"`
	Level                                  int                                      `json:"level"`
	ExerciseNumber                         int                                      `json:"exercise"`
	PhonologicRhymeProductionSentenceItems []PhonologicRhymeProductionSentenceItems `json:"items"`
}

type PhonologicRhymeProductionSentenceItems struct {
	Sentence string `json:"sentence"`
	Word     string `json:"word"`
}

type DecodingTextReading struct {
	Section                  int                        `json:"section"`
	Level                    int                        `json:"level"`
	ExerciseNumber           int                        `json:"exercise"`
	DecodingTextReadingItems []DecodingTextReadingItems `json:"items"`
}

type DecodingTextReadingItems struct {
	FirstQuestion  string   `json:"firstquestion"`
	TextData       string   `json:"textdata"`
	SecondQuestion string   `json:"secondquestion"`
	RedWords       []string `json:"redwords"`
}

type MemoryVisualSingleChoice struct {
	Section                       int                             `json:"section"`
	Level                         int                             `json:"level"`
	ExerciseNumber                int                             `json:"exercise"`
	MemoryVisualSingleChoiceItems []MemoryVisualSingleChoiceItems `json:"items"`
}

type MemoryVisualSingleChoiceItems struct {
	FirstQuestion  string   `json:"firstquestion"`
	Choice1        []string `json:"choice1"`
	SecondQuestion string   `json:"secondquestion"`
	Choice2        []string `json:"choice2"`
	Correct        []int    `json:"correct"`
}

type Answer struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type PhonologicExercises struct {
	PhonologicRhymePair               []PhonologicRhymePair
	PhonologicRhymeMatch              []PhonologicRhymeMatch
	PhonologicRhymeMultipleMatch      []PhonologicRhymeMultipleMatch
	PhonologicRhymeSentence           []PhonologicRhymeSentence
	PhonologicRhymeProductionSentence []PhonologicRhymeProductionSentence
	DecodingTextReading               []DecodingTextReading
	MemoryVisualSingleChoice          []MemoryVisualSingleChoice
}

var answers []Answer

var phonologicRhymePairExercise []PhonologicRhymePair
var phonologicRhymeMatchExercise []PhonologicRhymeMatch
var phonologicRhymeMultipleMatchExercise []PhonologicRhymeMultipleMatch
var phonologicRhymeSentenceExercise []PhonologicRhymeSentence
var phonologicRhymeProductionSentenceExercise []PhonologicRhymeProductionSentence
var decodingTextReadingExercise []DecodingTextReading
var memoryVisualSingleChoice []MemoryVisualSingleChoice

var (
	sectors      = []string{"phonological_awareness", "phonological_awareness"}
	sections     = []int{1, 1}
	difficulties = []int{1, 1}
	games        = []int{1, 2}
	basePath     = "./data"
	gameData     = make(map[int][]Game)
)

func readJSONFile(filePath string, v interface{}) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	err = json.Unmarshal(file, v)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %s", err)
	}
}

func loadSession(sessionNumber int) (sessions []Session) {
	if sessionNumber < 1 {
		return
	}
	var data []Session
	fileName := fmt.Sprintf("data/Session%d.json", sessionNumber)
	readJSONFile(fileName, &data)

	for _, e := range data {
		session := Session{
			Section:  e.Section,
			Level:    e.Level,
			Exercise: e.Exercise,
		}
		sessions = append(sessions, session)
	}

	return sessions
}

func loadExercises(sessions []Session) {
	// Load aLL
	readJSONFile("data/PhonologicRhymePair.json", &phonologicRhymePairExercise)
	readJSONFile("data/PhonologicRhymeMatch.json", &phonologicRhymeMatchExercise)
	readJSONFile("data/PhonologicRhymeMultipleMatch.json", &phonologicRhymeMultipleMatchExercise)
	readJSONFile("data/PhonologicRhymeSentence.json", &phonologicRhymeSentenceExercise)
	readJSONFile("data/PhonologicRhymeProductionSentence.json", &phonologicRhymeProductionSentenceExercise)
	readJSONFile("data/DecodingTextReading.json", &decodingTextReadingExercise)
	readJSONFile("data/MemoryVisualSingleChoice.json", &memoryVisualSingleChoice)

	// Filter data according to Session input
	filtered_phonologicRhymePairExercise := make([]PhonologicRhymePair, 0)
	filtered_phonologicRhymeMatch := make([]PhonologicRhymeMatch, 0)
	filtered_phonologicRhymeMultipleMatch := make([]PhonologicRhymeMultipleMatch, 0)
	filtered_phonologicSentence := make([]PhonologicRhymeSentence, 0)
	filtered_phonologicRhymeProductionSentence := make([]PhonologicRhymeProductionSentence, 0)
	filtered_decodingTextReading := make([]DecodingTextReading, 0)
	filtered_memoryVisualSingleChoice := make([]MemoryVisualSingleChoice, 0)

	for _, session := range sessions {
		// filter PhonologicRhymePair exercises
		for _, exercise_rp := range phonologicRhymePairExercise {
			if exercise_rp.Section == session.Section && exercise_rp.Level == session.Level && exercise_rp.ExerciseNumber == session.Exercise {
				// log.Println("Selected Exercise: ", exercise_rp)
				filtered_phonologicRhymePairExercise = append(filtered_phonologicRhymePairExercise, exercise_rp)
			}
		}

		// filter PhonologicRhymeMatch exercises
		for _, exercise_rm := range phonologicRhymeMatchExercise {
			if exercise_rm.Section == session.Section && exercise_rm.Level == session.Level && exercise_rm.ExerciseNumber == session.Exercise {
				// log.Println("Selected Exercise: ", exercise_rm)
				filtered_phonologicRhymeMatch = append(filtered_phonologicRhymeMatch, exercise_rm)
			}
		}

		//filter PhonologicRhymeMultipleMatch exercises
		for _, exercise_rmm := range phonologicRhymeMultipleMatchExercise {
			if exercise_rmm.Section == session.Section && exercise_rmm.Level == session.Level && exercise_rmm.ExerciseNumber == session.Exercise {
				// log.Println("Selected Exercise: ", exercise_rmm)
				filtered_phonologicRhymeMultipleMatch = append(filtered_phonologicRhymeMultipleMatch, exercise_rmm)
			}
		}

		//filter PhonologicRhymeSentence exercises
		for _, exercise_rs := range phonologicRhymeSentenceExercise {
			if exercise_rs.Section == session.Section && exercise_rs.Level == session.Level && exercise_rs.ExerciseNumber == session.Exercise {
				//log.Println("Selected Exercise: ", exercise_rs)
				filtered_phonologicSentence = append(filtered_phonologicSentence, exercise_rs)
			}
		}

		//filter PhonologicRhymeProductionSentence exercises
		for _, exercise_rps := range phonologicRhymeProductionSentenceExercise {
			if exercise_rps.Section == session.Section && exercise_rps.Level == session.Level && exercise_rps.ExerciseNumber == session.Exercise {
				//log.Println("Selected Exercise: ", exercise_rps)
				filtered_phonologicRhymeProductionSentence = append(filtered_phonologicRhymeProductionSentence, exercise_rps)
			}
		}

		//filter DecodingTextReading exercises
		for _, exercise_dtr := range decodingTextReadingExercise {
			if exercise_dtr.Section == session.Section && exercise_dtr.Level == session.Level && exercise_dtr.ExerciseNumber == session.Exercise {
				// log.Println("Selected Exercise: ", exercise_dtr)
				filtered_decodingTextReading = append(filtered_decodingTextReading, exercise_dtr)
			}
		}

		//filter MemoryVisualSingleChoice exercises
		for _, exercise_mvsc := range memoryVisualSingleChoice {
			if exercise_mvsc.Section == session.Section && exercise_mvsc.Level == session.Level && exercise_mvsc.ExerciseNumber == session.Exercise {
				// log.Println("Selected Exercise: ", exercise_mvsc)
				filtered_memoryVisualSingleChoice = append(filtered_memoryVisualSingleChoice, exercise_mvsc)
			}
		}
	}

	phonologicRhymePairExercise = filtered_phonologicRhymePairExercise
	phonologicRhymeMatchExercise = filtered_phonologicRhymeMatch
	phonologicRhymeMultipleMatchExercise = filtered_phonologicRhymeMultipleMatch
	phonologicRhymeSentenceExercise = filtered_phonologicSentence
	phonologicRhymeProductionSentenceExercise = filtered_phonologicRhymeProductionSentence
	decodingTextReadingExercise = filtered_decodingTextReading
	memoryVisualSingleChoice = filtered_memoryVisualSingleChoice
}

func serveExercises(w http.ResponseWriter, r *http.Request) {
	sessionNumber := 1
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) > 2 {
		sessionNumber, _ = strconv.Atoi(parts[2])
	} else {
	}

	sessions := loadSession(sessionNumber)
	loadExercises(sessions)

	// Create an instance of ExerciseData with phonologicExercises
	exerciseData := PhonologicExercises{
		PhonologicRhymePair:               phonologicRhymePairExercise,
		PhonologicRhymeMatch:              phonologicRhymeMatchExercise,
		PhonologicRhymeMultipleMatch:      phonologicRhymeMultipleMatchExercise,
		PhonologicRhymeSentence:           phonologicRhymeSentenceExercise,
		PhonologicRhymeProductionSentence: phonologicRhymeProductionSentenceExercise,
		DecodingTextReading:               decodingTextReadingExercise,
		MemoryVisualSingleChoice:          memoryVisualSingleChoice,
	}

	// Parse the HTML template
	tmpl, err := template.New("index.html").Funcs(funcMap).ParseFiles("index.html")
	if err != nil {
		log.Fatalf("Failed to parse the HTML template: %v", err)
	}

	// Execute the HTML template with the ExerciseData instance
	err = tmpl.Execute(w, exerciseData)
	if err != nil {
		log.Fatalf("Failed to execute the HTML template: %v", err)
	}
}

func handleAnswer(w http.ResponseWriter, r *http.Request) {
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

// Define a custom function for modulus operation
var funcMap = template.FuncMap{
	"mod": func(a, b int) int {
		return a % b
	},
}

func loadData() error {
	//TODO: load sectors from file sessions.json and get section,difficulty and game from progression.json
	for i := range sections {
		sector := sectors[i]
		section := sections[i]
		difficulty := difficulties[i]
		game := games[i]

		filePath := filepath.Join(basePath, sector, fmt.Sprintf("%d_%d_%d.json", section, difficulty, game))

		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		var gameSlice []Game
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&gameSlice)
		if err != nil {
			return fmt.Errorf("failed to decode JSON data: %w", err)
		}

		// Add games to the appropriate array
		gameData[game] = append(gameData[game], gameSlice...)
	}

	return nil
}

func serveData(w http.ResponseWriter, r *http.Request) {
	loadData()
	// Parse the HTML template
	tmpl, err := template.New("new_index.html").ParseFiles("new_index.html")
	if err != nil {
		log.Fatalf("Failed to parse the HTML template: %v", err)
	}

	// Execute the HTML template with the combined games data
	err = tmpl.Execute(w, gameData)
	if err != nil {
		log.Fatalf("Failed to execute the HTML template: %v", err)
	}
}

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	//http.HandleFunc("/exercises/", serveExercises)
	http.HandleFunc("/exercises/", serveData)
	http.HandleFunc("/answer", handleAnswer)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "main.html")
	})

	err2 := loadData()
	if err2 != nil {
		log.Fatalf("Failed to load data: %v", err2)
	}
	fmt.Println("\n\ngame Data = ", gameData)

	log.Println("Starting server on port 8080")
	err1 := http.ListenAndServe(":8080", nil)
	if err1 != nil {
		log.Fatalf("Failed to start server: %v", err1)
	}

}
