package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Streams struct {
	// {{{
	Stream []struct {
		// Id   int64  `json:"_id"`
		Game    string `json:"game"`
		Viewers int    `json:"viewers"`
		// CreatedAt string `json:"created_at"`
		// Links     struct {
		// Self string `json:"self"`
		// } `json:"_links"`
		// Preview struct {
		// Small    string `json:"small"`
		// Medium   string `json:"medium"`
		// Large    string `json:"large"`
		// Template string `json:"template"`
		// } `json:"preview"`
		Channel struct {
			// Links struct {
			// Self          string `json:"self"`
			// Follows       string `json:"follows"`
			// Commercial    string `json:"commercial"`
			// StreamKey     string `json:"stream_key"`
			// Chat          string `json:"chat"`
			// Features      string `json:"features"`
			// Subscriptions string `json:"subscriptions"`
			// Editors       string `json:"editors"`
			// Videos        string `json:"videos"`
			// Teams         string `json:"teams"`
			// } `json:"_links"`
			// Background                   interface{} `json:"background"`
			// Banner                       interface{} `json:"banner"`
			// BroadcasterLanguage          string      `json:"broadcaster_language"`
			DisplayName string `json:"display_name"`
			// Game        string `json:"game"`
			Logo string `json:"logo"`
			// Mature                       bool        `json:"mature"`
			// Status string `json:"status"`
			// Partner                      bool        `json:"partner"`
			// Url string `json:"url"`
			// VideoBanner                  string      `json:"video_banner"`
			// Id                           int         `json:"_id"`
			// Name                         string      `json:"name"`
			// CreatedAt                    string      `json:"created_at"`
			// UpdatedAt                    string      `json:"updated_at"`
			// Delay                        int         `json:"delay"`
			// Followers                    int         `json:"followers"`
			// ProfileBanner                string      `json:"profile_banner"`
			// ProfileBannerBackgroundColor interface{} `json:"profile_banner_background_color"`
			// Views                        int         `json:"views"`
			// Language                     string      `json:"language"`
		} `json:"channel"`
	} `json:"streams"`
	// Total int `json:"_total"`
	// Links struct {
	// Self     string `json:"self"`
	// Next     string `json:"next"`
	// Featured string `json:"featured"`
	// Summary  string `json:"summary"`
	// Followed string `json:"followed"`
	// } `json:"_links"`
	// }}}
}

type MostKappa struct {
	Name   string `json:"name"`
	Kappas int    `json:"kappas"`
}

type HighestKPM struct {
	Name string    `json:"name"`
	Kpm  int       `json:"kpm"`
	Date time.Time `json:"kpm_date"`
}

type HighestAvg struct {
	Name string  `json:"name"`
	Avg  float64 `json:"avg"`
}

type Leaders struct {
	MostKappa  []MostKappa  `json:"most_kappa"`
	HighestKPM []HighestKPM `json:"highest_kpm"`
	HighestAvg []HighestAvg `json:"highest_avg"`
}

var LiveStreams [25]string
var PrevStreams [25]string

// It is OK if this returns nil. It will only prevent the DB from updating
func getTopStreams(first bool) *Streams {
	res, err := http.Get("https://api.twitch.tv/kraken/streams")
	if err != nil {
		log.Println("info: no response from api")
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("error: could not read api response")
		return nil
	}

	dat := new(Streams)
	if err := json.Unmarshal(body, dat); err != nil {
		log.Println("could not parse json")
		return nil
	}

	if !first {
		for i := 0; i < 25; i++ {
			PrevStreams[i] = LiveStreams[i]
		}
	}
	for i, _ := range dat.Stream {
		dat.Stream[i].Channel.DisplayName = strings.ToLower(dat.Stream[i].Channel.DisplayName)
		LiveStreams[i] = dat.Stream[i].Channel.DisplayName
	}
	log.Println("==> new wave")
	return dat
}

func JSONMarshal(v interface{}, safeEncoding bool) ([]byte, error) {
	b, err := json.MarshalIndent(v, "", "  ")

	if safeEncoding {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
}

func returnJSON(db *DB) string {
	d, err := db.Query()
	if err != nil {
		// TODO: do more
		return err.Error()
	}
	wrapper := &Wrapper{Streams: d}
	out, err := JSONMarshal(&wrapper, true)
	if err != nil {
		// TODO: do more
		return "error"
	}
	return string(out)
}

func returnLeadersJSON(db *DB) string {
	kappa, err := db.queryKappa()
	if err != nil {
		return "query kappa error"
	}
	kpm, err := db.queryHighestKPM()
	if err != nil {
		return "highest kpm error"
	}
	avg, err := db.queryHighestAvg()
	if err != nil {
		return "highest avg error"
	}
	data := &Leaders{MostKappa: kappa, HighestKPM: kpm, HighestAvg: avg}
	out, err := JSONMarshal(data, true)
	if err != nil {
		return "json marshal error"
	}
	return string(out)
}
