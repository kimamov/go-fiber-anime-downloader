package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
	"github.com/gofiber/fiber"
	"github.com/gofiber/template/html"
)

func GetJsonDocument(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// parse json from res.body
	var result map[string]string
	json.Unmarshal(body, &result)
	// use utfBody using goquery
	fmt.Println(result["html"])
	resultReader := strings.NewReader(result["html"])
	doc, err := goquery.NewDocumentFromReader(resultReader)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func GetDocument(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Convert the designated charset HTML to utf-8 encoded HTML.
	// `charset` being one of the charsets known by the iconv package.
	utfBody, err := iconv.NewReader(res.Body, "UTF-8", "utf-8")
	if err != nil {
		return nil, err
	}

	// use utfBody using goquery
	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func GetData() (map[int]string, error) {

	// Get the HMTML
	animeDocument, err := GetDocument("https://9anime.to/watch/tower-of-god-dub.kvjr/ojo9nqz")
	if err != nil {
		return nil, err
	}
	// get the data-id from the player div
	dataId, ok := animeDocument.Find("div#player").Attr("data-id")
	if !ok {
		return nil, nil
	}
	// use data-id to a list of episodes
	animeStreamsDocument, err := GetJsonDocument(fmt.Sprintf("https://9anime.to/ajax/film/servers?id=%s", dataId))
	if err != nil {
		return nil, err
	}

	// try to get the div containing the episodes hosted on streamtape
	var episodesMap map[int]string
	episodesMap = make(map[int]string)
	container := animeStreamsDocument.Find(`div[data-id="40"]`)
	container.Find("a").Each(func(i int, s *goquery.Selection) {
		episodeId, ok := s.Attr("data-id")
		if ok {
			episodesMap[i] = episodeId
		}
	})
	//fmt.Println(episodesMap)

	return episodesMap, nil
}

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(&fiber.Settings{
		Views: engine,
	})

	app.Get("/api/video", func(c *fiber.Ctx) {
		c.Send("succes")
		data, err := GetData()
		if err != nil {
			c.Send("failed to get video")
		} else {
			c.Send(data)
		}
	})

	//app.Static("/", "./public")
	app.Get("/*", func(c *fiber.Ctx) {
		_ = c.Render("index", fiber.Map{
			"Title": "hey there! ;)",
		})
	})

	app.Listen(3000)
}
