package filestat

import (
	"bufio"
	"github.com/gocql/gocql"
	"log"
	"os"
	"regexp"
	"sync"
)

type stream struct {
	lines       chan string
	file        *os.File
	num_parsers int
	session     gocql.Session
	//conn var cql-db connection
}

func Ingest(fname string, conn gocql.Session) {
	// maybe pass the open file instead
	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	//make this configurable
	nparsers := 2
	s := stream{
		lines:       make(chan string),
		file:        f,
		num_parsers: nparsers,
		session:     conn,
	}
	defer close(s.lines)
	var wg sync.WaitGroup
	wg.Add(nparsers*3 + 2)
	go startParsers(s)
	go s.Read()
	wg.Wait()
}

func (s stream) Read() {
	scanner := bufio.NewScanner(s.file)
	for scanner.Scan() {
		s.lines <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func startParsers(s stream) {
	var wg sync.WaitGroup
	wg.Add(s.num_parsers * 3)
	for i := 0; i < s.num_parsers; i++ {
		go parse(s)
	}
	wg.Wait()
}

func parse(s stream) {
	for line := range s.lines {
		s.InsertInfo(line)
	}
}

func (s stream) InsertInfo(line string) {
	//s.connection.
	h := hash(line)
	nc := len(line)
	nt := len(tokenize(line))
	go s.InsertCounts(h, nc, nt)
	go s.InsertKeywords(h, line) // could ln be passed as pointer?
}

func (s stream) InsertCounts(h uint32, nc int, nt int) {
	out := -1
	err := s.session.Query(`SELECT line_hash FROM LineInfo WHERE line_hash == ?`,
		h).Exec().Scan(&out)
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
}

func (s stream) InsertKeywords(h uint32, ln string) {
	for keyword := range s.Keywords() {
		if strings.Contains(ln, keyword) {
			s.session.Query(`UPDATE KeywordInfo SET line_hashes = ? + line_hashes WHERE keyword = ?`,
				hash(ln), keyword).Exec()
		}
	}
}

func (s stream) Keywords() *gocql.Iter {
	return s.session.Query(`SELECT keyword FROM KeywordInfo`).Exec().Inter()
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func tokenize(ln string) []string {
	reg := regexp.MustCompile(`\w+(?:'\w+)?|[^\w\s]`)
	return reg.FindAllStringIndex(ln, -1)
}
