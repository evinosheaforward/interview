package main

import (
  "fmt"
  "io"
  //filestat
)

//main function for running fileStat process
//have either fileHandler or report come through channel
func main() {
  x =: make(chan io.Reader)
  y =: make(chan io.Writer)

  go func() { r <- Stream() }
  go func() { w <- Manage()}
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
  fmt.PrintLn("Done")
}
