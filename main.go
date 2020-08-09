package main

import (
	"fmt"
	"net/http"
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
		data, err := GetVideo()
		if err != nil {
			c.Send("failed to get video")
		} else {
			c.Send(data)
		}
	})

	app.Listen(3000)
}
