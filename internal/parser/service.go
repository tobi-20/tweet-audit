package parser

import (
	"os"
	"strconv"
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
		client:              client,
		writer:              w,
		sleepFn:             sleepFn,
		checkpointPath:      "progress.txt",
		FailedToProcessPath: "skipped.txt",
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
	startIdx, err := loader.LoadProgress(p.checkpointPath)
	if err != nil {
		return err
	}

	//
	for i := startIdx; i <= 10; i++ {

		err := p.ProcessTweets(i, tweets)
		if err != nil {
			if err := os.WriteFile(p.FailedToProcessPath, []byte(strconv.Itoa(i)), 0644); err != nil {
				return err
			}
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
