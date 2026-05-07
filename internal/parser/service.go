package parser

import (
	"log"
	"os"
	"time"

	"tweet-audit/internal/loader"
	"tweet-audit/writer"

	"tweet-audit/model"
)

func NewContentParser(fileName string) (*ContentParser, error) {
	client, err := model.NewGeminiClient(os.Getenv("GEMINI_API_KEY"))
	if err != nil {
		return nil, err
	}

	w, err := writer.NewSetter(fileName)
	if err != nil {
		return nil, err
	}

	sleepFn := func() { time.Sleep(13 * time.Second) }

	return &ContentParser{
		client:  client,
		writer:  w,
		sleepFn: sleepFn,
	}, nil
}

func (p *ContentParser) Parse(path string) error {
	var tweets []ArchiveTweet
	//
	data, err := ReadFile(path)
	if err != nil {
		return err
	}
	//
	if tweets, err = ParseJson(string(data)); err != nil {
		return err
	}

	//  json.Unmarshal is the bridge between byte slice and tweets slice
	//https://x.com/gboye_tobiloba/status/204076442884966479

	//
	startIdx, err := loader.LoadProgress("progress.txt")
	if err != nil {
		return err
	}

	//
	for i := startIdx + 1; i <= 10; i++ {

		err := p.ProcessTweets(i, tweets)
		if err != nil {

			log.Println(err)
			continue
		}
		if i%10 == 0 {
			if err := p.writer.Flush(); err != nil {
				return err
			}
		}

		p.sleepFn()

	}

	return p.writer.Close()
}
