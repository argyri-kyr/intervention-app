// main.go
package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
)

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
	Sentence                     string                         `json:"sentence"`
	PhonologicRhymeSentenceItems []PhonologicRhymeSentenceItems `json:"items"`
}

type PhonologicRhymeSentenceItems struct {
	Sentence string   `json:"sentence"`
	Answers  []string `json:"answers"`
	Correct  int      `json:"correct"`
}

type DecodingTextReading struct {
	FirstQuestion  string   `json:"firstquestion"`
	Section        int      `json:"section"`
	Level          int      `json:"level"`
	ExerciseNumber int      `json:"exercise"`
	TextData       string   `json:"textdata"`
	SecondQuestion string   `json:"secondquestion"`
	RedWords       []string `json:"redwords"`
}

type Answer struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type PhonologicExercises struct {
	PhonologicRhymePair          []PhonologicRhymePair
	PhonologicRhymeMatch         []PhonologicRhymeMatch
	PhonologicRhymeMultipleMatch []PhonologicRhymeMultipleMatch
	PhonologicRhymeSentence      []PhonologicRhymeSentence
	DecodingTextReading          []DecodingTextReading
}

var answers []Answer

var phonologicRhymePairExercise []PhonologicRhymePair
var phonologicRhymeMatchExercise []PhonologicRhymeMatch
var phonologicRhymeMultipleMatchExercise []PhonologicRhymeMultipleMatch
var phonologicRhymeSentenceExercise []PhonologicRhymeSentence
var decodingTextReadingExercise []DecodingTextReading

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

func loadSession() (sessions []Session) {
	var data []Session
	readJSONFile("data/Session.json", &data)

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
	readJSONFile("data/DecodingTextReading.json", &decodingTextReadingExercise)

	// Filter data according to Session input
	filtered_phonologicRhymePairExercise := make([]PhonologicRhymePair, 0)
	filtered_phonologicRhymeMatch := make([]PhonologicRhymeMatch, 0)
	filtered_phonologicRhymeMultipleMatch := make([]PhonologicRhymeMultipleMatch, 0)
	filtered_phonologicSentence := make([]PhonologicRhymeSentence, 0)
	filtered_decodingTextReading := make([]DecodingTextReading, 0)

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

		//filter DecodingTextReading exercises
		for _, exercise_dtr := range decodingTextReadingExercise {
			if exercise_dtr.Section == session.Section && exercise_dtr.Level == session.Level && exercise_dtr.ExerciseNumber == session.Exercise {
				// log.Println("Selected Exercise: ", exercise_dtr)
				filtered_decodingTextReading = append(filtered_decodingTextReading, exercise_dtr)
			}
		}

	}

	phonologicRhymePairExercise = filtered_phonologicRhymePairExercise
	phonologicRhymeMatchExercise = filtered_phonologicRhymeMatch
	phonologicRhymeMultipleMatchExercise = filtered_phonologicRhymeMultipleMatch
	phonologicRhymeSentenceExercise = filtered_phonologicSentence
	decodingTextReadingExercise = filtered_decodingTextReading

}

func serveExercises(w http.ResponseWriter, r *http.Request) {
	// Create an instance of ExerciseData with phonologicExercises
	exerciseData := PhonologicExercises{
		PhonologicRhymePair:          phonologicRhymePairExercise,
		PhonologicRhymeMatch:         phonologicRhymeMatchExercise,
		PhonologicRhymeMultipleMatch: phonologicRhymeMultipleMatchExercise,
		PhonologicRhymeSentence:      phonologicRhymeSentenceExercise,
		DecodingTextReading:          decodingTextReadingExercise,
	}

	// Parse the HTML template
	tmpl, err := template.ParseFiles("index.html")
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

func main() {
	sessions := loadSession()

	loadExercises(sessions)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	http.HandleFunc("/exercises", serveExercises)
	http.HandleFunc("/answer", handleAnswer)

	log.Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
