package filestat

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"sync"
)

type Streamer interface {
	Stream()
}

type FileStream struct {
	lines 		chan string
	fileName 	string
}

func (s FileStream) FileName() string {
	return s.fileName
}

// Stream closes the stream's channel
// scanner.Scan() sends on the channel to the parsers
// the channel recieves inside parse()
func (s FileStream) Stream() {
	defer close(s.lines)
	indir := os.Getenv("INPUT_DIR")
	file, err := os.Open(path.Join(indir, s.fileName))
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s.lines <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		check(err)
	}
}

func FileStreams(indir string, fileNames chan string) {
	defer close(fileNames)
	files, err := ioutil.ReadDir(indir)
	check(err)
	for _, f := range files {
		fileNames <- f.Name()
	}
}

func StartRecievers(fileNames *chan string, conn *dbconn, extwg *sync.WaitGroup) {
	stop := StartTimer("Ingesting all files")
	defer stop()
	defer extwg.Done()
	nReaders, _ := strconv.Atoi(os.Getenv("NUM_FILE_READERS"))
	var wg sync.WaitGroup
	wg.Add(nReaders)
	for i := 0; i < nReaders; i++ {
		go recieveFile(fileNames, conn, &wg)
	}
	wg.Wait()
}

func recieveFile(fileNames *chan string, conn *dbconn, wg *sync.WaitGroup) {
	defer wg.Done()
	c := *conn
	pg, _ := c.(pgdb)
	var s FileStream
	for fname := range *fileNames {
		s = FileStream{
			lines: 		make(chan string),
			fileName: fname,}
		Ingest(s, &pg)
	}
}
