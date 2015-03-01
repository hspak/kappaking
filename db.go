package main

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func updateDB(db *sql.DB) error {
	insertDB(db, getTopStreams())
	// poll twitch every 5 minute
	ticker := time.NewTicker(time.Minute * 5)
	go func() {
		for _ = range ticker.C {
			insertDB(db, getTopStreams())
		}
	}()
	return nil
}

func insertDB(db *sql.DB, streams *Streams) error {
	for _, stream := range streams.Stream {
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
				kappa	INTEGER
			);`)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO
				newvals(name, viewers, game, logo, status, url, kappa)
				VALUES($3, $2, $1, $4, $5, $6, $7);`,
			stream.Game, stream.Viewers,
			stream.Channel.DisplayName, stream.Channel.Logo,
			stream.Channel.Status, stream.Channel.Url, 1)
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
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS streams (
		name 	VARCHAR(255) PRIMARY KEY,
		viewers INTEGER,
		game	VARCHAR(255),
		logo	VARCHAR(255),
		status	VARCHAR(255),
		url 	VARCHAR(255),
		kappa	INTEGER);`)
	return db, nil
}
