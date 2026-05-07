package writer

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"
)

func NewSetter(path string) (*Setter, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	w := csv.NewWriter(f)
	info, err := f.Stat()
	if info.Size() == 0 {
		w.Write([]string{"tweet_url", "deleted"})
	}
	if err != nil {
		return nil, err
	}

	return &Setter{
		file:   f,
		writer: w,
	}, nil
}

func (s *Setter) Set(tweetAddr string, shouldDelete bool) error {
	//fails fast perchance something goes wrong with the writer
	if tweetAddr == "" {
		return errors.New("cannot write invalid string literal")
	}
	return s.writer.Write([]string{tweetAddr, strconv.FormatBool(shouldDelete)})
}
func (s *Setter) Flush() error {

	s.writer.Flush()
	return s.writer.Error()
}
func (s *Setter) Close() error {

	s.writer.Flush()
	return s.file.Close()

}
