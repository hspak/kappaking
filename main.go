package main

import (
	"fmt"
	"log"
)

// TODO: need a way to detect streams going offline for bot and fetching db

// Non-terminal Goroutines:
//    - http server
//    - api call + db update    (5 minutes)
//    - irc channel list update (5 minutes)
//    - irc bot join            (channel blocked)
//    - kappa subtract          (1 minute)
//    - kappa update            (channel blocked)

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatal("Couldn't connect to db:postgres")
	}

	CacheDB.Fresh = false
	CacheDB.Data = nil

	fmt.Println("?")
	streamList := make(chan string, 25)

	updateDB(db, streamList)
	go serveWeb(db)
	launchBot(streamList)
}
