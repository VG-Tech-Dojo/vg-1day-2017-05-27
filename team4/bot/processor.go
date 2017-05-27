package bot

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/VG-Tech-Dojo/vg-1day-2017-05-27/team4/env"
	"github.com/VG-Tech-Dojo/vg-1day-2017-05-27/team4/model"
)

const (
	keywordApiUrlFormat = "https://webservice.recruit.co.jp/ab-road/spot/v1?key=%s&keyword=%s&format=json"
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

	SightSeeingProcessor struct{}

	SpotStruct struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Title       string `json:"title"`
		Url         string `json:"url"`
	}

	SightSeeingResult struct {
		ApiVersion string       `json:"api_version"`
		Spot       []SpotStruct `json:"spot"`
	}

	SightSeeingResponse struct {
		Results SightSeeingResult `json:"results"`
	}
)

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

func (p *SightSeeingProcessor) Process(msgIn *model.Message) *model.Message {
	r := regexp.MustCompile("\\Aspot (.*)\\z")
	matchedStrings := r.FindStringSubmatch(msgIn.Body)
	fmt.Printf("T:%s\n", matchedStrings[0])
	text := matchedStrings[1]

	json := SightSeeingResponse{}

	u := fmt.Sprintf(keywordApiUrlFormat, env.KeywordApiAppId, text)
	fmt.Printf("URL:%s\n", u)

	get(u, &json)

	fmt.Printf("V:%s", json.Results.ApiVersion)

	var body string
	if len(json.Results.Spot) > 0 {
		for _, res := range json.Results.Spot {
			if res.Name != "" {
				body += fmt.Sprintf("%s ", res.Name)
			}
		}
		body += "が有名です。"
	} else {
		body = text + "の観光地が見つかりませんでした。"
	}

	return &model.Message{
		Body: body,
	}
}
