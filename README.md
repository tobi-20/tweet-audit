# Go Tweet Processor — Architecture & Design Decisions

A system for sequential AI-powered tweet processing with checkpointing.
---
## Architecture Overview

```
Data → Parser → Brain (ProcessTweets) → External IO (AI + Writer)
```

## Core Structure

### AI Abstraction

```go
type AIClient interface {
    Generate(prompt string) (string, error)
}
```

Depend on the interface, not the implementation. Swapping clients means just a new struct, not touching the service.

### Process tweets Layer

```go
func (p *ContentParser) ProcessTweets(i int, tweets []ArchiveTweet) error {
	txt := tweets[i].Tweet.FullText
	//
	if txt == "" {
		return errors.New("the full text is empty")
	}
	if tweets[i].Tweet.ID == "" {
		return errors.New("tweet id missing")
	}

	if USERNAME == "" {
		return errors.New("username missing")
	}
	instr := FormatPrompt(txt)
	tweetAddr := "https://x.com/" + USERNAME + "/status/" + tweets[i].Tweet.ID

	//l
	log.Println(tweetAddr)

	resp, err := p.GetResponse(instr)
	if err != nil {
		return err
	}
	decision, err := shouldDelete(resp)
	if err != nil {
		return err
	}
	//

	if err := p.writer.Set(tweetAddr, decision); err != nil {
		return err
	}
	//
	if err := loader.SaveProgress(i); err != nil {
		return err
	}
	return nil
}

```

Single responsibility: Delegate AI work, persist progress.

### Entry Point

```go
func main() {
    p, err := parser.NewContentParser("flagged.csv")

	path := os.Args[1]
	err = p.Parse(path)
}
```

---