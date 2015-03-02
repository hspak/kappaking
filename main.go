package main

import "log"

// TODO: need a way to detect streams going offline for bot and fetching db

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
		log.Fatal("Couldn't connect to db:postgres")
	}

	CacheDB.Fresh = false
	CacheDB.Data = nil
	streamList := make(chan *BotAction, 25)

	updateDB(db, streamList)
	go serveWeb(db)
	launchBot(streamList)
}
