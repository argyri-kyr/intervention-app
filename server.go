package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const basePath = "./data"
const htmlFile = "exercises.html"
const serverPort = "8080"

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
	}

	gameData, _, err := LoadData(sessionNumber)
	if err != nil {
		log.Fatalf("Failed to load data: %v", err)
	}
	// for _, g := range gameData {
	// 	fmt.Println(g.Template)
	// }
	// Parse the HTML template
	funcMap := template.FuncMap{
		"join": strings.Join,
	}
	tmpl, err := template.New(htmlFile).Funcs(funcMap).ParseFiles(htmlFile)
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
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Set the output destination of the standard logger to the log file
	log.SetOutput(logFile)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	http.HandleFunc("/exercises/", serveData)
	http.HandleFunc("/answers/", HandleAnswer)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	log.Println("Starting server on port", serverPort)
	fmt.Println("Starting server on port", serverPort)
	fmt.Printf("Access app at http://localhost:%s/exercises/1\n", serverPort)
	err1 := http.ListenAndServe(":"+serverPort, nil)
	if err1 != nil {
		log.Fatalf("Failed to start server: %v", err1)
	}
}
