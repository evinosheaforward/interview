package filestat

import (
	"hash/fnv"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

func IngestFiles(conn *dbconn) {
	fileNames := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)
	go StartRecievers(&fileNames, conn, &wg)
	FileStreams(
		os.Getenv("INPUT_DIR"),
		fileNames)
	log.Println("Waiting on recievers.")
	wg.Wait()
}

func Ingest(st Streamer, pg *pgdb) {
	log.Println("Starting ingest subroutines.")
	defer log.Println("Ingesting finished.")
	var wg sync.WaitGroup
	wg.Add(1)
	go startParsers(st, pg, &wg)
	st.Stream()
	wg.Wait()
}

// Start subroutines for parsers for a given file
func startParsers(st Streamer, pg *pgdb, extwg *sync.WaitGroup) {
	defer extwg.Done()
	nParsers, _ := strconv.Atoi(os.Getenv("NUM_PARSERS"))
	var wg sync.WaitGroup
	wg.Add(nParsers)
	s := st.(FileStream)
	for i := 0; i < nParsers; i++ {
		go parse(s, pg, &wg)
	}
	wg.Wait()
}

// The magic happens with range s.lines
// that is what recieves on the channel
// channel is sent to from inside Stream()
func parse(s FileStream, pg *pgdb, wg *sync.WaitGroup) {
	defer log.Println("Parser finished.")
	defer wg.Done()
	for line := range s.lines {
		InsertInfo(pg, &line)
	}
}

// Accepts conn, does the thing.
// If added a CassDB then would already work
func InsertInfo(pg *pgdb, line *string) {
	h := hash(*line)
	nc := len(*line)
	nt := len(strings.Fields(*line))
	//var wg sync.WaitGroup
	dbWrapper := *pg
	//wg.Add(2)
	dbWrapper.InsertCounts(h, nc, nt)//, &wg)
	dbWrapper.InsertKeywords(line)//, &wg)
	//wg.Wait()
}

// simple hashing function for storing lines
// could think about a bigger hash depending on scale
func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
