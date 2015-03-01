package main

import (
	"fmt"
	"log"
)

func main() {
	// go launchBot()
	db, err := openDB()
	if err != nil {
		log.Fatal("Couldn't connect to db:postgres")
	}
	fmt.Println("run")
	err = insertDB(db, getTopStreams())
	if err != nil {
		log.Fatal(err)
	}
	// updateDB()
	fmt.Println("hm no err")
	// serveWeb()
}
