package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const basePath  = "./data"
const htmlFile = "exercises.html"

func serveData(w http.ResponseWriter, r *http.Request) {
	sessionNumber := 1
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) > 2 {
		var err error
		sessionNumber, err = strconv.Atoi(parts[2])
		if err == nil {
			// Perform actions with sessionNumber
		} else {
			sessionNumber = 1
			// fmt.Println("parts[2] is not an integer")
			// Do nothing or handle the error
		}
	} else {
	}

	gameData, err := LoadData(sessionNumber)
	if err != nil {
		log.Fatalf("Failed to load data: %v", err)
	}
	// for _, g := range gameData {
	// 	fmt.Println(g.Template)
	// }
	// Parse the HTML template
	tmpl, err := template.New(htmlFile).ParseFiles(htmlFile)
	if err != nil {
		log.Fatalf("Failed to parse the HTML template: %v", err)
	}

	// fmt.Println("game Data = ", gameData)
	// length := len(gameData)
	// fmt.Println("Length of gameData:", length)

	// Execute the HTML template with the combined games data
	err = tmpl.Execute(w, gameData)
	if err != nil {
		log.Fatalf("Failed to execute the HTML template: %v", err)
	}
}

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	http.HandleFunc("/exercises/", serveData)
	http.HandleFunc("/answers", HandleAnswer)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	log.Println("Starting server on port 8080")
	err1 := http.ListenAndServe(":8080", nil)
	if err1 != nil {
		log.Fatalf("Failed to start server: %v", err1)
	}
}
