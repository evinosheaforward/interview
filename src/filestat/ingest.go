package filestat

import (
	//"fmt"
	"bufio"
	"hash/fnv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)


type stream struct {
	io.Reader
	lines       chan string
	file        *os.File
	num_parsers int
	db          pgdb
	//db          cassdb
}

func Ingest(fname string) {
	// maybe pass the open file instead
	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	nparsers, _ := strconv.Atoi(os.Getenv("NUM_PARSERS"))
	s := stream{
		lines:       make(chan string),
		file:        f,
		num_parsers: nparsers,
		db:          NewDBConn(),
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
	h := hash(line)
	nc := len(line)
	nt := len(strings.Fields(line))
	go s.db.InsertCounts(h, nc, nt)
	go s.db.InsertKeywords(h, line) // could ln be passed as pointer?
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
