package main

import (
	"code.google.com/p/go.net/websocket"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	SQL_NBTILES_IN_SEC = "select count(id) from tiles where dthr >= NOW() - '1 second'::INTERVAL"
)

var (
	flcf     = flag.String("c", "rosmdash.json", "file conf")
	fllisten = flag.String("l", ":8080", "listen server")
)

type Db struct {
	Host     string `json:'host'`
	Port     int    `json:'port'`
	User     string `json:'user'`
	Password string `json:'password'`
	Name     string `json:'name'`
}

type Configuration struct {
	Db Db `json:'db'`
}

func newConfig(path string) (config Configuration, err error) {
	cf, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(cf, &config)
	return
}

func newDb(config Configuration) (db *sql.DB, err error) {
	cnx := fmt.Sprintf("host=%s user=%s password=%s dbname=%s", config.Db.Host, config.Db.User, config.Db.Password, config.Db.Name)
	db, err = sql.Open("postgres", cnx)
	return
}

func wsLastSec(db *sql.DB) func(*websocket.Conn) {
	return func(ws *websocket.Conn) {
		for {
			rows, err := db.Query(SQL_NBTILES_IN_SEC)
			if err != nil {
				log.Fatal(err)
			}

			defer rows.Close()
			var nbtiles int = 0
			for rows.Next() {
				rows.Scan(&nbtiles)
			}
			log.Printf("%d\n", nbtiles)
			msg, _ := json.Marshal(nbtiles)
			ws.Write(msg)
			time.Sleep(1 * time.Second)
		}
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tpl/index.html")
	t.Execute(w, nil)
}

func main() {
	flag.Parse()
	cf, err := newConfig(*flcf)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("using %s@%s\n", cf.Db.Name, cf.Db.Host)
	db, err := newDb(cf)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/ws/lastsec/", websocket.Handler(wsLastSec(db)))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", mainHandler)
	http.ListenAndServe(*fllisten, nil)
}
