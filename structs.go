package main

// bot
type KappaData struct {
	Name  string
	Count int
}

// request
type Data struct {
	DisplayName string `json:"display_name"`
	Game        string `json:"game"`
	Viewers     int    `json:"viewers"`
	CurrKpm     int    `json:"currkpm"`
	MaxKpm      int    `json:"maxkpm"`
	Kappa       int    `json:"kappa"`
	Logo        string `json:"logo"`
	Status      string `json:"status"`
	Url         string `json:"url"`
}

type Wrapper struct {
	Streams []Data
}

// db
type Cache struct {
	Fresh bool
	Data  []Data
}