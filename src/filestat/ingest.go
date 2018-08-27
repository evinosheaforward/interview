package filestat

import (
	"bufio"
	"io"
	"os"
)

type stream struct {
	lines       chan string
	file        os.File
	num_parsers int
	//conn var cql-db connection
}

func Ingest(fname string) {
	// maybe pass the open file instead
	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	s := stream{lines: make(chan string),
		file:        f,
		num_parsers: 2,
		// conn: db connection
	}
	defer close(s.lines)
	var wg sync.WaitGroup
	wg.Add(num_parsers*3 + 2)
	go startParsers(s)
	go s.Read(file)
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
	ln := make(chan string)
	defer close(ln)
	var wg sync.WaitGroup
	wg.Add(s.num_parsers * 3)
	for i := range s.num_parsers {
		go parse(s)
	}
	wg.Wait()
}

func parse(s stream) error {
	for line := range s.lines {
		ln <- line
		s.InsertInfo(ln)
	}
}

func (s stream) InsertInfo(ln string) error {
	//s.connection.
	h = hash(ln)
	nc = len(ln)
	nt = len(tokenize(ln))
	go s.InsertCounts(h, nc, nt)
	go s.InsertKeywords(h, ln) // could ln be passed as pointer?
}

func (s stream) InsertCounts(h int, nc int, nt int) {
	var out int
	err := s.session.Query(`SELECT line_hash FROM LineInfo WHERE line_hash == ?`,
		linehash).Exec().Scan(&out)
	if out == nil {
		s.session.Query(
			`INSERT INTO LineInfo (line_hash, num_chars, num_tokens) VALUES (?, ?, ?)`,
			h,
			nc,
			nt,
		).Exec()
	}
	err := s.session.Query(`UPDATE KeywordInfo SET count = count + 1 WHERE hash = ?`,
		hash(ln)).Exec()
}

func (s stream) InsertKeywords(h int, ln string) {
	for kwd := range s.Keywords() {
		if strings.Contains(ln, kwd) {
			s.session.Query(`UPDATE KeywordInfo SET line_hashes = ? + line_hashes WHERE keyword = ?`,
				hash(ln), kwd).Exec()
		}
	}
}

func (s stream) Keywords() *Iter {
	return s.session.Query(`SELECT keyword FROM KeywordInfo`).Exec().Inter()
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
