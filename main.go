package main

import "log"

// Non-terminal Goroutines:
//    - http server
//    - api call + db update    (5 minutes)
//    - irc channel list update (5 minutes)
//    - irc bot join            (channel blocked)
//    - kappa subtract          (1 minute)
//    - kappa update            (channel blocked)

func main() {
	db, err := NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.conn.Close()

	db.StartUpdateLoop()
	go launchBot(db)
	serveWeb(db)
}
