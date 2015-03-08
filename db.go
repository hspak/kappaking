package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var CacheDB Cache

func grabCounts(db *sql.DB, name string) (int, int, int, time.Time, error) {
	row, err := db.Query(`SELECT maxkpm, kappa, minutes, kpmdate FROM streams WHERE name=$1`, name)
	var maxkpm int
	var kappa int
	var minutes int
	var kpmdate time.Time
	if err != nil {
		return -1, -1, -1, time.Now().UTC(), err
	}
	if row.Next() {
		err = row.Scan(&maxkpm, &kappa, &minutes, &kpmdate)
		if err != nil {
			return -1, -1, -1, time.Now().UTC(), err
		}
	}
	return maxkpm, kappa, minutes, kpmdate, nil
}

func queryDB(db *sql.DB) ([]Data, error) {
	if CacheDB.Fresh {
		for i, dat := range CacheDB.Data {
			CacheDB.Data[i].CurrKpm = KPM[dat.DisplayName]
			CacheDB.Data[i].MaxKpm = MaxKPM[dat.DisplayName]
			CacheDB.Data[i].Kappa = TotalKappa[dat.DisplayName]
			CacheDB.Data[i].Minutes = Minutes[dat.DisplayName]
			CacheDB.Data[i].MaxKpmDate = DateKPM[dat.DisplayName].Format(time.RFC3339)
		}
		return CacheDB.Data, nil
	}

	rows, err := db.Query(`SELECT * FROM streams
		WHERE name IN (
			$1, $2 ,$3 ,$4 ,$5,
			$6, $7 ,$8 ,$9 ,$10,
			$11,$12,$13,$14,$15,
			$16,$17,$18,$19,$20,
			$21,$22,$23,$24,$25)
		ORDER BY kappa DESC`,
		LiveStreams[0], LiveStreams[1], LiveStreams[2], LiveStreams[3], LiveStreams[4],
		LiveStreams[5], LiveStreams[6], LiveStreams[7], LiveStreams[8], LiveStreams[9],
		LiveStreams[10], LiveStreams[11], LiveStreams[12], LiveStreams[13], LiveStreams[14],
		LiveStreams[15], LiveStreams[16], LiveStreams[17], LiveStreams[18], LiveStreams[19],
		LiveStreams[20], LiveStreams[21], LiveStreams[22], LiveStreams[23], LiveStreams[24])
	if err != nil {
		return nil, err
	}

	i := 0
	dat := make([]Data, 25)
	for rows.Next() {
		var name string
		var game string
		var logo string
		var viewers int
		var currkpm int
		var maxkpm int
		var kappa int
		var minutes int
		var date time.Time
		err = rows.Scan(&name, &viewers, &game, &logo,
			&currkpm, &maxkpm, &kappa, &minutes, &date)
		if err != nil {
			return nil, err
		}
		dat[i].DisplayName = name
		dat[i].Game = game
		dat[i].Logo = logo
		dat[i].Viewers = viewers
		dat[i].CurrKpm = currkpm
		dat[i].MaxKpm = maxkpm
		dat[i].Kappa = kappa
		dat[i].Minutes = minutes
		dat[i].MaxKpmDate = date.Format(time.RFC3339)
		i++
	}
	CacheDB.Fresh = true
	CacheDB.Data = dat
	return dat, nil
}

func addChanList(streamList chan *BotAction) {
	// part dead channels
	dead := true
	for _, prev := range PrevStreams {
		for _, curr := range LiveStreams {
			// still alive
			if prev == curr {
				dead = false
				break
			}
		}
		if dead {
			streamList <- &BotAction{Channel: prev, Join: false}
		}
		dead = true
	}
	for _, curr := range LiveStreams {
		streamList <- &BotAction{Channel: curr, Join: true}
	}
}

func updateDB(db *sql.DB, streamList chan *BotAction) error {
	topStreams := getTopStreams(true)
	if topStreams == nil {
		log.Fatal("first response from twitch cannot be nil")
	}
	err := insertDB(db, topStreams, true)
	if err != nil {
		log.Fatal(err)
	}

	// grab db values before first request
	for _, stream := range topStreams.Stream {
		name := stream.Channel.DisplayName
		MaxKPM[name], TotalKappa[name], Minutes[name], DateKPM[name], err = grabCounts(db, name)
		if err != nil {
			// TODO: error handling
			log.Println(err)
		}
	}

	for _, stream := range LiveStreams {
		streamList <- &BotAction{Channel: stream, Join: true}
	}

	// poll twitch every 5 minute
	ticker := time.NewTicker(time.Minute * 5)
	go func() {
		for _ = range ticker.C {
			topStreams = getTopStreams(false)
			if topStreams == nil {
				continue
			}

			err = insertDB(db, topStreams, false)
			if err != nil {
				log.Fatal(err)
			}
			addChanList(streamList)
			CacheDB.Fresh = false
		}
	}()
	return nil
}

func insertDB(db *sql.DB, streams *Streams, first bool) error {
	// no response from Twtich, no update
	if streams == nil {
		return nil
	}

	for _, stream := range streams.Stream {
		streamName := stream.Channel.DisplayName
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
				currkpm	INTEGER,
				maxkpm	INTEGER,
				kappa	INTEGER,
				minutes INTEGER,
				kpmdate TIMESTAMP
			);`)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO
				newvals(name, viewers, game, logo, currkpm, maxkpm, kappa, minutes, kpmdate)
				VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
			streamName, stream.Viewers, stream.Game, stream.Channel.Logo,
			KPM[streamName], MaxKPM[streamName], TotalKappa[streamName],
			Minutes[streamName], DateKPM[streamName])
		if err != nil {
			return err
		}

		_, err = tx.Exec(`LOCK TABLE streams IN EXCLUSIVE MODE;`)
		if err != nil {
			return err
		}

		// don't overwrite kappa values on load
		if first {
			_, err = tx.Exec(`
				UPDATE streams
				SET
					viewers = newvals.viewers,
					game = newvals.game,
					logo = newvals.logo
				FROM newvals
				WHERE newvals.name = streams.name;
			`)
		} else {
			_, err = tx.Exec(`
				UPDATE streams
				SET
					viewers = newvals.viewers,
					game = newvals.game,
					logo = newvals.logo,
					currkpm = newvals.currkpm,
					maxkpm = newvals.maxkpm,
					kappa = newvals.kappa,
					minutes = newvals.minutes,
					kpmdate = newvals.kpmdate
				FROM newvals
				WHERE newvals.name = streams.name;
			`)
		}
		if err != nil {
			return err
		}

		// if it's a new entry, use current time, don't input null
		_, err = tx.Exec(`
			INSERT INTO streams
			SELECT
				newvals.name,
				newvals.viewers,
				newvals.game,
				newvals.logo,
				newvals.currkpm,
				newvals.maxkpm,
				newvals.minutes,
				newvals.kappa,
				$1
			FROM newvals
			LEFT OUTER JOIN streams ON (streams.name = newvals.name)
			WHERE streams.name IS NULL;
		`, time.Now().UTC())
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
		currkpm INTEGER,
		maxkpm 	INTEGER,
		kappa 	INTEGER,
		minutes INTEGER,
		kpmdate	TIMESTAMP);`)
	if err != nil {
		return nil, err
	}
	return db, nil
}
