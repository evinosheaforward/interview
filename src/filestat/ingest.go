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
		Stream(io.Reader)
}

type stream struct {
		lines       chan string
		db          pgdb
		//db          cassdb
}

func (s stream) Stream(file io.Reader) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s.lines <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func Ingest(fname string) {
	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	s := stream{
		lines:       make(chan string),
		db:          NewDBConn(),
	}
	defer close(s.lines)
	log.Println("Staring subroutines")
	defer log.Println("Subroutines Finished.")
	var wg sync.WaitGroup
	wg.Add(2)
	nparsers, _ := strconv.Atoi(os.Getenv("NUM_PARSERS"))
	go startParsers(s, nparsers)
	go s.Stream(f)
	wg.Wait()
}

func startParsers(s stream, num_parsers int) {
	var wg sync.WaitGroup
	wg.Add(num_parsers)
	for i := 0; i < num_parsers; i++ {
		go parse(s)
	}
	wg.Wait()
}

func parse(s stream) {
	for line := range s.lines {
		InsertInfo(s, line)
	}
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
