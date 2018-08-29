package filestat

import (
	//"fmt"
	"bufio"
	"hash/fnv"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
)

type streamer interface {
	Stream(io.Reader, *sync.WaitGroup)
}

type stream struct {
	lines chan string
	db    pgdb
}

func (s stream) Stream(file io.Reader, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s.lines <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		check(err)
	}
}

func Ingest(fname string, conn dbconn) { // (s streamer)
	db, _ := conn.(pgdb)
	f, err := os.Open(fname)
	check(err)
	defer f.Close()
	s := stream{
		lines: make(chan string),
		db:    db,
	}
	defer close(s.lines)
	nparsers, _ := strconv.Atoi(os.Getenv("NUM_PARSERS"))
	log.Println("Starting ingest subroutines.")
	defer log.Println("Ingesting finished.")
	var wg sync.WaitGroup
	wg.Add(2)
	go startParsers(s, &wg, nparsers)
	go s.Stream(f, &wg)
	wg.Wait()
}

func startParsers(s stream, wg *sync.WaitGroup, num_parsers int) {
	defer wg.Done()
	var wg2 sync.WaitGroup
	wg2.Add(num_parsers)
	for i := 0; i < num_parsers; i++ {
		go parse(s, &wg2)
	}
}

func parse(s stream, wg *sync.WaitGroup) {
	defer wg.Done()
	defer log.Println("Parser finished.")
	for line := range s.lines {
		InsertInfo(s, line)
	}
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
