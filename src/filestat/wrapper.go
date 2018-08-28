package filestat

import (
	"bufio"
	"fmt"
	"log"
	"database/sql"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

type dbconn interface {
		InsertCounts()
		InsertKeywords()
}

type pgdb struct {
		dbconn
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
	if err != nil {
			log.Fatal(err)
			panic(err)
	}
	for iter.Next() {
			iter.Scan(&keyword)
			if strings.Contains(strings.ToLower(ln), keyword) {
				pg.db.Exec(
					"UPDATE KeywordInfo SET times_found = times_found + 1 WHERE keyword = $1",
					keyword,)
			}
	}
}
/* struct for cass db that is a dbconn
type cassdb struct {
  dbconn
  db gocql.session
}

func (db cassdb) InsertCounts(h uint32, nc int, nt int) {
	out := -1
	iter := s.session.Query("SELECT line_hash FROM LineInfo WHERE line_hash == ?", h).Iter()
	iter.Scan(&out)
	if out != -1 {
		s.session.Query(
		"INSERT INTO LineInfo (line_hash, num_chars, num_tokens) VALUES (?, ?, ?)",
		h,
		nc,
		nt,
		).Exec()
	}
	err := s.session.Query(
			"UPDATE KeywordInfo SET count = count + 1 WHERE hash = ?",
			h).Exec()
	if err != nil {
		fmt.Println("Bad update on LineInfo count")
	}
}

func (db cassdb) InsertKeywords(h uint32, ln string) {
	var keyword string
	keywords := db.session.Query("SELECT keyword FROM KeywordInfo")
	for iter := keywords().Iter() {
	  iter.Scan(&keyword)
		if strings.Contains(ln, keyword) {
			s.session.Query("UPDATE KeywordInfo SET line_hashes = ? + line_hashes WHERE keyword = ?",
				hash(ln), keyword).Exec()
		}
	}
}
*/
// Should accept streamer, check, continue
func InsertInfo(s stream, line string) {
	h := hash(line)
	nc := len(line)
	nt := len(strings.Fields(line))
	go s.db.InsertCounts(h, nc, nt)
	go s.db.InsertKeywords(h, line) // could ln be passed as pointer?
}

func SetupDB() {
	name := os.Getenv("DBNAME")
	db, err := sql.Open("postgres", fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
			os.Getenv("USER"),
		  os.Getenv("PASSWORD"),
		  name,
		  os.Getenv("HOST"),
		  os.Getenv("PORT"),
	))
	if err != nil {
			log.Fatal(err)
			panic(err)
	}
   defer db.Close()
   _, err = db.Exec(
		 `CREATE TABLE IF NOT EXISTS LineInfo (line_hash bigint,
					 													       PRIMARY KEY (line_hash),
					 																 line_len int,
					 																 num_tokens int,
					 																 times_found int);
		 CREATE TABLE IF NOT EXISTS KeywordInfo (keyword TEXT,
					 													          PRIMARY KEY (keyword),
					 																    times_found int)`)
   if err != nil {
       panic(err)
   }
	 SetupKeywords(db)
}

func SetupKeywords(db *sql.DB) {
		f, err := os.Open(os.Getenv("KEYWORD_FILE"))
		if err != nil {
			log.Fatal(err)
		}
		scanner := bufio.NewScanner(f)
		var line string
		qry := "INSERT INTO KeywordInfo (keyword, times_found) VALUES ($1, $2)"
		for scanner.Scan() {
				line = scanner.Text()
				db.Exec(qry, strings.ToLower(line), 0)
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
	if err != nil {
			log.Println("Couldn't open connection to postgre database")
	}
	return pgdb{ db: session }
}
