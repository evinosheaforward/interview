package main

import (
	"log"
	"time"

	"filestat"
)

func main() {
	// Make sure DB setup
	time.Sleep(10 * time.Second)
	conn := filestat.NewDBConn()
	filestat.SetupDB(conn)
	filestat.IngestFiles(&conn)
	// Wait for ingests to finish
	log.Println("All Ingest Finished")
	time.Sleep(10 * time.Second)
	filestat.Report(conn)
}
