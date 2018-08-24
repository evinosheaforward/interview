package main

import (
  "fmt"
  "io"
  "https://github.com/evinosheaforward/interview/tree/master/src/filestat"
)

//main function for running fileStat process
//have either fileHandler or report come through channel
func main() {
	/*
	x =: make(chan io.Reader)
  y =: make(chan io.Writer)

  go func() { r <- FileStreamer() }
  go func() { w <- NewsManager()}
  while true:
    // I want FileStreamer and NewsManager to run concurrently
    select {
    case file <- r:
      // run the ingest code for the file
      go Ingest(file)
      fmt.PrintLn("Ingesting file: %s", file)
    case reporter <- w:
      err = reporter.Report()
      if err {
        fmt.PrintLn("Uh oh!")
      }
      break
    }
	*/
  fmt.PrintLn(filestat.Ingest())
}
