package main

import "time"

type Stored struct {
	maxKPM     map[string]int
	dateKPM    map[string]time.Time
	totalKappa map[string]int
	minutes    map[string]int
}

type Cache struct {
	Fresh bool
	kpm   map[string]int
	store Stored
}

func NewCache() *Cache {
	return &Cache{
		kpm: make(map[string]int),
		store: Stored{
			maxKPM:     make(map[string]int),
			totalKappa: make(map[string]int),
			minutes:    make(map[string]int),
		},
	}
}

func (c *Cache) UpdateKPM(streamer string, kpm int) {
	c.kpm[streamer] = kpm
}

func (c *Cache) UpdateTotalKappa(streamer string, kappa int) {
	c.store.totalKappa[streamer] = kappa
}

func (c *Cache) UpdateMinutes(streamer string, minutes int) {
	c.store.minutes[streamer] = minutes
}

func (c *Cache) UpdateMaxDate(streamer string, max int, date time.Time) {
	if max > c.store.maxKPM[streamer] {
		c.store.maxKPM[streamer] = max
		c.store.dateKPM[streamer] = date
	}
}
