package writer

import (
	"encoding/csv"
	"os"
)

type Writer interface {
	Set(tweetAddr string, shouldDelete bool) error
	Flush() error
	Close() error
}

type Setter struct {
	file   *os.File
	writer *csv.Writer
}
