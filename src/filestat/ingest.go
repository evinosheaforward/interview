package filestat

import (
	"bufio"
  "io"
  "os"
)

func RunFile
	// maybe pass the open file instead
	file, err := os.Open(f)
	if err != nil {
			log.Fatal(err)
	}
	defer file.Close()

type stream interface {
	Read()
}

type stream struct {
	chan string lines
}

func (s stream) Read() {
		defer close(s)
		scanner := bufio.NewScanner(s)
		for scanner.Scan() {
				lines <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
				log.Fatal(err)
		}
}

func Ingest(f io.Reader) error {
  lines := make(chan string)
	defer close(lines)
	//var wg sync.WaitGroup
	s := stream{ lines }
  go startParsers(lines)
  go s.Read()
}

func startParsers(lines <-chan string) {
    ln := make(chan string)
		defer close(ln)
		num_parsers = 2
		p = parser{ ln }
		for i := range num_parsers {
			go p.Parse(lines)
		}

type parser interface {
	Parse()
	Info()
	Keywords()
}

type parser struct {
	string line
}

func (p parser) Parse(lines <-chan string) error {
		for line := range lines {
				ln <- line
				//?go storeInfo(ln)
				//?go ingestKwds(ln)
				p.Info(ln)
				p.Keywords(ln)
		}
}

func (p parser) Info(line string) error {
	h = hash(line)

	nc = len(line)
	nt = len(tokenize(line))
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func (ln line) Record(ln string) {
	h = hash(ln)
	if check(h) {
		nc = len(ln)
		nt = len(tokenize(ln))
		go kwdInfo(*ln)
	}
	  updateCount(h)
}

func updateCount(h int) {
	//update the counts wherever kwd comes up
}

func kwdInfo(ln *string, h int) {
  for kwd := range keywords {
		if strings.Contains(ln, kwd) {
        err = addHash(h)
    }
	}
}

func addHash(h int) error {

}
