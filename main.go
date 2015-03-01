package main

import "log"

func main() {
	// go launchBot()
	db, err := openDB()
	if err != nil {
		log.Fatal("Couldn't connect to db:postgres")
	}

	CacheDB.Fresh = false
	CacheDB.Data = nil

	updateDB(db)
	// fmt.Println(returnJSON(db))
	serveWeb(db)
}
