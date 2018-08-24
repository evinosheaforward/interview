package filestat

import (
	"bufio"
  "io"
  "os"
)

func Ingest(f io.Reader) {
  //magically run some parallelized / concurrent file reading
}

type chucker interface {
	readChunk()
}

//This is not right at all... but outlines the idea
func readChunk(f *os.File, start int, stop int) {
  for line := range f.Read(start, stop) {
		//call process to put length of line, hash? the line, num tokens in db
		go storeInfo(line)
		//call another process to handle the kwd stuff for the line and store
		go findKeywords(line)
  }
	/*
	scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        jobs <- scanner.Text()
    }
    close(jobs)
	*/
}

func (i ingester) ReadWrite() error {
	// want to break into multiple file readers
	// when a filereader recieved on write
  go readChunk()

	// Collect all the results...
	// First, make sure we close the result channel when everything was processed
	go func() {
   wg.Wait()
   close(results)
  }()

  // Now, add up the results from the results channel until closed
  counts := 0
  for v := range results {
   counts += v
  }
