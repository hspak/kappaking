package main

// request
type Data struct {
	DisplayName string `json:"display_name"`
	Game        string `json:"game"`
	Viewers     int    `json:"viewers"`
	Logo        string `json:"logo"`
	CurrKpm     int    `json:"currkpm"`
	MaxKpm      int    `json:"maxkpm"`
	MaxKpmDate  string `json:"maxkpm_date"`
	Kappa       int    `json:"kappa"`
	Minutes     int    `json:"minutes"`
}

type Wrapper struct {
	Streams []Data
}
