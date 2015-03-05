package main

import (
	"log"
	"time"
)

// Non-terminal Goroutines:
//    - http server
//    - api call + db update    (5 minutes)
//    - irc channel list update (5 minutes)
//    - irc bot join            (channel blocked)
//    - kappa subtract          (1 minute)
//    - kappa update            (channel blocked)

type BotAction struct {
	Channel string
	Join    bool
}

func main() {
	// defer profile.Start(profile.CPUProfile).Stop()
	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	CacheDB.Fresh = false
	CacheDB.Data = nil
	KPM = make(map[string]int)
	MaxKPM = make(map[string]int)
	DateKPM = make(map[string]time.Time)
	TotalKappa = make(map[string]int)
	Minutes = make(map[string]int)

	streamList := make(chan *BotAction, 25)

	updateDB(db, streamList)
	time.Sleep(time.Second)
	go serveWeb(db)
	launchBot(db, streamList)
}
