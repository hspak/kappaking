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

// 1. Pulls initial streamer list from twitch
// 2. Fill cache if previous data exists
// 3. Inserts current data into database
// 4. Add channels to join for the bot
// 5. Start loop
func (db *DB) StartUpdateLoop() {
	topStreams := getTopStreams(true)
	if topStreams == nil {
		log.Fatal("first response from twitch cannot be nil")
	}
	for _, stream := range topStreams.Stream {
		err := db.SetupCache(stream.Channel.DisplayName)
		if err != nil {
			log.Println(err)
		}
	}
	err := db.Insert(topStreams, true)
	if err != nil {
		log.Fatal(err)
	}
	for _, stream := range LiveStreams {
		db.streamList <- &BotAction{Channel: stream, Join: true}
	}
	db.ready = true
	db.updateLoop()
}

// Loop
//   1. Grabs new list of streamers
//   2. Inserts current cache into DB
//   3. Marks current cache dirty
//   4. Marks channels to join/part
func (db *DB) updateLoop() error {
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
			db.cache.Fresh = false
			db.addChanList()
		}
	}()
	return nil
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

func (db *DB) SetupCache(name string) error {
	query := `SELECT maxkpm, kappa, minutes, kpmdate FROM streams WHERE name=$1`
	row, err := db.conn.Query(query, name)
	if err != nil {
		return err
	}
	var max int
	var tot int
	var min int
	var date time.Time
	if row.Next() {
		err = row.Scan(&max, &tot, &min, &date)
		if err != nil {
			return err
		}
		db.cache.Store.MaxKPM[name] = max
		db.cache.Store.TotalKappa[name] = tot
		db.cache.Store.Minutes[name] = min
		db.cache.Store.DateKPM[name] = date
	}
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

		// The new cache needs to be updated
		db.cache.Store.DateKPM[db.currData[i].DisplayName] = date
		db.cache.Store.MaxKPM[db.currData[i].DisplayName] = db.currData[i].MaxKpm
		db.cache.Store.Minutes[db.currData[i].DisplayName] = db.currData[i].Minutes
		db.cache.Store.TotalKappa[db.currData[i].DisplayName] = db.currData[i].Kappa
		i++
	}
	db.cache.Fresh = true
	return db.currData, nil
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
				name 	 VARCHAR(255),
				viewers  INTEGER,
				game	 VARCHAR(255),
				logo	 VARCHAR(255),
				maxkpmt  INTEGER,
				kappat   INTEGER,
				minutest INTEGER,
				kpmdatet TIMESTAMP);`)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO
				newvals(name, viewers, game, logo, maxkpmt, kappat, minutest, kpmdatet)
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
				FROM newvals
				WHERE newvals.name = streams.name;
			`)
			_, err = tx.Exec(`
				UPDATE streams
				SET minutes = newvals.minutest
				FROM newvals
				WHERE newvals.name = streams.name AND newvals.minutest > minutes;
			`)
			_, err = tx.Exec(`
				UPDATE streams
				SET maxkpm = newvals.maxkpmt
				FROM newvals
				WHERE newvals.name = streams.name AND newvals.maxkpmt > maxkpm;
			`)
			_, err = tx.Exec(`
				UPDATE streams
				SET kappa = newvals.kappat
				FROM newvals
				WHERE newvals.name = streams.name AND newvals.kappat > kappa;
			`)
			_, err = tx.Exec(`
				UPDATE streams
				SET kappa = newvals.kappat
				FROM newvals
				WHERE newvals.name = streams.name AND newvals.kappat > kappa;
			`)
			_, err = tx.Exec(`
				UPDATE streams
				SET kpmdate = newvals.kpmdatet
				FROM newvals
				WHERE newvals.name = streams.name AND newvals.kpmdatet > kpmdate;
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
				newvals.maxkpmt,
				newvals.minutest,
				newvals.kappat,
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

// {{{ Leaderboards
func (db *DB) queryKappa() ([]MostKappa, error) {
	query := `SELECT name, kappa FROM streams ORDER BY kappa DESC LIMIT 20;`
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	kappas := make([]MostKappa, 20)
	i := 0
	for rows.Next() {
		err := rows.Scan(&kappas[i].Name, &kappas[i].Kappas)
		i++
		if err != nil {
			return nil, err
		}
	}
	return kappas, nil
}

func (db *DB) queryHighestKPM() ([]HighestKPM, error) {
	query := `SELECT name, maxkpm, kpmdate FROM streams ORDER BY maxkpm DESC LIMIT 20;`
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	maxkpm := make([]HighestKPM, 20)
	i := 0
	for rows.Next() {
		err := rows.Scan(&maxkpm[i].Name, &maxkpm[i].Kpm, &maxkpm[i].Date)
		i++
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}
	return maxkpm, nil
}

func (db *DB) queryHighestAvg() ([]HighestAvg, error) {
	query := `SELECT name, kappa, minutes FROM streams ORDER BY COALESCE(kappa/NULLIF(minutes,0),0) DESC LIMIT 20;`
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	avg := make([]HighestAvg, 20)
	kappa := make([]int, 20)
	i := 0
	for rows.Next() {
		err := rows.Scan(&avg[i].Name, &kappa[i], &avg[i].Min)
		avg[i].Avg = float64(kappa[i]) / float64(avg[i].Min)
		i++
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}
	return avg, nil
}

// }}}
