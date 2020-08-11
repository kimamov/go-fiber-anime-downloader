package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
	"github.com/gofiber/fiber"
)

func SortString(s string) string {
	charArray := strings.Split(s, "")
	sort.Strings(charArray)
	return strings.Join(charArray, "")
}

func GetVideo() (string, error) {

	// Get the HMTML
	resp, err := http.Get("https://9anime.to/watch/tower-of-god-dub.kvjr/ojo9nqz")
	if err != nil {
		return "", err
	}

	// Convert HTML into goquery document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// Find something
	out := "nothing found"
	doc.Find("#player").Each(func(i int, s *goquery.Selection) {
		fmt.Println(goquery.OuterHtml(s))

		dataId, _ := s.Attr("data-id")
		// use dataId to get episodes and streams
		// https://9anime.to/ajax/film/servers?id={title_id}

		// after getting all the episodes search for the once on streamtape
		// use the episode id to get the url of the iframe
		// https://9anime.to/ajax/episode/info?id={self.url}&server=40

		// request the iframe and download the video
		out = fmt.Sprintf("player with id: %d has content %s\n", i, dataId)
	})

	/* iframeContainer := doc.Find("#player")

	iframeSrc, ok := iframeContainer.Find("iframe").Attr("src")

	if ok != true {
		return "", nil
	} */

	return out, nil
}

type JSONBody struct {
	Html string `json:html`
}

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

func GetDataA() (string, error) {

	// Get the HMTML
	animeDocument, err := GetDocument("https://9anime.to/watch/tower-of-god-dub.kvjr/ojo9nqz")
	if err != nil {
		return "", err
	}
	out := "nothing found"
	// get the data-id from the player div
	dataId, ok := animeDocument.Find("div#player").Attr("data-id")
	if !ok {
		return out, nil
	}
	// use data-id to a list of episodes
	animeStreamsDocument, err := GetDocument(fmt.Sprintf("https://9anime.to/ajax/film/servers?id=%s", dataId))
	if err != nil {
		return "", err
	}
	return fmt.Sprintln(goquery.OuterHtml(animeStreamsDocument.Contents())), nil
	//return fmt.Sprintln(goquery.OuterHtml(animeStreamsDocument.Contents())), nil
	// try to get the div containing the episodes hosted on streamtape
	//var episodesMap map[int]string
	episodesMap := make(map[int]string)
	animeStreamsDocument.Find("div").Each(func(i int, s *goquery.Selection) {
		dataName, _ := s.Attr("data-name")
		//fmt.Println(goquery.OuterHtml(s))
		//matched, _ := regexp.MatchString(`40`, dataName)
		if strings.Contains(dataName, "40") {
			fmt.Println("data-name=", dataName)
			fmt.Println("data-name=", i)
			streamLinks := s.Children()
			out = fmt.Sprintln(goquery.OuterHtml(streamLinks))
			streamLinks.Find("a").Each(func(i int, s *goquery.Selection) {
				episodesMap[i] = fmt.Sprintln( /* goquery.OuterHtml(s) */ i)
			})
			//fmt.Println(episodesMap)
		}

	})
	//fmt.Println(goquery.OuterHtml(episodesContainer))
	//out = fmt.Sprintf("player with id: 1 has content %s\n", dataId)

	return out, nil
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
	app := fiber.New()

	app.Get("/anagram/:firstword/:secondword", func(c *fiber.Ctx) {
		sortedWordOne := SortString(c.Params("firstword"))
		sortedWordTwo := SortString(c.Params("secondword"))

		if sortedWordOne == sortedWordTwo {
			c.Send("is anagram")
		} else {
			c.Send("is no anagram")
		}
	})

	app.Get("/video", func(c *fiber.Ctx) {
		c.Send("succes")
		data, err := GetData()
		if err != nil {
			c.Send("failed to get video")
		} else {
			c.Send(data)
		}
	})

	app.Get("/proxyanime", func(c *fiber.Ctx) {
		resp, err := http.Get("https://9anime.to/watch/terror-in-resonance-dub.818n/800n1v3")
		if err != nil {
			c.Send("failed to get")
		}
		c.Send(resp.Body)
	})

	app.Get("/proxydata", func(c *fiber.Ctx) {
		resp, err := GetJsonDocument("https://9anime.to/ajax/film/servers?id=kvjr")
		if err != nil {
			c.Send("failed to get")
		}
		a := resp.Find(`div[data-id="40"]`)
		c.Send(goquery.OuterHtml(a))
	})

	app.Listen(3000)
}
