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
		for i, dat := range CacheDB.Data {
			CacheDB.Data[i].Kappa = KPM[dat.DisplayName]
		}
		return CacheDB.Data, nil
	}

	rows, err := db.Query(`SELECT * FROM streams ORDER BY currkpm, viewers DESC LIMIT 50`)
	if err != nil {
		return nil, err
	}

	i := 0
	dat := make([]Data, 50)
	for rows.Next() {
		var name string
		var game string
		var status string
		var url string
		var logo string
		var viewers int
		var currkpm int
		var maxkpm int
		var kappa int
		err = rows.Scan(&name, &viewers, &game, &logo, &status, &url, &currkpm, &maxkpm, &kappa)
		if err != nil {
			return nil, err
		}
		dat[i].DisplayName = name
		dat[i].Game = game
		dat[i].Status = status
		dat[i].Url = url
		dat[i].Logo = logo
		dat[i].Viewers = viewers
		dat[i].CurrKpm = KPM[name]
		dat[i].MaxKpm = 1234
		dat[i].Kappa = 9001
		i++
	}
	CacheDB.Fresh = true
	CacheDB.Data = dat
	return dat, nil
}

func addChanList(streamList chan *BotAction) {
	for _, dat := range CacheDB.Data {
		if dat.DisplayName == "" {
			continue
		}

		fmt.Println("check", dat.DisplayName)
		if live, exist := LiveStreams[dat.DisplayName]; live && exist {
			streamList <- &BotAction{Channel: dat.DisplayName, Join: true}
		} else {
			streamList <- &BotAction{Channel: dat.DisplayName, Join: false}
		}
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
		fmt.Println(streamName, "count:", KPM[streamName])
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
				currkpm	INTEGER,
				maxkpm	INTEGER,
				kappa	INTEGER
			);`)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO
				newvals(name, viewers, game, logo, status, url, currkpm, maxkpm, kappa)
				VALUES($3, $2, $1, $4, $5, $6, $7, $8, $9);`,
			stream.Game, stream.Viewers,
			streamName, stream.Channel.Logo,
			stream.Channel.Status, stream.Channel.Url, KPM[streamName], 1234, 9001)
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
				currkpm = newvals.currkpm,
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
				newvals.currkpm,
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
		currkpm INTEGER,
		maxkpm 	INTEGER,
		kappa	INTEGER);`)
	if err != nil {
		return nil, err
	}
	return db, nil
}
