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
	resp, err := http.Get("https://9anime.to/watch/tower-of-god-dub.kvjr/ll3pknq")
	if err != nil {
		return "", err
	}

	// Convert HTML into goquery document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// Find something
	out := ""
	doc.Find("#controls").Each(func(i int, s *goquery.Selection) {
		playerIframe := s.Find(".autoplay").Text()

		out = fmt.Sprintf("player with id: %d has content %s\n", i, playerIframe)
	})

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
		data, err := GetVideo()
		if err != nil {
			c.Send("failed to get video")
		} else {
			c.Send(data)
		}
	})

	app.Listen(3000)
}
