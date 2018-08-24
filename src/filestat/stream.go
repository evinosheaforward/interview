package filestat

/*
import (
  "io"
)

// The queue of files
type fileStreamer interface {
  push()
  pop()
}

// Just the list of file
type fileQueue struct {
  queue []io.Reader()
}

// Stream function could be more complicated
// Function that can be called externally to get file from queue
// Want to have separate process add files to queue
func (f fileQueue) Stream() io.Reader() {
  //Get the file
  return f.queue.pop()
}

/* Possible future output channel for fileQueue...
func sq(out <-chan int) <-chan int {
    out <- //send something to out ch
}
*/
