package parser

import (
	"errors"
	"os"
	"testing"
	"tweet-audit/internal/loader"
)

type MockClient struct {
	Response string
	Err      error
}

type MockWriter struct {
	tweetAddr    string
	shouldDelete bool
	Err          error
}

func (m *MockWriter) Set(tweetAddr string, shouldDelete bool) error {
	//fails fast perchance something goes wrong with the writer
	m.tweetAddr = tweetAddr
	m.shouldDelete = shouldDelete
	return m.Err
}
func (m *MockWriter) Flush() error {

	return m.Err
}
func (m *MockWriter) Close() error {

	return m.Err

}

type MockParser struct {
	Err error
}

func (m *MockClient) Analyze(instr string) (string, error) {

	return m.Response, m.Err
}

// Expect to not get errors
func TestParseJson_Happy(t *testing.T) {
	input := `[{
    "tweet" : {
      "id_str" : "2",
       "full_text" : "R"      
    }
  }]`
	tweets, err := ParseJson(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(tweets) != 1 {
		t.Fatalf("expected len to be 1 got %d instead", len(tweets))
	}

	if tweets[0].Tweet.ID != "2" || tweets[0].Tweet.FullText != "R" {
		t.Fatalf("expected id to be 2 and full text to be R, got id to be %v and full text to be %v", tweets[0].Tweet.ID, tweets[0].Tweet.FullText)
	}
}
func TestParseJson_InvalidJson(t *testing.T) {
	input := `[{
    "tweet" : {
      "id_str" : "2",
       "full_text" : "R"      
    }
  }`
	_, err := ParseJson(input)
	if err == nil {
		t.Fatal("expected error got nil instead")
	}

}

func TestLoadProgress_Happy(t *testing.T) {
	path := "progress_test.txt"
	if err := os.WriteFile(path, []byte("1"), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(path)
	n, err := loader.LoadProgress(path)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("expected %d got %d", 1, n)
	}
}
func TestLoadProgress_InvalidCounter(t *testing.T) {
	path := "progress_test1.txt"
	if err := os.WriteFile(path, []byte("1i"), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(path)
	_, err := loader.LoadProgress(path)
	if err == nil {
		t.Fatal("expected error but got none")
	}

}

func TestProcessTweets_Happy(t *testing.T) {
	tweets := []ArchiveTweet{{
		Tweet: Tweet{
			ID:       "2",
			FullText: "R",
		},
	}}
	i := 0
	writer := &MockWriter{}
	p := &ContentParser{
		client: &MockClient{
			Response: `{"decision":"keep"}`,
			Err:      nil,
		},
		writer:  writer,
		sleepFn: func() {},
	}

	if err := p.ProcessTweets(i, tweets); err != nil {
		t.Fatal(err)
	}
	if writer.tweetAddr != "https://x.com/gboye_tobiloba/status/2" {
		t.Fatalf("expected https://x.com/gboye_tobiloba/status/2 got %v", writer.tweetAddr)
	}
	if writer.shouldDelete {
		t.Fatalf("expected false got %v", writer.shouldDelete)
	}
}
func TestProcessTweets_InvalidID(t *testing.T) {
	tweets := []ArchiveTweet{{
		Tweet: Tweet{
			ID:       "",
			FullText: "R",
		},
	}}
	i := 0
	writer := &MockWriter{}
	p := &ContentParser{
		client: &MockClient{
			Response: `{"decision":"keep"}`,
			Err:      nil,
		},
		writer:  writer,
		sleepFn: func() {},
	}

	if err := p.ProcessTweets(i, tweets); err == nil {
		t.Fatal("expected error for invalid tweet ID")
	}

}
func TestProcessTweets_InvalidFullText(t *testing.T) {
	tweets := []ArchiveTweet{{
		Tweet: Tweet{
			ID:       "444444444234556789",
			FullText: "",
		},
	}}
	i := 0
	writer := &MockWriter{}
	p := &ContentParser{
		client: &MockClient{
			Response: `{"decision":"keep"}`,
			Err:      nil,
		},
		writer:  writer,
		sleepFn: func() {},
	}

	if err := p.ProcessTweets(i, tweets); err == nil {
		t.Fatal("expected error for invalid tweet ID")
	}

}
func TestGetResponse_Happy(t *testing.T) {
	p := &ContentParser{
		client: &MockClient{
			Response: `{"decision":"keepg"}`,
			Err:      nil,
		},
	}
	instr := "gg"
	d, err := p.GetResponse(instr)
	if err != nil {
		t.Fatal(err)
	}
	if d != "keepg" {
		t.Fatalf("expected keepg got %v", d)
	}
}
func TestGetResponse_InvalidPrompt(t *testing.T) {
	p := &ContentParser{
		client: &MockClient{
			Response: `{"decision":"keepg"}`,
			Err:      nil,
		},
	}
	instr := ""
	_, err := p.GetResponse(instr)
	if err == nil {
		t.Fatal(err)
	}

}

func TestShouldDelete_Happy(t *testing.T) {
	ok, err := shouldDelete("keep")
	if err != nil {
		t.Fatal(err)
	}
	if ok != false {
		t.Fatal(err)
	}
}
func TestShouldDelete_Fail(t *testing.T) {
	_, err := shouldDelete("yes")
	if err == nil {
		t.Fatal(err)
	}
}

func TestProcessTweets_WriteFailed(t *testing.T) {

	tweets := []ArchiveTweet{{
		Tweet: Tweet{
			ID:       "2",
			FullText: "R",
		},
	}}
	i := 0
	setter := &MockWriter{
		Err: errors.New("Write failed"),
	}

	p := &ContentParser{
		client: &MockClient{
			Response: `{"decision":"keep"}`,
			Err:      nil,
		},
		writer:  setter,
		sleepFn: func() {},
	}

	if err := p.ProcessTweets(i, tweets); err == nil {
		t.Fatal("expected error got none")
	}

}
