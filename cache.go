package main

import "time"

type Stored struct {
	MaxKPM     map[string]int
	DateKPM    map[string]time.Time
	TotalKappa map[string]int
	Minutes    map[string]int
}

type Cache struct {
	Fresh bool
	KPM   map[string]int
	Store Stored
}

func NewCache() *Cache {
	return &Cache{
		KPM: make(map[string]int),
		Store: Stored{
			MaxKPM:     make(map[string]int),
			DateKPM:    make(map[string]time.Time),
			TotalKappa: make(map[string]int),
			Minutes:    make(map[string]int),
		},
	}
}
