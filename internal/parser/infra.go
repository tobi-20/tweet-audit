package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"tweet-audit/internal/loader"
	"tweet-audit/model"
)

func FlagTweet(post string) bool {
	post = strings.ToLower(post)
	keywords := []string{"feminism", "fuck", "racism", "pride", "arsenal", "simp", "republican", "malema", "Candace Owens", "jordan Peterson", "slander", "crazy", "kill"}

	for _, v := range keywords {
		if strings.Contains(post, strings.ToLower(v)) {
			return false
		}
	}
	return true
}

func ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path) // ReadFile reads the named file and returns the contents, data is of type []byte
	if err != nil {
		return nil, err
	}
	return data, nil
}

func ParseJson(data string) ([]ArchiveTweet, error) {
	var tweets []ArchiveTweet
	fileContents := data

	start := strings.Index(fileContents, "[") //returns the index of "["
	if start == -1 {
		return nil, errors.New("No json array found")
	}

	part := fileContents[start:] // file contents from "[" to the end
	return tweets, json.Unmarshal([]byte(part), &tweets)

}

func FormatPrompt(text string) string {

	return fmt.Sprintf(`Give me a json response for the prompt you get. The response takes this shape: {"decision": "action"}. action is either "keep" or "delete". 
 The rationale behind your responses should be you keep:
 -neutral Content
 -football Content
 -software engineering or computer science related Content.
 
 Then you delete:
 -content that  consist of swear words.
 -content that contain pidgin english.
 -political content.
 -content that contain violent words. %v `, text)
}

func (p *ContentParser) GetResponse(instr string) (string, error) {
	var response model.ModelResponse
	if instr == "" {

	}
	res, err := p.client.Analyze(instr)
	if err != nil {
		return "", err

	}
	if err := json.Unmarshal([]byte(res), &response); err != nil {
		return "", err

	}

	return response.Decision, nil
}

func shouldDelete(decision string) (bool, error) {
	if decision == "delete" {
		return true, nil
	}
	if decision == "keep" {
		return false, nil
	}
	return false, errors.New("unexpected response")
}

func (p *ContentParser) ProcessTweets(i int, tweets []ArchiveTweet) error {
	txt := tweets[i].Tweet.FullText
	//
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
