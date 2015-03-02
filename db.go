package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var CacheDB Cache

func queryDB(db *sql.DB) ([]Data, error) {
	if CacheDB.Fresh {
		fmt.Println("UPDATE KPM")
		for i, dat := range CacheDB.Data {
			CacheDB.Data[i].Kappa = KPM[dat.DisplayName]
		}
		fmt.Println("return cached")
		return CacheDB.Data, nil
	}

	rows, err := db.Query(`SELECT * FROM streams ORDER BY viewers DESC LIMIT 25`)
	if err != nil {
		return nil, err
	}

	i := 0
	dat := make([]Data, 25)
	for rows.Next() {
		var name string
		var game string
		var status string
		var url string
		var logo string
		var viewers int
		var kpm int
		var kappa int
		err = rows.Scan(&name, &viewers, &game, &logo, &status, &url, &kpm, &kappa)
		if err != nil {
			return nil, err
		}
		dat[i].DisplayName = name
		dat[i].Game = game
		dat[i].Status = status
		dat[i].Url = url
		dat[i].Logo = logo
		dat[i].Viewers = viewers
		dat[i].Kpm = KPM[name]
		dat[i].Kappa = 9001
		i++
	}
	CacheDB.Fresh = true
	CacheDB.Data = dat
	fmt.Println("return db")
	return dat, nil
}

func addChanList(streamList chan *BotAction) {
	for _, dat := range CacheDB.Data {
		fmt.Println("adding", dat.DisplayName)
		streamList <- &BotAction{Channel: dat.DisplayName, Join: true}
	}
}

func updateDB(db *sql.DB, streamList chan *BotAction) error {
	err := insertDB(db, getTopStreams(true))
	if err != nil {
		log.Fatal(err)
	}

	// wait long enough for the first queryDB
	// so that CacheDB is not empty
	go func() {
		time.Sleep(time.Second * 6)
		addChanList(streamList)
	}()

	// poll twitch every 5 minute
	ticker := time.NewTicker(time.Minute * 5)
	go func() {
		for _ = range ticker.C {
			err = insertDB(db, getTopStreams(false))
			if err != nil {
				log.Fatal(err)
			}
			addChanList(streamList)
			CacheDB.Fresh = false
		}
	}()
	return nil
}

func insertDB(db *sql.DB, streams *Streams) error {
	// no response from Twtich, no update
	if streams == nil {
		return nil
	}

	for _, stream := range streams.Stream {
		streamName := strings.ToLower(stream.Channel.DisplayName)
		tx, err := db.Begin()
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			CREATE TEMPORARY TABLE newvals(
				name 	VARCHAR(255),
				viewers INTEGER,
				game	VARCHAR(255),
				logo	VARCHAR(255),
				status	VARCHAR(255),
				url		VARCHAR(255),
				maxkpm	INTEGER,
				kappa	INTEGER
			);`)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO
				newvals(name, viewers, game, logo, status, url, maxkpm, kappa)
				VALUES($3, $2, $1, $4, $5, $6, $7, $8);`,
			stream.Game, stream.Viewers,
			streamName, stream.Channel.Logo,
			stream.Channel.Status, stream.Channel.Url, KPM[streamName], 9001)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`LOCK TABLE streams IN EXCLUSIVE MODE;`)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			UPDATE streams
			SET
				viewers = newvals.viewers,
				game = newvals.game,
				logo = newvals.logo,
				url = newvals.url,
				status = newvals.status,
				maxkpm = newvals.maxkpm,
				kappa = newvals.kappa
			FROM newvals
			WHERE newvals.name = streams.name;
		`)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO streams
			SELECT
				newvals.name,
				newvals.viewers,
				newvals.game,
				newvals.logo,
				newvals.status,
				newvals.url,
				newvals.maxkpm,
				newvals.kappa
			FROM newvals
			LEFT OUTER JOIN streams ON (streams.name = newvals.name)
			WHERE streams.name IS NULL;
		`)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`DROP TABLE newvals`)
		if err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "user=kappaking dbname=kappaking sslmode=disable")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS streams(
		name 	VARCHAR(255) PRIMARY KEY,
		viewers INTEGER,
		game	VARCHAR(255),
		logo	VARCHAR(255),
		status	VARCHAR(255),
		url 	VARCHAR(255),
		maxkpm 	INTEGER,
		kappa	INTEGER);`)
	if err != nil {
		return nil, err
	}
	return db, nil
}
