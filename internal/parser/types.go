package parser

import (
	"tweet-audit/model"
	"tweet-audit/writer"
)

const USERNAME = "gboye_tobiloba"

type Tweet struct {
	ID       string `json:"id_str"`
	FullText string `json:"full_text"`
}

type ArchiveTweet struct {
	Tweet Tweet `json:"tweet"`
}
type Processor interface {
	Parse(path string) error
}

type ContentParser struct {
	client  model.AIClient
	writer  writer.Writer
	sleepFn func()
}
