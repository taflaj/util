// linereader.go

package reader

import (
	"bufio"
	"log"
	"os"

	"github.com/golang-collections/collections/stack"
)

// LineReader reads lines from an input file
type LineReader struct {
	filename string
	done     bool
	buffer   chan string
	count    int
	stack    *stack.Stack
}

func read(lr *LineReader) {
	if file, err := os.Open(lr.filename); err != nil {
		log.Panic(err)
	} else {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lr.buffer <- scanner.Text()
			// lr.lines++
		}
		if err := scanner.Err(); err != nil {
			log.Panic(err)
		}
		close(lr.buffer)
		lr.done = true
	}
}

// NewLineReader creates a new LineReader with default buffer size
func NewLineReader(filename string) *LineReader {
	return NewLineReaderWithSize(filename, 100)
}

// NewLineReaderWithSize creates a new LineReader with the specified buffer size
func NewLineReaderWithSize(filename string, size int) *LineReader {
	lr := LineReader{}
	lr.filename = filename
	lr.done = false
	lr.buffer = make(chan string, size)
	lr.count = 0
	lr.stack = stack.New()
	go read(&lr)
	return &lr
}

// GetFilename returns the file name
func (lr *LineReader) GetFilename() string {
	return lr.filename
}

// IsDone returns the reader done status
func (lr *LineReader) IsDone() bool {
	return lr.done
}

// GetLines returns the number of lines read so far
func (lr *LineReader) GetLines() int {
	return lr.count
}

// ReadLine returns the next line from the file
func (lr *LineReader) ReadLine() (int, string, bool) {
	if lr.stack.Len() > 0 {
		return lr.count, lr.stack.Pop().(string), true
	}
	line, ok := <-lr.buffer
	lr.count++
	return lr.count, line, ok
}

// UnreadLine unreads a line, pushing it into the queue
func (lr *LineReader) UnreadLine(line string) {
	lr.stack.Push(line)
}
