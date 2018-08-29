package main

import (
    "fmt"
    "net/http"
    "log"
		"time"

		"filestat"
)

func main() {
		time.Sleep(7 * time.Second)
    http.HandleFunc("/", handle)
		db := filestat.NewDBConn()
		filestat.SetupDB(db)
		fname := "/data/input/simple.txt"
		fmt.Println("Ingesting file: ", fname)
		filestat.Ingest(fname, db)
		time.Sleep(7 * time.Second)
		filestat.Report(db)
		log.Fatal(http.ListenAndServe(":8080", nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Docker")
	fmt.Println("OH DOH DOH")
	log.Println("OH DOH DOH")
}

	/*
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
	*/
