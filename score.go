package main

import (
	"log"
	"sort"
	"strings"
)


func compareSlices(a []string, b string) bool {
	trimmedb := strings.ReplaceAll(b, "[", "")
	trimmedb = strings.ReplaceAll(trimmedb, "]", "")
	bSlice := strings.Split(trimmedb, ",")

	if len(a) != len(bSlice) {
		// log.Println("a:", len(a))
		// log.Println("b:", len(bSlice))
		// log.Println("Lengths do not match")
		return false
	}

	for i, v := range a {
		a[i] = strings.TrimSpace(v)
	}
	sort.Strings(a)

	for i, v := range bSlice {
		bSlice[i] = strings.TrimSpace(v)
	}
	sort.Strings(bSlice)
	// log.Println("a:", a)
	// log.Println("b:", bSlice)
	for i, v := range a {
		if v != bSlice[i] {
			// log.Println("v:", v)
			// log.Println("bSlice[i]:", bSlice[i])
			// log.Println("Mismatch at index", i)
			return false
		}
	}
	return true
}

func getScore(sessions []Sessions, progression Progression, sessionNumber int, gameNumber int, answers []Answer) float64 {
	log.Println("Calculating score for session", sessionNumber)
	sectorDetails := loadSessionData(sessions, progression, sessionNumber)
	gameData, err := loadAllGameData(sectorDetails)
	if err != nil {
		log.Println("Error loading game data:", err)
	}

	score := 0.0
	maxExamples := 0
	var selectedGames []Game
	for _, answer := range answers {
		// log.Println("Answer:", answer)
		for _, games := range gameData {
			// log.Println("Games:", games)
			if games[gameNumber-1].Question[0] == answer.Question {
				// log.Println("Found game")
				selectedGames = games
				// log.Println(sectorDetails[gameNumber-1])
				maxExamples = sectorDetails[gameNumber-1].MaxExamples
				// log.Println("Max examples:", maxExamples, " for game number:", gameNumber)
				break
			}
		}
	}

	for _, answer := range answers {
		for _, game := range selectedGames {
			// log.Println("MainWord:", game.MainWord.Text, " and submitted MainWord:", answer.MainwordText)
			if compareSlices(game.MainWord.Text, answer.MainwordText) {
				//Check if the answer is correct
				// log.Println("Question: ", game.Question[0], " and submitted question: ", answers[answerIndex].Question)
				// log.Println("Options: ", game.MainWord.Text, " and submitted options: ", answers[answerIndex].MainwordText)
				submittedAnswers := strings.Split(answer.Answer, ",")
				correctOptionAnswer := 0
				for s := range submittedAnswers {
					for a := range game.Answer {
						trimmedAnswer := strings.ReplaceAll(submittedAnswers[s], "[", "")
						trimmedAnswer = strings.ReplaceAll(trimmedAnswer, "]", "")
						trimmedAnswer = strings.ReplaceAll(trimmedAnswer, ".png", "")
						// log.Println("Answer: ", game.Answer[a], " and submitted answer: ", trimmedAnswer)
						if game.Answer[a] == trimmedAnswer {
							correctOptionAnswer++
						}
					}
				}
				if correctOptionAnswer == len(game.Answer) {
					// log.Println("Correct answer")
					score++
					continue
				}
			}
		}
	}

	// log.Println("All answers checked")
	// Get the MaxExamples for the sector, in order to calculate the score in %
	score = (score / float64(maxExamples)) * 100
	// log.Println("Score:", score)
	return score
	//decodeAudio(answers[0].Answer)
}