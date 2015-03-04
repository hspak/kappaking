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

func grabKappa(db *sql.DB, name string) (int, int, error) {
	row, err := db.Query(`SELECT maxkpm, kappa FROM streams WHERE name=$1`, name)
	maxkpm := 0
	kappa := 0
	if err != nil {
		fmt.Println("damn it")
		return 0, 0, err
	}
	if row.Next() {
		err = row.Scan(&maxkpm, &kappa)
		if err != nil {
			fmt.Println("fuck")
			return 0, 0, err
		}
	}
	return maxkpm, kappa, nil
}

func queryDB(db *sql.DB) ([]Data, error) {
	if CacheDB.Fresh {
		// log.Println("DB: returning cache")
		for i, dat := range CacheDB.Data {
			CacheDB.Data[i].CurrKpm = KPM[dat.DisplayName]
			CacheDB.Data[i].MaxKpm = MaxKPM[dat.DisplayName]
			CacheDB.Data[i].Kappa = TotalKappa[dat.DisplayName]
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
		ORDER BY viewers DESC`,
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
		dat[i].CurrKpm = currkpm
		dat[i].MaxKpm = maxkpm
		dat[i].Kappa = kappa
		i++
	}
	CacheDB.Fresh = true
	CacheDB.Data = dat
	// log.Println("DB: returning db")
	return dat, nil
}

func deadStream(name string) (bool, int) {
	for i, curr := range LiveStreams {
		if name == curr {
			fmt.Println("not stream", name)
			return false, i
		}
	}
	fmt.Println("dead stream", name)
	return true, -1
}

func addChanList(streamList chan *BotAction) {
	for i := 0; i < 25; i++ {
		if dead, j := deadStream(PrevStreams[i]); dead {
			streamList <- &BotAction{Channel: LiveStreams[j], Join: false}
		} else {
			streamList <- &BotAction{Channel: LiveStreams[i], Join: true}
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
		for _, stream := range LiveStreams {
			streamList <- &BotAction{Channel: stream, Join: true}
		}
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
			stream.Channel.Status, stream.Channel.Url, KPM[streamName], MaxKPM[streamName], TotalKappa[streamName])
		// TODO: maxkpm, kappa values
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
