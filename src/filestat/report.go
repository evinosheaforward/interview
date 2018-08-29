package filestat

import (
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
)

func Report(db dbconn) {
	conn, _ := db.(pgdb)
	// slight speedup (depending on bottleneck)
	var wg sync.WaitGroup
	wg.Add(1)
	go WriteStats(conn, &wg)
	kwds := conn.SelectKeywords()
	wg.Wait()
	WriteKeywords(kwds)
	log.Println("Report Finished.")
}

func WriteStats(conn dbconn, wg *sync.WaitGroup) {
	db, _ := conn.(pgdb)
	f, err := os.OpenFile(os.Getenv("OUTFILE"),
		os.O_WRONLY,
		0644)
	check(err)
	defer f.Close()
	countData := db.SelectCounts()
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
	wg.Done()
}

func WriteDupes(countData lineData, f *os.File) {
	outline := fmt.Sprintf("num dupes\t%d\n", countData.CountDupes())
	if _, err := f.WriteString(outline); err != nil {
		check(err)
	}
}

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

func WriteKeywords(kwds map[string]int) {
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
