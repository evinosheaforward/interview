package filestat

import (
	"fmt"
	"log"
	"database/sql"
	"os"
	"strings"
	"strconv"

	"github.com/lib/pq"
)

type dbconn interface {
		InsertCounts()
		InsertKeywords()
}

type pgdb struct {
		dbconn
		db *sql.DB
}

//type cassdb struct {
//   dbconn
//   db gocql.session
//}

func NewDBConn() pgdb { //cassdb {
		//return NewCassDB()
		return NewPGDB()
}

func SetupDB() {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	pass := os.Getenv("PASSWORD")
	name := os.Getenv("DBNAME")
	host := os.Getenv("HOST")
	user := os.Getenv("USER")
	db, err := sql.Open("postgres", fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
			port, pass, name, host, user,
	))
	if err != nil {
			log.Fatal(err)
			panic(err)
	}
   defer db.Close()

   _,err = db.Exec("CREATE DATABASE " + name)
   if err != nil {
       panic(err)
   }

   _,err = db.Exec("USE " + name)
   if err != nil {
       panic(err)
   }

   _,err = db.Exec(
		 `CREATE TABLE IF NOT EXISTS LineInfo (line_hash int,
		 													       PRIMARY KEY (line_hash),
		 																 line_len int,
		 																 num_tokens int,
		 																 times_found int);

		 CREATE TABLE IF NOT EXISTS KeywordInfo (keyword string,
		 													          PRIMARY KEY (keyword),
		 																    times_found int)`)
   if err != nil {
       panic(err)
   }

}

func NewPGDB() pgdb {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	session, err := sql.Open("postgres", fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
			os.Getenv("USER"),
			os.Getenv("PASSWORD"),
			os.Getenv("DBNAME"),
			os.Getenv("HOST"),
			port,
	))
	if err != nil {
			log.Println("Couldn't open connection to postgre database")
	}
	return pgdb{ db: session }
}

func (pg pgdb) InsertCounts(h uint32, nc int, nt int) {
		found := -1
		qry := `SELECT line_hash FROM LineInfo WHERE line_hash == $1`
		iter, err := pg.db.Query(qry, h)
		if err != nil {
			log.Println("Couldn't insert user row into DB")
		}
		iter.Scan(&found)
		if found != -1 {
			pg.db.Exec(
				`INSERT INTO LineInfo (line_hash, num_chars, num_tokens) VALUES ($1, $2, $3)`,
				h,
				nc,
				nt,
			)
		}
		_, err2 := pg.db.Exec(
		    `UPDATE KeywordInfo SET count = count + 1 WHERE hash = $1`,
				h,)
		if err2 != nil {
			 log.Println("Bad update on LineInfo count")
		}
}
/*
func (db cassdb) InsertCounts(h uint32, nc int, nt int) {
	out := -1
	iter := s.session.Query(`SELECT line_hash FROM LineInfo WHERE line_hash == ?`, h).Iter()
	iter.Scan(&out)
	if out != -1 {
		s.session.Query(
			`INSERT INTO LineInfo (line_hash, num_chars, num_tokens) VALUES (?, ?, ?)`,
			h,
			nc,
			nt,
		).Exec()
	}
	err := s.session.Query(`UPDATE KeywordInfo SET count = count + 1 WHERE hash = ?`,
		h).Exec()
	if err != nil {
		 fmt.Println("Bad update on LineInfo count")
	}
}
*/

func (pg pgdb) InsertKeywords(h uint32, ln string) {
	var keyword string
	iter, _ := pg.db.Query(`SELECT keyword FROM KeywordInfo`)
	for iter.Next() {
			iter.Scan(&keyword)
			if strings.Contains(ln, keyword) {
				pg.db.Exec(`UPDATE KeywordInfo SET count = count + 1 WHERE keyword = $1`,
								keyword)
			}
	}
}
/*
func (db cassdb) InsertKeywords(h uint32, ln string) {
	var keyword string
	for iter := s.Keywords() {
	  iter.Scan(&keyword)
		if strings.Contains(ln, keyword) {
			s.session.Query(`UPDATE KeywordInfo SET line_hashes = ? + line_hashes WHERE keyword = ?`,
				hash(ln), keyword).Exec()
		}
	}
}
*/

/*
func (db cassdb) Keywords() *gocql.Iter {
	return db.session.Query(`SELECT keyword FROM KeywordInfo`).Iter()
}
*/
