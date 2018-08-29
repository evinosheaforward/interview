package filestat

import (
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
)

func Report(conn dbconn) {
	stop := StartTimer("Reporting")
	defer stop()
	var wg sync.WaitGroup
	wg.Add(1)
	go WriteStats(conn, &wg)
	WriteKeywords(conn, &wg)
	log.Println("Report Finished.")
}

// This function does the heavy lifting
// grabs the data with the connection
// and calls the writeing methods
func WriteStats(conn dbconn, wg *sync.WaitGroup) {
	stop := StartTimer("Writing Statistics")
	defer stop()
	defer wg.Done()
	dbWrapper, _ := conn.(pgdb)
	f, err := os.OpenFile(os.Getenv("OUTFILE"),
		os.O_WRONLY,
		0644)
	check(err)
	defer f.Close()
	countData := dbWrapper.SelectCounts()
	lenData := make(statData, len(countData))
	tokenData := make(statData, len(countData))
	for i, countDatum := range countData {
		lenData[i] = statDatum{
			countDatum.NumChars(),
			countDatum.Count()}
		tokenData[i] = statDatum{
			countDatum.NumTokens(),
			countDatum.Count()}
	}
	WriteDupes(countData, f)
	WriteMedStd(lenData, f, "length")
	WriteMedStd(tokenData, f, "tokens")
}

// This is fine for memory because its not being copied
func WriteDupes(countData lineData, f *os.File) {
	outline := fmt.Sprintf("num dupes\t%d\n", countData.CountDupes())
	if _, err := f.WriteString(outline); err != nil {
		check(err)
	}
}

// This function just calculates the Std and Med and writes them to file
func WriteMedStd(data statData, f *os.File, name string) {
	outline := fmt.Sprintf("med %s\t%f\n", name, Median(data))
	if _, err := f.WriteString(outline); err != nil {
		check(err)
	}
	outline = fmt.Sprintf("std %s\t%f\n", name, Std(data))
	if _, err := f.WriteString(outline); err != nil {
		check(err)
	}
}

// Waits on the other writer
// Get the keywords from the database,
// sort the keys, then write the data.
func WriteKeywords(conn dbconn, wg *sync.WaitGroup) {
	dbWrapper, _ := conn.(pgdb)
	kwds := dbWrapper.SelectKeywords()
	wg.Wait()
	f, err := os.OpenFile(os.Getenv("OUTFILE"),
		os.O_APPEND|os.O_WRONLY,
		os.ModeAppend)
	check(err)
	defer f.Close()
	keys := make([]string, len(kwds))
	idx := 0
	for k, _ := range kwds {
		keys[idx] = k
		idx++
	}
	sort.Strings(keys)
	var outline string
	for _, k := range keys {
		outline = fmt.Sprintf("%s\t\t%d\n", k, kwds[k])
		if _, err := f.WriteString(outline); err != nil {
			log.Fatal(err)
		}
	}
}
