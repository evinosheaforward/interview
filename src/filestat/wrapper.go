package filestat

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

type dbconn interface {
	InsertCounts(uint32, int, int)
	InsertKeywords(uint32, string)
	SelectCounts() lineData
	SelectKeywords() map[string]int
}

type pgdb struct {
	db *sql.DB
}

func (pg pgdb) InsertCounts(h uint32, nc int, nt int) {
	pg.db.Exec(fmt.Sprintf(
		`INSERT INTO LineInfo
						(line_hash, line_len, num_tokens, times_found)
						VALUES (%d, %d, %d, 1)
				ON CONFLICT (line_hash) DO UPDATE
						SET times_found = LineInfo.times_found + 1`,
		h, nc, nt))
}

func (pg pgdb) InsertKeywords(h uint32, ln string) {
	var keyword string
	iter, err := pg.db.Query("SELECT keyword FROM KeywordInfo")
	check(err)
	defer iter.Close()
	for iter.Next() {
		if err := iter.Scan(&keyword); err != nil {
			check(err)
		}
		if strings.Contains(strings.ToLower(ln), keyword) {
			pg.db.Exec(
				"UPDATE KeywordInfo SET times_found = times_found + 1 WHERE keyword = $1",
				keyword)
		}
	}
}

func (pg pgdb) SelectCounts() lineData {
	iter, err := pg.db.Query(
		"SELECT line_len, num_tokens, times_found FROM LineInfo",
	)
	defer iter.Close()
	check(err)
	var data lineData
	var nc int
	var nt int
	var n int
	for iter.Next() {
		if err := iter.Scan(&nc, &nt, &n); err != nil {
			check(err)
		}
		data = append(data, lineInfo{
			numChars:  nc,
			numTokens: nt,
			count:     n})
	}
	return data
}

func (pg pgdb) SelectKeywords() map[string]int {
	iter, err := pg.db.Query("SELECT * FROM KeywordInfo")
	check(err)
	defer iter.Close()
	var keyword string
	var times_found int
	m := make(map[string]int)
	for iter.Next() {
		if err := iter.Scan(&keyword, &times_found); err != nil {
			check(err)
		}
		m[keyword] = times_found
	}
	return m
}

type lineInfo struct {
	numChars  int
	numTokens int
	count     int
}

func (datum lineInfo) NumChars() int {
	return datum.numChars
}

func (datum lineInfo) NumTokens() int {
	return datum.numTokens
}

func (datum lineInfo) Count() int {
	return datum.count
}

type lineData []lineInfo

func (data lineData) CountDupes() int {
	dupeCount := 0
	for _, lineInfo := range data {
		if lineInfo.count > 1 {
			dupeCount += 1
		}
	}
	return dupeCount
}

// Should accept streamer, check, continue
func InsertInfo(s stream, line string) {
	h := hash(line)
	nc := len(line)
	nt := len(strings.Fields(line))
	go s.db.InsertCounts(h, nc, nt)
	go s.db.InsertKeywords(h, line) // could ln be passed as pointer?
}

func SetupDB(conn dbconn) {
	pg, _ := conn.(pgdb)
	_, err := pg.db.Exec(
		`CREATE TABLE IF NOT EXISTS LineInfo (line_hash bigint,
						 													       PRIMARY KEY (line_hash),
						 																 line_len int,
						 																 num_tokens int,
						 																 times_found int);
			 CREATE TABLE IF NOT EXISTS KeywordInfo (keyword TEXT,
						 													         PRIMARY KEY (keyword),
						 																   times_found int)`)
	check(err)
	SetupKeywords(pg.db)
}

func SetupKeywords(db *sql.DB) {
	f, err := os.Open(os.Getenv("KEYWORD_FILE"))
	check(err)
	scanner := bufio.NewScanner(f)
	var line string
	qry := "INSERT INTO KeywordInfo (keyword, times_found) VALUES ($1, 0)"
	for scanner.Scan() {
		line = scanner.Text()
		db.Exec(qry, strings.ToLower(line))
	}
}

func NewDBConn() pgdb { //cassdb {
	//return NewCassDB()
	return NewPGDB()
}

func NewPGDB() pgdb {
	session, err := sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("USER"),
		os.Getenv("PASSWORD"),
		os.Getenv("DBNAME"),
		os.Getenv("HOST"),
		os.Getenv("PORT"),
	))
	check(err)
	return pgdb{db: session}
}
