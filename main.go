package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Print("specify keywork and webhook url.")
		os.Exit(1)
	}

	// var keyword string = "ビジネス"
	// var webhook string = config.Webhook

	var keyword string = os.Args[1]
	var webhook string = os.Args[2]

	today := time.Now()
	const dayLayout = "20060102"
	todayStr := today.Format(dayLayout)

	fmt.Print(todayStr + "\n")
	doc, err := goquery.NewDocument("https://tv.yahoo.co.jp/search/?q=" + keyword + "&d=" + todayStr)
	if err != nil {
		fmt.Print("document not found. ")
		os.Exit(1)
	}

	postText := ""
	linkUrl := ""
	postText += "本日のTV番組情報！"
	postText += "検索ワード : " + keyword + "\n"

	doc.Find(".programlist > li").Each(func(_ int, s *goquery.Selection) {
		postText += "_"
		s.Find(".leftarea > p > em").Each(func(_ int, em *goquery.Selection) {
			postText += em.Text() + " "
		})
		atag := s.Find(".rightarea > p > a").First()
		postText += s.Find(".rightarea > p > span").First().Text()
		postText += "_"
		postText += " *[" + atag.Text() + "]* :"
		linkUrl, _ = atag.Attr("href")
		postText += s.Find(".rightarea > p").Filter(":not(:has(a,span))").Text()

		postText += "\n" + "> https://tv.yahoo.co.jp" + linkUrl + "\n"
	})

	if postText == "" {
		fmt.Print("program not found. ")
		os.Exit(0)
	}

	fmt.Print(postText)

	jsonStr := "{\"text\":\"" + postText + "\"}"
	fmt.Print(jsonStr)
	req, _ := http.NewRequest(
		"POST",
		webhook,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err)
	}
	defer resp.Body.Close()
}
