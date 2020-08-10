package main

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

func GetData() (string, error) {

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

	out := "nothing found"
	// get the data-id from the player div
	dataId, ok := doc.Find("div#player").Attr("data-id")

	if !ok {
		return out, nil
	}
	// use data-id to a list of episodes
	resp, err = http.Get(fmt.Sprintf("https://9anime.to/ajax/film/servers?id=%s", dataId))
	if err != nil {
		return "", err
	}
	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// try to get the div containing the episodes hosted on streamtape
	doc.Find("div").Each(func(i int, s *goquery.Selection) {
		dataName, _ := s.Attr("data-name")

		matched, _ := regexp.MatchString(`40`, dataName)
		if matched {
			fmt.Println("data-name=", dataName)
			//fmt.Println(goquery.OuterHtml(s))

			fmt.Println(s.Find("a").Contents().Text())
		}

	})
	//fmt.Println(goquery.OuterHtml(episodesContainer))
	out = fmt.Sprintf("player with id: 1 has content %s\n", dataId)

	return out, nil
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

	app.Listen(3000)
}
