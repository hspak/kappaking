package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type BotAction struct {
	Channel string
	Join    bool
}

type DB struct {
	conn       *sql.DB
	streamList chan *BotAction
	cache      *Cache
	currData   []Data
	ready      bool
}

func NewDB() (*DB, error) {
	db, err := sql.Open("postgres", "user=hsp dbname=kappaking sslmode=disable")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS streams(
			name 	VARCHAR(255) PRIMARY KEY,
			viewers INTEGER,
			game	VARCHAR(255),
			logo	VARCHAR(255),
			maxkpm 	INTEGER,
			kappa 	INTEGER,
			minutes INTEGER,
			kpmdate	TIMESTAMP);`)
	if err != nil {
		return nil, err
	}
	return &DB{
		conn:       db,
		streamList: make(chan *BotAction, 25),
		cache:      NewCache(),
		currData:   make([]Data, 25),
		ready:      false,
	}, nil
}

func (db *DB) StartUpdateLoop() {
	topStreams := getTopStreams(true)
	if topStreams == nil {
		log.Fatal("first response from twitch cannot be nil")
	}
	err := db.Insert(topStreams, true)
	if err != nil {
		log.Fatal(err)
	}
	for _, stream := range topStreams.Stream {
		name := stream.Channel.DisplayName
		db.SetupCache(name)
		if err != nil {
			log.Println(err)
		}
	}
	for _, stream := range LiveStreams {
		db.streamList <- &BotAction{Channel: stream, Join: true}
	}
	db.ready = true
	db.updateLoop()
}

func (db *DB) updateLoop() error {
	// poll twitch every 5 minute
	ticker := time.NewTicker(time.Minute)
	go func() {
		for _ = range ticker.C {
			topStreams := getTopStreams(false)
			if topStreams == nil {
				continue
			}

			err := db.Insert(topStreams, false)
			if err != nil {
				log.Fatal(err)
			}
			db.addChanList()
			db.cache.Fresh = false
		}
	}()
	return nil
}

func (db *DB) SetupCache(name string) error {
	query := `SELECT maxkpm, kappa, minutes, kpmdate FROM streams WHERE name=$1`
	row, err := db.conn.Query(query, name)
	if err != nil {
		return err
	}
	if row.Next() {
		err = row.Scan(&db.cache.Store.MaxKPM,
			&db.cache.Store.TotalKappa,
			&db.cache.Store.Minutes,
			&db.cache.Store.DateKPM)
		if err != nil {
			return err
		}
	}
	fmt.Println("Cache MaxKPM:", db.cache.Store.MaxKPM)
	fmt.Println("Cache TotalKappa:", db.cache.Store.TotalKappa)
	fmt.Println("Cache Minute:", db.cache.Store.Minutes)
	fmt.Println("Cache DateKPM:", db.cache.Store.DateKPM)
	return nil
}

func (db *DB) Query() ([]Data, error) {
	// wait for the first twitch request before allowing queries
	for !db.ready {
		time.Sleep(time.Second)
	}

	if db.cache.Fresh {
		for i, dat := range db.currData {
			db.currData[i].CurrKpm = db.cache.KPM[dat.DisplayName]
			db.currData[i].MaxKpm = db.cache.Store.MaxKPM[dat.DisplayName]
			db.currData[i].Kappa = db.cache.Store.TotalKappa[dat.DisplayName]
			db.currData[i].Minutes = db.cache.Store.Minutes[dat.DisplayName]
			db.currData[i].MaxKpmDate = db.cache.Store.DateKPM[dat.DisplayName].Format(time.RFC3339)
		}
		return db.currData, nil
	}

	query := `SELECT * FROM streams WHERE name IN (
		$1, $2 ,$3 ,$4 ,$5,
		$6, $7 ,$8 ,$9 ,$10,
		$11,$12,$13,$14,$15,
		$16,$17,$18,$19,$20,
		$21,$22,$23,$24,$25) ORDER BY kappa DESC`
	rows, err := db.conn.Query(query,
		LiveStreams[0], LiveStreams[1],
		LiveStreams[2], LiveStreams[3],
		LiveStreams[4], LiveStreams[5],
		LiveStreams[6], LiveStreams[7],
		LiveStreams[8], LiveStreams[9],
		LiveStreams[10], LiveStreams[11],
		LiveStreams[12], LiveStreams[13],
		LiveStreams[14], LiveStreams[15],
		LiveStreams[16], LiveStreams[17],
		LiveStreams[18], LiveStreams[19],
		LiveStreams[20], LiveStreams[21],
		LiveStreams[22], LiveStreams[23], LiveStreams[24])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		var date time.Time
		err = rows.Scan(
			&db.currData[i].DisplayName,
			&db.currData[i].Viewers,
			&db.currData[i].Game,
			&db.currData[i].Logo,
			&db.currData[i].MaxKpm,
			&db.currData[i].Kappa,
			&db.currData[i].Minutes,
			&date)
		if err != nil {
			return nil, err
		}
		db.currData[i].MaxKpmDate = date.Format(time.RFC3339)
		i++
	}
	db.cache.Fresh = true
	return db.currData, nil
}

func (db *DB) addChanList() {
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
			db.streamList <- &BotAction{Channel: prev, Join: false}
		}
		dead = true
	}
	for _, curr := range LiveStreams {
		db.streamList <- &BotAction{Channel: curr, Join: true}
	}
}

func (db *DB) Insert(streams *Streams, first bool) error {
	// no response from Twtich, no update
	if streams == nil {
		return nil
	}

	for _, stream := range streams.Stream {
		streamName := stream.Channel.DisplayName
		tx, err := db.conn.Begin()
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			CREATE TEMPORARY TABLE newvals(
				name 	VARCHAR(255),
				viewers INTEGER,
				game	VARCHAR(255),
				logo	VARCHAR(255),
				maxkpm	INTEGER,
				kappa	INTEGER,
				minutes INTEGER,
				kpmdate TIMESTAMP);`)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO
				newvals(name, viewers, game, logo, maxkpm, kappa, minutes, kpmdate)
				VALUES($1, $2, $3, $4, $5, $6, $7, $8);`,
			streamName, stream.Viewers, stream.Game, stream.Channel.Logo,
			db.cache.Store.MaxKPM[streamName],
			db.cache.Store.TotalKappa[streamName],
			db.cache.Store.Minutes[streamName],
			db.cache.Store.DateKPM[streamName])
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
					game    = newvals.game,
					logo    = newvals.logo
				FROM newvals
				WHERE newvals.name = streams.name;
			`)
		} else {
			_, err = tx.Exec(`
				UPDATE streams
				SET
					viewers = newvals.viewers,
					game    = newvals.game,
					logo    = newvals.logo,
					maxkpm  = newvals.maxkpm,
					kappa   = newvals.kappa,
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
