package bot

import (
	"regexp"
	"strings"

	"fmt"
	"github.com/VG-Tech-Dojo/vg-1day-2017-05-27/kaneko/env"
	"github.com/VG-Tech-Dojo/vg-1day-2017-05-27/kaneko/model"
	"net/url"
	//"os"
	"encoding/json"
	"log"
)

const (
	keywordApiUrlFormat = "https://jlp.yahooapis.jp/KeyphraseService/V1/extract?appid=%s&sentence=%s&output=json"
	talkApiUrl = "https://api.a3rt.recruit-tech.co.jp/talk/v1/smalltalk"
)

type (
	// Processor はmessageを受け取り、投稿用messageを作るインターフェースです
	Processor interface {
		Process(message *model.Message) *model.Message
	}

	// HelloWorldProcessor は"hello, world!"メッセージを作るprocessorの構造体です
	HelloWorldProcessor struct{}

	// OmikujiProcessor は"大吉", "吉", "中吉", "小吉", "末吉", "凶"のいずれかをランダムで作るprocessorの構造体です
	OmikujiProcessor struct{}

	// メッセージ本文からキーワードを抽出するprocessorの構造体です
	KeywordProcessor struct{}

	GachaProcessor struct{}

	TalkProcessor struct {}
)
type talkApiResults struct {
	Perplexity string `json:"perplexity"`
	Reply string `json:"reply"`
}
type talkApiResponce struct {
	Status       int    `json:"status"`
	Message     string `json:"message"`
	Results []talkApiResults `json:"results"`
}

// Process は"hello, world!"というbodyがセットされたメッセージのポインタを返します
func (p *HelloWorldProcessor) Process(msgIn *model.Message) *model.Message {
	return &model.Message{
		Body: msgIn.Body + ", world!",
	}
}

// Process は"大吉", "吉", "中吉", "小吉", "末吉", "凶"のいずれかがbodyにセットされたメッセージへのポインタを返します
func (p *OmikujiProcessor) Process(msgIn *model.Message) *model.Message {
	fortunes := []string{
		"大吉",
		"吉",
		"中吉",
		"小吉",
		"末吉",
		"凶",
	}
	result := fortunes[randIntn(len(fortunes))]
	return &model.Message{
		Body: result,
	}
}

// ガチャボット
func (p *GachaProcessor) Process(msgIn *model.Message) *model.Message {
	fortunes := []string{
		"SSレア",
		"Sレア",
		"レア",
		"ノーマル",
	}
	result := fortunes[randIntn(len(fortunes))]
	return &model.Message{
		Body: result,
	}
}

// Process はメッセージ本文からキーワードを抽出します
func (p *KeywordProcessor) Process(msgIn *model.Message) *model.Message {
	r := regexp.MustCompile("\\Akeyword (.*)\\z")
	matchedStrings := r.FindStringSubmatch(msgIn.Body)
	text := matchedStrings[1]

	url := fmt.Sprintf(keywordApiUrlFormat, env.KeywordApiAppId, text)

	json := map[string]int{}
	get(url, &json)

	keywords := []string{}
	for keyword := range map[string]int(json) {
		keywords = append(keywords, keyword)
	}

	return &model.Message{
		Body: "キーワード：" + strings.Join(keywords, ", "),
	}
}

// Process はメッセージ本文からキーワードを抽出します
func (p *TalkProcessor) Process(msgIn *model.Message) *model.Message {
	r := regexp.MustCompile("\\Atalk (.*)\\z")
	matchedStrings := r.FindStringSubmatch(msgIn.Body)
	text := matchedStrings[1]

	// json := map[string]int{}
	result_json := []byte{}
	values := url.Values{}
	values.Add("apikey", env.TalkApiAppId)
	values.Add("query", text)
	post(talkApiUrl, values, &result_json)

	fmt.Printf("%#v", result_json)
	var responce []talkApiResponce
	if err := json.Unmarshal(result_json, &responce); err != nil {
		log.Fatal(err)
	}

	return &model.Message{
		Body: responce[0].Results[0].Reply,
	}
}
