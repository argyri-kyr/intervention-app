package main

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
	Question []string `json:"question"`
	MainWord MainWord `json:"mainWord"`
	Options  Options  `json:"options"`
	Answer   []string `json:"answer"`
}

type GameToHtml struct {
	Template string `json:"template"`
	Games    []Game `json:"game"`
}

type Sessions struct {
	Session int      `json:"session"`
	Sectors []string `json:"sectors"`
}

type SectorEntry struct {
	Section     int    `json:"section"`
	Difficulty  int    `json:"difficulty"`
	GameNo      int    `json:"game"`
	MaxExamples int    `json:"maxExamples"`
	Template    string `json:"template"`
}

type Progression struct {
	PhonologicalAwareness []SectorEntry `json:"phonological_awareness"`
	Decoding              []SectorEntry `json:"decoding"`
	ReadingComprehension  []SectorEntry `json:"reading_comprehension"`
	Memory                []SectorEntry `json:"memory"`
	ProcessingSpeed       []SectorEntry `json:"processing_speed"`
}

type SectorDetails struct {
	Sector      string
	Section     int
	Difficulty  int
	GameNo      int
	MaxExamples int
	Template    string
}

type Session struct {
	Section  int `json:"section"`
	Level    int `json:"level"`
	Exercise int `json:"exercise"`
}

type Answer struct {
	Question     string `json:"question"`
	MainwordText string `json:"mainwordtext"`
	Answer       string `json:"answer"`
}
